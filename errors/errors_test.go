package errors_test

import (
	"fmt"
	"io"

	"github.com/jopbrown/gobase/errors"
)

func showErr(err error) {
	// fmt.Print(err)
	// fmt.Printf("%s", err)
	// fmt.Printf("%q", err)
	// fmt.Printf("%v", err)
	fmt.Printf("%+v", err)
	// fmt.Printf("%+2v", err)
}

func ExampleError() {
	err := errors.Error("the error:", "unable to got resource")
	showErr(err)

	// Output:
	// * the error:unable to got resource
	// 	* github.com/jopbrown/gobase/errors/errors_test.go:20 errors_test.ExampleError
}

func ExampleErrorOmit() {
	err := errors.Error()
	showErr(err)

	// Output:
	// * something is wrong
	// 	* github.com/jopbrown/gobase/errors/errors_test.go:29 errors_test.ExampleErrorOmit
}

func ExampleErrorf() {
	err := errors.Errorf("the error:%s", "unable to got resource")
	showErr(err)

	// Output:
	// * the error:unable to got resource
	// 	* github.com/jopbrown/gobase/errors/errors_test.go:38 errors_test.ExampleErrorf
}

func ExampleErrorAt() {
	err := errors.ErrorAt(io.EOF, "the error:", "unable to got resource")
	showErr(err)

	// Output:
	// * EOF
	// * the error:unable to got resource
	// 	* github.com/jopbrown/gobase/errors/errors_test.go:47 errors_test.ExampleErrorAt
}

func ExampleErrorAt2() {
	err := errors.ErrorAt(io.EOF, "1st layer error")
	err = errors.ErrorAt(err, "2nd layer error")
	err = errors.ErrorAt(err, "3rd layer error")
	err = errors.ErrorAt(err, "4th layer error")
	showErr(err)

	// Output:
	// * EOF
	// * 1st layer error
	// 	* github.com/jopbrown/gobase/errors/errors_test.go:57 errors_test.ExampleErrorAt2
	// * 2nd layer error
	// 	* github.com/jopbrown/gobase/errors/errors_test.go:58 errors_test.ExampleErrorAt2
	// * 3rd layer error
	// 	* github.com/jopbrown/gobase/errors/errors_test.go:59 errors_test.ExampleErrorAt2
	// * 4th layer error
	// 	* github.com/jopbrown/gobase/errors/errors_test.go:60 errors_test.ExampleErrorAt2
}

func ExampleErrorAtOmit() {
	err := errors.ErrorAt(io.EOF)
	err = errors.ErrorAt(err)
	err = errors.ErrorAt(err)
	err = errors.ErrorAt(err)
	showErr(err)

	// Output:
	// * EOF
	// * something is wrong
	// 	* github.com/jopbrown/gobase/errors/errors_test.go:76 errors_test.ExampleErrorAtOmit
	// * something is wrong
	// 	* github.com/jopbrown/gobase/errors/errors_test.go:77 errors_test.ExampleErrorAtOmit
	// * something is wrong
	// 	* github.com/jopbrown/gobase/errors/errors_test.go:78 errors_test.ExampleErrorAtOmit
	// * something is wrong
	// 	* github.com/jopbrown/gobase/errors/errors_test.go:79 errors_test.ExampleErrorAtOmit
}

func ExampleErrorAtf() {
	err := io.EOF
	for i := 0; i < 4; i++ {
		err = errors.ErrorAtf(err, "%d err", i+1)
	}
	showErr(err)

	// Output:
	// * EOF
	// * 1 err
	// 	* github.com/jopbrown/gobase/errors/errors_test.go:97 errors_test.ExampleErrorAtf
	// * 2 err
	// 	* github.com/jopbrown/gobase/errors/errors_test.go:97 errors_test.ExampleErrorAtf
	// * 3 err
	// 	* github.com/jopbrown/gobase/errors/errors_test.go:97 errors_test.ExampleErrorAtf
	// * 4 err
	// 	* github.com/jopbrown/gobase/errors/errors_test.go:97 errors_test.ExampleErrorAtf
}
