package delivery

import (
	"airport-tools-backend/pkg/e"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorToHttpRes(err error, c *gin.Context) {
	log.Println(err)

	var code int
	var message string

	switch {
	case errors.Is(err, e.ErrUserNotFound):
		code = http.StatusNotFound
		message = "user not found"
	case errors.Is(err, e.ErrToolSetNotFound):
		code = http.StatusNotFound
		message = "tool set not found"
	case errors.Is(err, e.ErrTransactionNotFound):
		code = http.StatusNotFound
		message = "transaction not found"
	default:
		code = http.StatusInternalServerError
		message = "internal server error"
	}

	c.JSON(code, gin.H{"error": message})
}
