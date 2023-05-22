package strutil

type Matcher interface {
	MatchString(str string) bool
}

type MatchHandler func(str string) bool

func (h MatchHandler) MatchString(str string) bool {
	return h(str)
}

type MultiMatcher []Matcher

func MakeMultiMatcher(matchers ...Matcher) MultiMatcher {
	return MultiMatcher(matchers)
}

func (mm MultiMatcher) All() Matcher {
	return MatchHandler(func(str string) bool {
		for _, m := range mm {
			if !m.MatchString(str) {
				return false
			}
		}
		return true
	})
}

func (mm MultiMatcher) Any() Matcher {
	return MatchHandler(func(str string) bool {
		for _, m := range mm {
			if m.MatchString(str) {
				return true
			}
		}
		return false
	})
}

func (mm MultiMatcher) None() Matcher {
	return MatchHandler(func(str string) bool {
		for _, m := range mm {
			if m.MatchString(str) {
				return false
			}
		}
		return true
	})
}
