package errx

import (
	"fmt"
	"strings"
)

// Markdown 生成错误码注册表的 Markdown 文档（按错误码排序），
// 可直接交付前端、API 网关或审计使用。
func Markdown() string {
	infos := Codes()
	var b strings.Builder
	b.WriteString("# 错误码表\n\n")
	b.WriteString("| 错误码 | 说明 |\n")
	b.WriteString("| --- | --- |\n")
	for _, info := range infos {
		fmt.Fprintf(&b, "| %s | %s |\n", info.Code, info.Description)
	}
	return b.String()
}
