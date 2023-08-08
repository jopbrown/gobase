package errors_test

import (
	"io"
	"os"

	"github.com/jopbrown/gobase/errors"
)

func ExampleJoin() {
	err0 := io.EOF
	err1 := errors.Error("err1")
	err2 := errors.Error("err2")
	err3 := errors.Error("err3")
	err4 := errors.Error("err4")
	err5 := errors.Error("err5")
	err6 := os.ErrNotExist
	err7 := os.ErrClosed
	err8 := errors.Error("err8")
	err := errors.Join(err0, err1, err2)
	err = errors.Join(err, err3)
	err456 := errors.Join(err4, err5, err6)
	err = errors.Join(err, err456, err7, err8)

	trimErr(err)

	// Output:
	// 1. EOF
	// 2. err1
	// 	* err1
	// 		* github.com/jopbrown/gobase/errors/multi_test.go:12 errors_test.ExampleJoin
	// 3. err2
	// 	* err2
	// 		* github.com/jopbrown/gobase/errors/multi_test.go:13 errors_test.ExampleJoin
	// 4. err3
	// 	* err3
	// 		* github.com/jopbrown/gobase/errors/multi_test.go:14 errors_test.ExampleJoin
	// 5. err4
	// err5
	// file does not exist
	// 	1. err4
	// 		* err4
	// 			* github.com/jopbrown/gobase/errors/multi_test.go:15 errors_test.ExampleJoin
	// 	2. err5
	// 		* err5
	// 			* github.com/jopbrown/gobase/errors/multi_test.go:16 errors_test.ExampleJoin
	// 	3. file does not exist
	// 6. file already closed
	// 7. err8
	// 	* err8
	// 		* github.com/jopbrown/gobase/errors/multi_test.go:19 errors_test.ExampleJoin
}
