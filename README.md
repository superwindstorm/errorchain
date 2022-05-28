# errorchain

Usage:
1. 对每个package，定义
    const PACKAGE_NAME = "gomain.com/gcl/errorchain"
或者使用反射和init:
    var PACKAGE_NAME string
	func init() {
	    type DUMMY struct{}
	    PACKAGE_NAME = reflect.TypeOf(DUMMY{}).PkgPath()
    }
2. 调用其它库得到error，包装为Error返回:
	e := a_error_occure()
    return &errorchain.Error{
	    Code:    0x00000001,
	    Source:  e,
	    Package: PACKAGE_NAME,
	    Func:    "function_name",
 	    Info:    "",
    }
3. 自己得到error，返回Error:
	invalid input happens
    return &errorchain.Error{
	    Code:    0x00000002,
	    Source:  nil,
	    Package: PACKAGE_NAME,
	    Func:    "function_name",
 	    Info:    "invalid input",
    }