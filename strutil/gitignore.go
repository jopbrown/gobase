package strutil

import (
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/jopbrown/gobase/errors"
	"github.com/jopbrown/gobase/must"
)

type GitIgnoreMatcher struct {
	re  *regexp.Regexp
	raw string

	isAlwaysNotMatch bool
	isNegative       bool

	folderToMatch *regexp.Regexp
}

const (
	patternStartDoubleAsterisk = `([^/]*/)*`
	patternEndDoubleAsterisk   = `(/[^/]*)*`
	patternSingleAsterisk      = `([^/]*)`
	patternQuestion            = `([^/])`
)

/*
CompileGitIgnore compile gitignore pattern into regexp.

This implement not support rule 3.1, 4.1 and 9.3.

The pattern format rule copy from https://git-scm.com/docs/gitignore#_pattern_format

1. A blank line matches no files, so it can serve as a separator for readability.

2. A line starting with # serves as a comment.
 1. Put a backslash ("\") in front of the first hash for patterns that begin with a hash.

3. Trailing spaces are ignored.
 1. unless they are quoted with backslash ("\").

4. An optional prefix "!" which negates the pattern;
 1. any matching file excluded by a previous pattern will become included again. It is not possible to re-include a file if a parent directory of that file is excluded. Git doesn’t list excluded directories for performance reasons, so any patterns on contained files have no effect, no matter where they are defined.
 2. Put a backslash ("\") in front of the first "!" for patterns that begin with a literal "!", for example, "\!important!.txt".

5. The slash / is used as the directory separator. Separators may occur at the beginning, middle or end of the .gitignore search pattern.

6. If there is a separator at the beginning or middle (or both) of the pattern, then the pattern is relative to the directory level of the particular .gitignore file itself. Otherwise the pattern may also match at any level below the .gitignore level.

7. If there is a separator at the end of the pattern then the pattern will only match directories, otherwise the pattern can match both files and directories.

8. For example, a pattern doc/frotz/ matches doc/frotz directory, but not a/doc/frotz directory; however frotz/ matches frotz and a/frotz that is a directory (all paths are relative from the .gitignore file).

9. An asterisk "*" matches anything except a slash. The character "?" matches any one character except "/". The range notation, e.g. [a-zA-Z], can be used to match one of the characters in a range. See fnmatch(3) and the FNM_PATHNAME flag for a more detailed description.

 1. A '?' (not between brackets) matches any single character.

 2. A '*' (not between brackets) matches any string, including the empty string.

 3. Other special rules. i.g. [A-Fa-f0-9]

10. Two consecutive asterisks ("**") in patterns matched against full pathname may have special meaning:

 1. A leading "**" followed by a slash means match in all directories. For example, "** /foo" matches file or directory "foo" anywhere, the same as pattern "foo". "** /foo/bar" matches file or directory "bar" anywhere that is directly under directory "foo".

 2. A trailing "/**" matches everything inside. For example, "abc/**" matches all files inside directory "abc", relative to the location of the .gitignore file, with infinite depth.

 3. A slash followed by two consecutive asterisks then a slash matches zero or more directories. For example, "a/** /b" matches "a/b", "a/x/b", "a/x/y/b" and so on.

 4. Other consecutive asterisks are considered regular asterisks and will match according to the previous rules.
*/
func CompileGitIgnore(line string) (*GitIgnoreMatcher, error) {
	var err error

	matcher := &GitIgnoreMatcher{}
	matcher.raw = line

	// rule 3
	expr := strings.TrimSpace(line)

	// rule 1, 2
	if expr == "" || expr[0] == '#' {
		matcher.isAlwaysNotMatch = true
		return matcher, nil
	}

	// rule 4
	if expr[0] == '!' {
		matcher.isNegative = true
		expr = expr[1:]
	} else {
		// rule 2.1, 4.2
		if len(expr) > 1 && expr[0] == '\\' && (expr[1] == '#' || expr[1] == '!') {
			expr = expr[1:]
		}
	}

	// rule 6
	if i := strings.Index(expr, `/`); i >= 0 && i != len(expr)-1 {
		if expr[0] == '/' {
			expr = expr[1:]
		}
	} else {
		expr = "**/" + expr
	}

	// rule 7
	if expr[len(expr)-1] == '/' {
		matcher.folderToMatch, err = ComplieGlob(path.Base(expr))
		if err != nil {
			return nil, errors.ErrorAt(err)
		}

		if !strings.HasSuffix(expr[:len(expr)-1], "/**") {
			expr += "**"
		}
	} else if strings.HasSuffix(expr, "/**") {
		// rule 10.2
		matcher.folderToMatch, err = ComplieGlob(path.Base(expr[:len(expr)-3]))
		if err != nil {
			return nil, errors.ErrorAt(err)
		}
	} else {
		expr += "/**"
	}

	expr = regexp.QuoteMeta(expr)

	expr = `^` + expr + `$`

	// rule 10.2
	expr = replaceGlobRule(expr, "/**", `/\*\*`, patternEndDoubleAsterisk)

	// rule 10.1
	expr = replaceGlobRule(expr, "**/", `\*\*/`, patternStartDoubleAsterisk)

	// rule 9.2
	expr = replaceGlobRule(expr, "*", `\*`, patternSingleAsterisk)

	// rule 9.1
	expr = replaceGlobRule(expr, "?", `\?`, patternQuestion)

	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, errors.ErrorAt(err)
	}
	matcher.re = re

	return matcher, nil
}

func MustCompileGitIgnore(line string) *GitIgnoreMatcher {
	return must.Value(CompileGitIgnore(line))
}

func CompileGitIgnoreLines(lines ...string) (MultiMatcher, error) {
	mm := make(MultiMatcher, 0, len(lines))
	for _, line := range lines {
		m, err := CompileGitIgnore(line)
		if err != nil {
			return nil, errors.ErrorAt(err)
		}

		mm = append(mm, m)
	}

	return mm, nil
}

func MustCompileGitIgnoreLines(lines ...string) MultiMatcher {
	mm := make(MultiMatcher, 0, len(lines))
	for _, line := range lines {
		mm = append(mm, MustCompileGitIgnore(line))
	}

	return mm
}

func (gm *GitIgnoreMatcher) MatchString(fname string) bool {
	if gm.isAlwaysNotMatch {
		return false
	}

	fname = strings.ReplaceAll(fname, `\`, "/")
	if len(fname) > 0 && fname[0] == '/' {
		fname = fname[1:]
	}

	if gm.folderToMatch != nil && !isFileInFolder(gm.folderToMatch, fname) {
		return false
	}

	match := gm.re.MatchString(fname)
	if gm.isNegative && match {
		match = !match
	}

	return match
}

func isDir(name string) bool {
	finfo, err := os.Stat(name)
	if err != nil {
		return false
	}

	return finfo.IsDir()
}

func isFileInFolder(reFolder *regexp.Regexp, fname string) bool {
	dir, base := path.Split(path.Clean(fname))

	if reFolder.MatchString(base) {
		if isDir(fname) {
			return true
		}
	}

	elems := strings.Split(dir, "/")
	for _, elem := range elems {
		if reFolder.MatchString(elem) {
			return true
		}
	}
	return false
}
