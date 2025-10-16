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

// ErrorToHttpRes формирует HTTP-ответ на основе переданной ошибки

// TODO: заменить ошибки на русские
func ErrorToHttpRes(err error, c *gin.Context) {
	log.Println(err)

	var res HTTPError

	switch {
	case errors.Is(err, e.ErrUserNotFound):
		res.Code = http.StatusNotFound
		res.Message = "user not found"
	case errors.Is(err, e.ErrToolSetNotFound):
		res.Code = http.StatusNotFound
		res.Message = "assigned tool set not found"
	case errors.Is(err, e.ErrTransactionNotFound):
		res.Code = http.StatusNotFound
		res.Message = "transaction not found"
	case errors.Is(err, e.ErrTransactionUnfinished):
		res.Code = http.StatusBadRequest
		res.Message = "return previous tool set before taking a new one"
	case errors.Is(err, e.ErrInvalidRequestBody):
		res.Code = http.StatusBadRequest
		res.Message = "invalid request body"
	case errors.Is(err, e.ErrTransactionAllFinished):
		res.Code = http.StatusBadRequest
		res.Message = "no active tool set to return"
	case errors.Is(err, e.ErrTransactionLimit):
		res.Code = http.StatusConflict
		res.Message = "3 неудачные попытки сканирования. Данные отправлены на проверку QA"
	case errors.Is(err, e.ErrTransactionCheckQA):
		res.Code = http.StatusConflict
		res.Message = "Вы не можете получить новые инструменты, пока вас проверяет QA"
	case errors.Is(err, e.ErrUserExists):
		res.Code = http.StatusConflict
		res.Message = "Пользователь с таким табельным номером уже существует"
	case errors.Is(err, e.ErrRequestNotSupported):
		res.Code = http.StatusBadRequest
		res.Message = "Некорректное значение параметра 'status'. Допустимое значение: 'qa'"
	case errors.Is(err, e.ErrUserRoleNotFound):
		res.Code = http.StatusNotFound
		res.Message = "Такой роли не существует"
	default:
		res.Code = http.StatusInternalServerError
		res.Message = "internal server error"
	}

	c.JSON(res.Code, res)
}
