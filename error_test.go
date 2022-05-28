package errorchain_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/superwindstorm/errorchain"
)

const PACKAGE_NAME = "github.com/superwindstorm/errorchain"

// or use init function to auto generate a PACKAGE_NAME.
// var PACKAGE_NAME string

// func init() {
// 	type DUMMY struct{}
// 	PACKAGE_NAME = reflect.TypeOf(DUMMY{}).PkgPath()
// 	fmt.Println(PACKAGE_NAME)
// }

func errorOccures() *errorchain.Error {
	return &errorchain.Error{
		Code:   0x00000001,
		Source: nil,
		Pkg:    PACKAGE_NAME,
		Func:   "errorOccures",
		Msg:    "invalid input",
	}
}

func TestError1(t *testing.T) {
	e := errorOccures()
	e3 := &errorchain.Error{
		Code:   0x00000002,
		Source: e,
		Pkg:    PACKAGE_NAME,
		Func:   "TestError1",
	}
	t.Log(e3)
}

func errorOccureFromStd() *errorchain.Error {
	e1 := errors.New("error caused by stdlib.")
	return &errorchain.Error{
		Code:   0x00000001,
		Source: e1,
		Pkg:    PACKAGE_NAME,
		Func:   "errorOccureFromStd",
		Msg:    "invalid input",
	}
}

func TestError2(t *testing.T) {
	e := errorOccureFromStd()
	e3 := &errorchain.Error{
		Code:   0x00000002,
		Source: e,
		Pkg:    PACKAGE_NAME,
		Func:   "TestError2",
		// leave info empty
	}
	e4 := errorchain.New(e3, 0x00000003, PACKAGE_NAME, "TestError2", "")
	t.Log(e3)
	t.Log(e4)

	fmt.Println(e3.PretyString())
	fmt.Println(e4.PretyString())
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
	fmt.Println(e1.PretyString())
	fmt.Println()
	fmt.Println(e2.PretyString())
}
