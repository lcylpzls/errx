package errx

import (
	"strings"
	"testing"
)

func TestMarkdown(t *testing.T) {
	RegisterCode("MD_A", "错误 A")
	RegisterCode("MD_B", "错误 B")

	out := Markdown()
	if !strings.Contains(out, "# 错误码表") {
		t.Error("Markdown 缺少标题")
	}
	if !strings.Contains(out, "| 错误码 | 说明 |") {
		t.Error("Markdown 缺少表头")
	}
	if !strings.Contains(out, "| MD_A | 错误 A |") {
		t.Error("Markdown 缺少已注册错误码")
	}
	// 排序：MD_A 应在 MD_B 之前
	if strings.Index(out, "MD_A") > strings.Index(out, "MD_B") {
		t.Error("Markdown 错误码未按序排列")
	}
}
