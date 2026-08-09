package errx

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// Code 是稳定的业务错误码，建议使用大写下划线风格（如 "USER_NOT_FOUND"）。
type Code string

// CodeUnknown 是未指定错误码时的默认值。
const CodeUnknown Code = "UNKNOWN"

// CodeInfo 描述一个错误码及其含义，用于文档生成与审计。
type CodeInfo struct {
	Code        Code
	Description string
	// Kind 是注册的分类，未注册时为 KindUnknown。
	Kind Kind
}

var (
	codeMu    sync.RWMutex
	codes     = map[Code]string{CodeUnknown: "未知错误"}
	codeKinds = map[Code]Kind{}
)

// RegisterCode 注册错误码及其说明。重复注册以最后一次为准。
// 同码同描述幂等；同码不同描述视为编程错误并 panic，
// 防止不同模块静默覆盖同一错误码。
// 建议在包 init 或程序启动阶段完成注册。
func RegisterCode(code Code, description string) {
	if code == "" {
		return
	}
	codeMu.Lock()
	defer codeMu.Unlock()
	if prev, ok := codes[code]; ok {
		if prev != description {
			panic(fmt.Sprintf(
				"errx: 错误码 %q 已注册为 %q,重复注册为 %q",
				code, prev, description))
		}
		return
	}
	codes[code] = description
}

// RegisterCodeKind 注册错误码的分类，供 NewCode / NewCodef 自动使用。
// 建议与 RegisterCode 成对调用；未注册分类的错误码构造时使用 KindUnknown。
func RegisterCodeKind(code Code, kind Kind) {
	if code == "" {
		return
	}
	codeMu.Lock()
	codeKinds[code] = kind
	codeMu.Unlock()
}

// CodeKind 返回错误码注册的分类；未注册返回 KindUnknown。
func CodeKind(code Code) Kind {
	codeMu.RLock()
	defer codeMu.RUnlock()
	if kind, ok := codeKinds[code]; ok {
		return kind
	}
	return KindUnknown
}

// Describe 返回错误码的注册说明；未注册时返回空字符串。
func Describe(code Code) string {
	codeMu.RLock()
	defer codeMu.RUnlock()
	return codes[code]
}

// Codes 返回全部已注册错误码的快照，按错误码排序。
func Codes() []CodeInfo {
	codeMu.RLock()
	defer codeMu.RUnlock()

	out := make([]CodeInfo, 0, len(codes))
	for code, desc := range codes {
		out = append(out, CodeInfo{
			Code:        code,
			Description: desc,
			Kind:        codeKinds[code],
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Code < out[j].Code })
	return out
}

// NewCode 基于已注册错误码构造结构化错误，自动使用注册的 Kind。
// 未注册分类的错误码使用 KindUnknown。
func NewCode(code Code, msg string) *Error {
	return New(CodeKind(code), code, msg)
}

// NewCodef 基于已注册错误码构造带格式化消息的结构化错误。
func NewCodef(code Code, format string, args ...any) *Error {
	return Newf(CodeKind(code), code, format, args...)
}

// WrapCode 包装底层错误并附加已注册错误码，自动使用注册分类。
// err 为 nil 时返回 nil，便于直接 return errx.WrapCode(err, ...) 写法。
func WrapCode(err error, code Code, msg string) *Error {
	if err == nil {
		return nil
	}
	return Wrap(err, CodeKind(code), code, msg)
}

// WrapCodef 包装底层错误并附加已注册错误码与格式化消息。
func WrapCodef(err error, code Code, format string, args ...any) *Error {
	if err == nil {
		return nil
	}
	return Wrapf(err, CodeKind(code), code, format, args...)
}

// CodesMarkdown 生成全库错误码手册 Markdown：按前缀分组，
// 每行包含错误码、分类与说明，可直接写入仓库文档。
func CodesMarkdown() string {
	infos := Codes()
	var b strings.Builder
	b.WriteString("# 错误码手册\n\n")
	b.WriteString("> 由 errx 注册表自动生成,新增错误码请通过 `RegisterCode` 注册。\n\n")

	// 按前缀分组(以 _ 前部分或前 4 字符为组)。
	groups := map[string][]CodeInfo{}
	var groupKeys []string
	for _, info := range infos {
		key := groupPrefix(info.Code)
		if _, ok := groups[key]; !ok {
			groupKeys = append(groupKeys, key)
		}
		groups[key] = append(groups[key], info)
	}
	sort.Strings(groupKeys)
	for _, key := range groupKeys {
		items := groups[key]
		sort.Slice(items, func(i, j int) bool { return items[i].Code < items[j].Code })
		fmt.Fprintf(&b, "## %s\n\n", key)
		b.WriteString("| 错误码 | 分类 | 说明 |\n")
		b.WriteString("| --- | --- | --- |\n")
		for _, item := range items {
			fmt.Fprintf(&b, "| `%s` | %s | %s |\n",
				item.Code, item.Kind.String(), item.Description)
		}
		b.WriteString("\n")
	}
	return b.String()
}

// groupPrefix 提取错误码分组前缀：取 _ 前部分；无下划线时取前 4 字符。
func groupPrefix(code Code) string {
	s := string(code)
	if idx := strings.IndexByte(s, '_'); idx > 0 {
		return s[:idx]
	}
	if len(s) > 4 {
		return s[:4]
	}
	return s
}
