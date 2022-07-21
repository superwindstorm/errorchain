package errorchain

import (
	"errors"
	"fmt"
	"testing"
)

const PACKAGE_NAME = "github.com/superwindstorm/errorchain_test"

// or use init function to auto generate a PACKAGE_NAME.
// var PACKAGE_NAME string

// func init() {
// 	type DUMMY struct{}
// 	PACKAGE_NAME = reflect.TypeOf(DUMMY{}).PkgPath()
// 	fmt.Println(PACKAGE_NAME)
// }

func errorOccures() *Error {
	// return &Error{
	// 	Code:   0x00000001,
	// 	Source: nil,
	// 	Pkg:    PACKAGE_NAME,
	// 	Func:   "errorOccures",
	// 	Msg:    "invalid input",
	// }
	return NewUtil(nil, 0x00000001, "invalid input")
}

func TestError1(t *testing.T) {
	e := NewUtil(errorOccures(), 0x00000002, "oops, some error occurred")
	e2 := Wrapper(e, 1)

	fmt.Println(e2.PrettyString())

}

func errorOccureFromStd() *Error {
	return NewUtil(errors.New("error caused by stdlib."), 0x00000001, "invalid input")
}

func TestError2(t *testing.T) {
	e := errorOccureFromStd()
	// 手动填写Pkg和Func
	// e2 := &Error{
	// 	Code:   0x00000002,
	// 	Source: e,
	// 	Pkg:    PACKAGE_NAME,
	// 	Func:   "TestError2",
	// 	// leave info empty
	// }

	// or
	// 使用runtime自动获取。
	e2 := NewUtil(e, 0x00000002, "info")
	e3 := NewUtil(e2, 0x00000003, "no such file error")
	s := e3.PrettyString()
	fmt.Println(s)
}

// two recursive errors
func TestErrorLoop(t *testing.T) {
	SetMaxRecurse(100)
	e1 := &Error{
		code: 0x00000001,
		pkg:  PACKAGE_NAME,
		fn:   "TestErrorLoop",
		msg:  "error 1",
	}
	e2 := &Error{
		code:   0x00000002,
		source: e1,
		pkg:    PACKAGE_NAME,
		fn:     "TestErrorLoop",
		msg:    "error 2",
	}
	e1.source = e2
	// Stop recurce after 8 times.
	fmt.Println(e1.PrettyString())
	fmt.Println()
	fmt.Println(e2.PrettyString())
}
