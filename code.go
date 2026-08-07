package errx

import (
	"sort"
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
}

var (
	codeMu sync.RWMutex
	codes  = map[Code]string{
		CodeUnknown: "未知错误",
	}
)

// RegisterCode 注册错误码及其说明。重复注册以最后一次为准。
// 建议在包 init 或程序启动阶段完成注册。
func RegisterCode(code Code, description string) {
	if code == "" {
		return
	}
	codeMu.Lock()
	codes[code] = description
	codeMu.Unlock()
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
		out = append(out, CodeInfo{Code: code, Description: desc})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Code < out[j].Code })
	return out
}
