// Package httpx 提供 errx 与标准库 net/http 的适配层：
// 将结构化错误映射为 HTTP 状态码与 JSON 响应体。
package httpx

import (
	"encoding/json"
	"net/http"

	"github.com/lcylpzls/errx"
)

// Status 返回错误对应的 HTTP 状态码；无结构化错误时返回 500。
func Status(err error) int {
	if e, ok := errx.As(err); ok {
		return e.HTTPStatus()
	}
	return http.StatusInternalServerError
}

// response 是统一的错误响应体。
type response struct {
	Code    string `json:"code"`
	Kind    string `json:"kind"`
	Message string `json:"message"`
}

// WriteJSON 将错误以 JSON 形式写入 ResponseWriter：
// 设置 Content-Type、状态码（由错误分类映射），响应体含 code/kind/message。
// err 为 nil 时输出 500 与未知分类，保证上层可直接安全调用。
func WriteJSON(w http.ResponseWriter, err error) {
	code, ok := errx.CodeOf(err)
	if !ok || code == "" {
		code = errx.CodeUnknown
	}
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	resp := response{
		Code:    string(code),
		Kind:    errx.KindOf(err).String(),
		Message: msg,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(Status(err))
	_ = json.NewEncoder(w).Encode(resp)
}
