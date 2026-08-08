package httpx_test

import (
	"fmt"
	"net/http/httptest"

	"github.com/lcylpzls/errx"
	"github.com/lcylpzls/errx/httpx"
)

func ExampleWriteJSON() {
	rec := httptest.NewRecorder()
	httpx.WriteJSON(rec, errx.New(errx.KindNotFound, "USER_NOT_FOUND", "用户不存在"))
	fmt.Println(rec.Code)
	fmt.Print(rec.Body.String())
	// Output:
	// 404
	// {"code":"USER_NOT_FOUND","kind":"not_found","message":"USER_NOT_FOUND: 用户不存在"}
}
