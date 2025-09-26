package v1

import (
	"airport-tools-backend/pkg/e"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ErrorToHttpRes(err error, c *gin.Context) {
	log.Println(err)

	var res HTTPError

	switch {
	case errors.Is(err, e.ErrUserNotFound):
		res.Code = http.StatusNotFound
		res.Message = "user not found"
	case errors.Is(err, e.ErrToolSetNotFound):
		res.Code = http.StatusNotFound
		res.Message = "tool set not found"
	case errors.Is(err, e.ErrTransactionNotFound):
		res.Code = http.StatusNotFound
		res.Message = "transaction not found"
	case errors.Is(err, e.ErrTransactionUnfinished):
		res.Code = http.StatusBadRequest
		res.Message = "you have an unfinished issue"
	case errors.Is(err, e.ErrInvalidRequestBody):
		res.Code = http.StatusBadRequest
		res.Message = "invalid request body"
	default:
		res.Code = http.StatusInternalServerError
		res.Message = "internal server error"
	}

	c.JSON(res.Code, res)
}
