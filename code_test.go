package errx

import (
	"testing"
)

func TestRegisterCode(t *testing.T) {
	RegisterCode("USER_NOT_FOUND", "用户不存在")
	if got := Describe("USER_NOT_FOUND"); got != "用户不存在" {
		t.Errorf("Describe 不符：%s", got)
	}

	// 重复注册覆盖
	RegisterCode("USER_NOT_FOUND", "用户不存在或已注销")
	if got := Describe("USER_NOT_FOUND"); got != "用户不存在或已注销" {
		t.Errorf("重复注册覆盖不符：%s", got)
	}
}

func TestRegisterCodeEmpty(t *testing.T) {
	RegisterCode("", "空码")
	// 不应 panic，且不影响已有注册表
	if got := Describe("UNKNOWN"); got != "未知错误" {
		t.Errorf("默认错误码说明被破坏：%s", got)
	}
}

func TestDescribeUnknown(t *testing.T) {
	if got := Describe("NOT_REGISTERED"); got != "" {
		t.Errorf("未注册错误码应返回空：%s", got)
	}
}

func TestCodesSorted(t *testing.T) {
	RegisterCode("B_CODE", "B")
	RegisterCode("A_CODE", "A")
	codes := Codes()
	if len(codes) < 3 {
		t.Fatalf("Codes 数量不足：%d", len(codes))
	}
	for i := 1; i < len(codes); i++ {
		if codes[i-1].Code >= codes[i].Code {
			t.Errorf("Codes 未排序：%v", codes)
		}
	}
}
