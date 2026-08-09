package errx

import (
	"errors"
	"strings"
	"testing"
)

func TestRegisterCode(t *testing.T) {
	RegisterCode("USER_NOT_FOUND", "用户不存在")
	if got := Describe("USER_NOT_FOUND"); got != "用户不存在" {
		t.Errorf("Describe 不符：%s", got)
	}

	// 同码同描述幂等,不 panic 且不改变描述。
	RegisterCode("USER_NOT_FOUND", "用户不存在")
	if got := Describe("USER_NOT_FOUND"); got != "用户不存在" {
		t.Errorf("幂等注册不应覆盖原描述：%s", got)
	}
}

func TestRegisterCodeConflictPanic(t *testing.T) {
	code := Code("TEST_CONFLICT_ONLY")
	RegisterCode(code, "原描述")
	defer func() {
		if r := recover(); r == nil {
			t.Error("同码不同描述应 panic")
		} else if !strings.Contains(r.(string), string(code)) {
			t.Errorf("panic 信息应包含错误码:%v", r)
		}
	}()
	RegisterCode(code, "新描述")
}

func TestRegisterCodeKindAndCodeKind(t *testing.T) {
	RegisterCodeKind("KINDED_CODE", KindTimeout)
	if got := CodeKind("KINDED_CODE"); got != KindTimeout {
		t.Errorf("CodeKind = %s,want timeout", got)
	}
	RegisterCodeKind("", KindInvalid) // 空码忽略
	if got := CodeKind("NOT_REGISTERED_KIND"); got != KindUnknown {
		t.Errorf("未注册分类应为 unknown:%s", got)
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
	for _, c := range codes {
		if c.Code == "KINDED_CODE" && c.Kind != KindTimeout {
			t.Error("Codes 应携带 Kind 元数据")
		}
	}
}

func TestNewCode(t *testing.T) {
	RegisterCode("NEW_CODE", "新码描述")
	RegisterCodeKind("NEW_CODE", KindUnavailable)
	e := NewCode("NEW_CODE", "服务不可用")
	if e.Code() != "NEW_CODE" || e.Kind() != KindUnavailable {
		t.Errorf("NewCode 元数据不符:%v %v", e.Code(), e.Kind())
	}
	if e.Error() != "NEW_CODE: 服务不可用" {
		t.Errorf("NewCode 消息不符:%s", e.Error())
	}
	// 未注册 Kind 的码 → KindUnknown
	RegisterCode("PLAIN_CODE", "普通码")
	e = NewCode("PLAIN_CODE", "x")
	if e.Kind() != KindUnknown {
		t.Errorf("未注册分类应为 unknown:%s", e.Kind())
	}
	// 格式化构造
	e = NewCodef("NEW_CODE", "第 %d 次失败", 3)
	if e.Error() != "NEW_CODE: 第 3 次失败" {
		t.Errorf("NewCodef 消息不符:%s", e.Error())
	}
}

func TestWrapCode(t *testing.T) {
	RegisterCode("WRAP_CODE", "包装码")
	RegisterCodeKind("WRAP_CODE", KindTimeout)
	cause := New(KindBusiness, "CAUSE", "原因")
	e := WrapCode(cause, "WRAP_CODE", "包装消息")
	if e == nil || e.Code() != "WRAP_CODE" || e.Kind() != KindTimeout {
		t.Errorf("WrapCode 元数据不符:%v %v", e.Code(), e.Kind())
	}
	if !errors.Is(e, cause) {
		t.Error("WrapCode 应保留原因链")
	}
	// nil 原因返回 nil
	if got := WrapCode(nil, "WRAP_CODE", "x"); got != nil {
		t.Error("nil 原因应返回 nil")
	}
	if got := WrapCodef(nil, "WRAP_CODE", "x"); got != nil {
		t.Error("WrapCodef 的 nil 原因应返回 nil")
	}
	// 格式化
	e = WrapCodef(cause, "WRAP_CODE", "第 %d 次", 2)
	if e.Error() != "WRAP_CODE: 第 2 次: CAUSE: 原因" {
		t.Errorf("WrapCodef 消息不符:%s", e.Error())
	}
}

func TestCodesMarkdown(t *testing.T) {
	RegisterCode("DOC_TEST_CODE", "文档测试码")
	RegisterCodeKind("DOC_TEST_CODE", KindConflict)
	RegisterCode("DOC_TEST_AA", "文档测试码二") // 同组两个码,覆盖组内排序
	md := CodesMarkdown()
	if !strings.Contains(md, "# 错误码手册") {
		t.Error("缺少标题")
	}
	if !strings.Contains(md, "DOC_TEST_CODE") {
		t.Error("缺少测试错误码")
	}
	if !strings.Contains(md, "DOC_TEST_AA") {
		t.Error("缺少同组第二个错误码")
	}
	if !strings.Contains(md, "conflict") {
		t.Error("缺少 Kind 分类")
	}
	if !strings.Contains(md, "| 错误码 | 分类 | 说明 |") {
		t.Error("缺少表头")
	}
	if !strings.Contains(md, "## DOC") {
		t.Error("缺少分组")
	}
}

func TestGroupPrefix(t *testing.T) {
	cases := []struct {
		code Code
		want string
	}{
		{"DBX_OPEN_FAILED", "DBX"},
		{"HTX_INVALID_CONFIG", "HTX"},
		{"UNKNOWN", "UNKN"},
		{"AB", "AB"},
	}
	for _, c := range cases {
		if got := groupPrefix(c.code); got != c.want {
			t.Errorf("groupPrefix(%q) = %q,want %q", c.code, got, c.want)
		}
	}
}
