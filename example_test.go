package errx_test

import (
	"errors"
	"fmt"

	"github.com/lcylpzls/errx"
)

func ExampleNew() {
	err := errx.New(errx.KindNotFound, "USER_NOT_FOUND", "用户不存在")
	fmt.Println(err)
	// Output:
	// USER_NOT_FOUND: 用户不存在
}

func ExampleWrap() {
	err := errx.Wrap(errors.New("连接超时"), errx.KindTimeout, "DB_TIMEOUT", "数据库超时")
	fmt.Println(err)
	// Output:
	// DB_TIMEOUT: 数据库超时: 连接超时
}

func ExampleJoin() {
	err := errx.Join(
		errx.New(errx.KindBusiness, "ORDER_FAIL", "下单失败"),
		errx.New(errx.KindInternal, "DB_DOWN", "数据库不可用"),
	)
	fmt.Println(err)
	// Output:
	// 2 个错误：
	//   1) ORDER_FAIL: 下单失败
	//   2) DB_DOWN: 数据库不可用
}
