// package errorchain is the chained error.
// For there have error and errors, we name the package errorchain.

package errorchain

import (
	"fmt"
	"runtime"
	"strings"
)

// MAX_RECURSE_LEVEL 最多迭代追溯Source次数.
var MAX_RECURSE_LEVEL = 8

//
func SetMaxRecurse(n int) {
	MAX_RECURSE_LEVEL = n
}

// Error is a type implements the error interface.
// Source: source of the error, returned by the calling other functions.
// Code: unique error code for each error.
// Msg: additional information about the error.
// Pkg: the package where error occurred.
// Func: the function where error occurred.
type Error struct {
	Source error
	Code   uint32
	Msg    string
	Pkg    string
	Func   string
}

// New return a Error with given params.
func New(source error, code uint32, msg string, pkg string, funcName string) *Error {
	return &Error{
		Source: source,
		Msg:    msg,
		Code:   code,
		Pkg:    pkg,
		Func:   funcName,
	}
}

// NewUtil is the same with New, but use runtime to find the Pkg and Func.
// runtime.Caller may failed (and I don't know when it will fail).
// So use NewUtil as your risk.
// TODO: When runtime.Caller failed?
func NewUtil(source error, code uint32, msg string) *Error {
	e := &Error{
		Source: source,
		Msg:    msg,
		Code:   code,
	}
	if counter, _, _, ok := runtime.Caller(1); ok {
		e.Pkg = runtime.FuncForPC(counter).Name()
	}
	// If failed, the Pkg and Func are nil.
	return e
}

func (e *Error) Error() string {
	return e.String()
}

func (e *Error) Cause() error {
	return e.Source
}

func (e *Error) PrettyString() string {
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
	sb.WriteString(fmt.Sprintf("{code: 0x%08x, msg: \"%s\", in: \"", e.Code, e.Msg))
	sb.WriteString(fmt.Sprintf("{code: 0x%08x, msg: \"%s\", in: \"", e.Code, e.Msg))
	if len(e.Pkg) > 0 {
		sb.WriteString(e.Pkg)
	}
	if len(e.Func) > 0 {
		sb.WriteString(".")
		sb.WriteString(e.Func)
	}
	sb.WriteString("\"}")
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
