package strutil

import (
	"regexp"
	"strings"

	"github.com/jopbrown/gobase/errors"
)

func ComplieGlob(pattern string) (*regexp.Regexp, error) {
	pattern = regexp.QuoteMeta(pattern)
	pattern = replaceGlobRule(pattern, `*`, `\*`, `.*?`)
	pattern = replaceGlobRule(pattern, `?`, `\?`, `.`)
	pattern = "^" + pattern + "$"
	return regexp.Compile(pattern)
}

func MustComplieGlob(pattern string) *regexp.Regexp {
	return errors.Must1(ComplieGlob(pattern))
}

func replaceGlobRule(expr, replace, placehold, pattern string) string {
	findIndex := 0
	sb := &strings.Builder{}
	for {
		if findIndex >= len(expr) {
			break
		}

		i := strings.Index(expr[findIndex:], placehold)
		if i < 0 {
			sb.WriteString(expr[findIndex:])
			break
		}

		sb.WriteString(expr[findIndex : findIndex+i])

		findIndex += i

		if findIndex >= 1 && expr[findIndex-1] == '\\' {
			sb.WriteString(replace)
		} else {
			sb.WriteString(pattern)
		}

		findIndex += len(placehold)

	}

	return sb.String()
}
