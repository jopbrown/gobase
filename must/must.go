package must

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Value[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}

func Error[T any](v T, err error) error {
	return err
}

func Value2[T, U any](v T, u U, err error) (T, U) {
	if err != nil {
		panic(err)
	}

	return v, u
}

func Error2[T, U any](v T, u U, err error) error {
	return err
}
