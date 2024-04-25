package output

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"slow/util/types"
)

func Success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "",
		"data": data,
	})
}

func Fail(c echo.Context, err error) error {
	code := -1
	msg := "system error"
	var customErr types.Error
	if ok := errors.As(err, &customErr); ok {
		code = customErr.Code
		msg = customErr.Msg
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": struct{}{},
	})
}
