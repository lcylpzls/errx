// Package logx 提供 errx 与 github.com/lcylpzls/logx 的适配层：
// 将结构化错误转换为 logx 字段组，随日志输出错误码、分类与上下文。
package logx

import (
	"github.com/lcylpzls/errx"
	core "github.com/lcylpzls/logx"
)

// Fields 将错误转换为 logx 字段组：
// 固定输出 err.code / err.kind，以及结构化错误携带的 KV 字段与消息。
//
// 用法：
//
//	logger.Error("下单失败", errxlogx.Fields(err))
func Fields(err error) core.FieldGroup {
	code, _ := errx.CodeOf(err)
	kind := errx.KindOf(err)

	fs := []core.Field{
		core.String("err.code", string(code)),
		core.String("err.kind", kind.String()),
		core.Bool("err.retryable", errx.Retryable(err)),
	}
	if desc := errx.Describe(code); desc != "" {
		fs = append(fs, core.String("err.code_desc", desc))
	}

	if e, ok := errx.As(err); ok {
		if msg := e.Message(); msg != "" {
			fs = append(fs, core.String("err.message", msg))
		}
		for _, kv := range e.Fields() {
			fs = append(fs, core.Any(kv.Key, kv.Value))
		}
	}

	return core.Fields(fs...)
}
