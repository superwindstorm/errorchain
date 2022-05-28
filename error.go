// package errorchain is the chained error.
// For there have error and errors, we name the package errorchain.

package errorchain

import (
	"fmt"
	"strings"
)

// MAX_RECURSE_LEVEL 最多迭代追溯Source次数.
var MAX_RECURSE_LEVEL = 8

//
func SetMaxRecurse(n int) {
	MAX_RECURSE_LEVEL = n
}

type Error struct {
	Source error  // source of the error, returned by the calling other functions.
	Code   uint32 // unique error code for each error.
	Msg    string // additional information about the error.
	Pkg    string // the package where error occurred.
	Func   string // the function where error occurred.
}

func New(source error, code uint32, msg string, pkg string, funcName string) *Error {
	return &Error{
		Source: source,
		Msg:    msg,
		Code:   code,
		Pkg:    pkg,
		Func:   funcName,
	}
}
func (e *Error) Error() string {
	return e.String()
}

func (e *Error) Unwrap() error {
	return e.Source
}

func (e *Error) PretyString() string {
	sb := new(strings.Builder)
	e.display(sb, true, 0)
	return sb.String()
}

func (e *Error) String() string {
	sb := new(strings.Builder)
	e.display(sb, false, 0)
	return sb.String()
}

func (e *Error) display(sb *strings.Builder, lf bool, n int) {
	// 防御循环迭代stack溢出
	if n >= MAX_RECURSE_LEVEL {
		sb.WriteString("...")
		return
	}
	sb.WriteString(fmt.Sprintf("{code: 0x%08x, msg: \"%s\", in: \"%s.%s\"}", e.Code, e.Msg, e.Pkg, e.Func))
	if e.Source != nil {
		sb.WriteString(", caused by")
		if lf {
			sb.WriteString("\n\t")
		} else {
			sb.WriteString(" ")
		}
		if ee, ok := e.Source.(*Error); ok {
			ee.display(sb, lf, n+1)
		} else {
			sb.WriteString(e.Source.Error())
		}
	}
}
