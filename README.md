# errorchain

Usage:
1. 对每个package，定义
```
const PACKAGE_NAME = "github.com/superwindstorm/errorchain"
```
或者使用反射和init:
```
var PACKAGE_NAME string
func init() {
    type DUMMY struct{}
    PACKAGE_NAME = reflect.TypeOf(DUMMY{}).PkgPath()
}
```
2. 调用其它库得到error，包装为Error返回:
```
e := a_error_occure()
return &errorchain.Error{
    Code:    0x00000001,
    Source:  e,
    Package: PACKAGE_NAME,
    Func:    "function_name",
    Info:    "",
}
```
3. 自己得到error，返回Error:
```
invalid input happens
return &errorchain.Error{
    Code:    0x00000002,
    Source:  nil,
    Package: PACKAGE_NAME,
    Func:    "function_name",
    Info:    "invalid input",
}
```
4. 输出示例
```
{code: 0x00000003, msg: "no such file error", in: "github.com/superwindstorm/errorchain_test.TestError2"}, caused by
	{code: 0x00000002, msg: "", in: "github.com/superwindstorm/errorchain_test.TestError2"}, caused by
	{code: 0x00000001, msg: "invalid input", in: "github.com/superwindstorm/errorchain_test.errorOccureFromStd"}, caused by
	error caused by stdlib.
```