package errors

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Must1[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}

func Must2[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	if err != nil {
		panic(err)
	}

	return v1, v2
}

func Must3[T1, T2, T3 any](v1 T1, v2 T2, v3 T3, err error) (T1, T2, T3) {
	if err != nil {
		panic(err)
	}

	return v1, v2, v3
}

func Should1[T any](v T, err error) (T, bool) {
	return v, err != nil
}

func Should2[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2, bool) {
	return v1, v2, err != nil
}

func Should3[T1, T2, T3 any](v1 T1, v2 T2, v3 T3, err error) (T1, T2, T3, bool) {
	return v1, v2, v3, err != nil
}

func Has(err error) bool {
	return err != nil
}

func Has1[T any](v T, err error) bool {
	return err != nil
}

func Has2[T1, T2 any](v1 T1, v2 T2, err error) bool {
	return err != nil
}

func Has3[T1, T2, T3 any](v1 T1, v2 T2, v3 T3, err error) bool {
	return err != nil
}

func Ignore1[T any](v T, err error) T {
	return v
}

func Ignore2[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	return v1, v2
}

func Ignore3[T1, T2, T3 any](v1 T1, v2 T2, v3 T3, err error) (T1, T2, T3) {
	return v1, v2, v3
}

func Get1[T any](v T, err error) error {
	return err
}

func Get2[T1, T2 any](v1 T1, v2 T2, err error) error {
	return err
}

func Get3[T1, T2, T3 any](v1 T1, v2 T2, v3 T3, err error) error {
	return err
}

func Catch(fn func()) (err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		if e, ok := r.(error); ok {
			err = e
			return
		}

		err = Error(r)
	}()
	fn()
	return
}
