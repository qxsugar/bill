package router

import "github.com/qxsugar/pkg/kit"

// wrapErr 把 service 层错误转换为统一响应错误：
// 已是业务错误（kit.BusinessError）则原样透传，否则按内部错误包装。
func wrapErr(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(kit.BusinessError); ok {
		return err
	}
	return kit.NewInternalError().WithErr(err)
}
