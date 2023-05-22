package strslice

import "strings"

type Slice []string

func From(ss []string) Slice {
	return ss
}

func FromArgs(ss ...string) Slice {
	return ss
}

func Fields(str string) Slice {
	return strings.Fields(str)
}

func Split(str, sep string) Slice {
	return strings.Split(str, sep)
}

func SplitAndTrimSpace(str, sep string) Slice {
	return Slice(strings.Split(str, sep)).Map(strings.TrimSpace)
}

func FieldsFunc(str string, f func(rune) bool) Slice {
	return strings.FieldsFunc(str, f)
}

func (ss Slice) Map(fn func(s string) string) Slice {
	newss := make(Slice, 0, len(ss))
	for _, s := range ss {
		newss = append(newss, fn(s))
	}
	return newss
}

func (ss Slice) Filter(fn func(s string) bool) Slice {
	newss := make(Slice, 0, len(ss))
	for _, s := range ss {
		if fn(s) {
			newss = append(newss, s)
		}
	}
	return newss
}

func (ss Slice) Reject(fn func(s string) bool) Slice {
	newss := make(Slice, 0, len(ss))
	for _, s := range ss {
		if !fn(s) {
			newss = append(newss, s)
		}
	}
	return newss
}

func (ss Slice) Join(sep string) string {
	return strings.Join(ss, sep)
}

func (ss Slice) ToSlice() []string {
	return []string(ss)
}
