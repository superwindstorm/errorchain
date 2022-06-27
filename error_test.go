package errorchain_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/superwindstorm/errorchain"
)

const PACKAGE_NAME = "github.com/superwindstorm/errorchain_test"

// or use init function to auto generate a PACKAGE_NAME.
// var PACKAGE_NAME string

// func init() {
// 	type DUMMY struct{}
// 	PACKAGE_NAME = reflect.TypeOf(DUMMY{}).PkgPath()
// 	fmt.Println(PACKAGE_NAME)
// }

func errorOccures() *errorchain.Error {
	// return &errorchain.Error{
	// 	Code:   0x00000001,
	// 	Source: nil,
	// 	Pkg:    PACKAGE_NAME,
	// 	Func:   "errorOccures",
	// 	Msg:    "invalid input",
	// }
	return errorchain.NewUtil(nil, 0x00000001, "invalid input")
}

func TestError1(t *testing.T) {
	e := errorchain.NewUtil(errorOccures(), 0x00000002, "oops, some error occurred")
	e2 := errorchain.Wrapper(e, 1)

	fmt.Println(e2.PrettyString())

}

func errorOccureFromStd() *errorchain.Error {
	return errorchain.NewUtil(errors.New("error caused by stdlib."), 0x00000001, "invalid input")
}

func TestError2(t *testing.T) {
	e := errorOccureFromStd()
	// 手动填写Pkg和Func
	// e2 := &errorchain.Error{
	// 	Code:   0x00000002,
	// 	Source: e,
	// 	Pkg:    PACKAGE_NAME,
	// 	Func:   "TestError2",
	// 	// leave info empty
	// }

	// or
	// 使用runtime自动获取。
	e2 := errorchain.NewUtil(e, 0x00000002, "info")
	e3 := errorchain.NewUtil(e2, 0x00000003, "no such file error")
	s := e3.PrettyString()
	fmt.Println(s)
}

// two recursive errors
func TestErrorLoop(t *testing.T) {
	errorchain.SetMaxRecurse(100)
	e1 := &errorchain.Error{
		Code: 0x00000001,
		Pkg:  PACKAGE_NAME,
		Func: "TestErrorLoop",
		Msg:  "error 1",
	}
	e2 := &errorchain.Error{
		Code:   0x00000002,
		Source: e1,
		Pkg:    PACKAGE_NAME,
		Func:   "TestErrorLoop",
		Msg:    "error 2",
	}
	e1.Source = e2
	// Stop recurce after 8 times.
	fmt.Println(e1.PrettyString())
	fmt.Println()
	fmt.Println(e2.PrettyString())
}
