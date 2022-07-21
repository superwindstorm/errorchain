// package errorchain is the chained error.
// For there have error and errors, we name the package errorchain.

package errorchain

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// MAX_RECURSE_LEVEL 最多迭代追溯Source次数.
var MAX_RECURSE_LEVEL = 16

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
	source error
	code   uint32
	msg    string
	pkg    string
	fn     string
	tag    int // a wrapper error could use tag to distinguish between wrapper errors.
}

// New return a Error with given params.
func New(source error, code uint32, msg string, pkg string, funcName string) *Error {
	return &Error{
		source: source,
		msg:    msg,
		code:   code,
		pkg:    pkg,
		fn:     funcName,
	}
}

// NewUtil is the same with New, but use runtime to find the Pkg and Func.
// runtime.Caller may failed (and I don't know when it will fail).
// So use NewUtil as your risk.
// TODO: When runtime.Caller failed?
func NewUtil(source error, code uint32, msg string) *Error {
	e := &Error{
		source: source,
		msg:    msg,
		code:   code,
	}
	if counter, _, _, ok := runtime.Caller(1); ok {
		e.pkg = runtime.FuncForPC(counter).Name()
	}
	// If failed, the Pkg and Func are nil.
	return e
}

// Just a wrapper around source, use 1,2,3 as tags(not 0).
func Wrapper(source error, tag int) *Error {
	e := &Error{
		source: source,
		tag:    tag,
	}
	if counter, _, _, ok := runtime.Caller(1); ok {
		e.pkg = runtime.FuncForPC(counter).Name()
	}
	return e
}

func (e *Error) Code() uint32 {
	return e.code
}

func (e *Error) Error() string {
	return e.String()
}

func (e *Error) Cause() error {
	return e.source
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
	sb.WriteString(fmt.Sprintf("{code: 0x%08x, msg: \"%s\", in: \"", e.code, e.msg))
	if len(e.pkg) > 0 {
		sb.WriteString(e.pkg)
	}
	if len(e.fn) > 0 {
		sb.WriteString(".")
		sb.WriteString(e.fn)
	}
	if e.tag > 0 {
		sb.WriteString(":")
		sb.WriteString(strconv.Itoa(e.tag))
	}
	sb.WriteString("\"}")
	if e.source != nil {
		sb.WriteString(", caused by")
		if lf {
			sb.WriteString("\n\t")
		} else {
			sb.WriteString(" ")
		}
		if ee, ok := e.source.(*Error); ok {
			ee.display(sb, lf, n+1)
		} else {
			sb.WriteString(e.source.Error())
		}
	}
}
