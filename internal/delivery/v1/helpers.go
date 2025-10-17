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
func ErrorToHttpRes(err error, c *gin.Context) {
	log.Println(err)
	var res HTTPError

	switch {
	case errors.Is(err, e.ErrUserNotFound):
		res.Code = http.StatusNotFound
		res.Message = "Пользователь не найден"
	case errors.Is(err, e.ErrToolSetNotFound):
		res.Code = http.StatusNotFound
		res.Message = "Назначенный набор инструментов не найден"
	case errors.Is(err, e.ErrTransactionNotFound):
		res.Code = http.StatusNotFound
		res.Message = "Транзакция не найдена"
	case errors.Is(err, e.ErrTransactionUnfinished):
		res.Code = http.StatusBadRequest
		res.Message = "Верните предыдущий набор инструментов, прежде чем брать новый"
	case errors.Is(err, e.ErrInvalidRequestBody):
		res.Code = http.StatusBadRequest
		res.Message = "Неверное тело запроса"
	case errors.Is(err, e.ErrTransactionAllFinished):
		res.Code = http.StatusBadRequest
		res.Message = "Вы не получали инструменты, чтобы их возвращать"
	case errors.Is(err, e.ErrTransactionLimit):
		res.Code = http.StatusConflict
		res.Message = "Превышено 3 попытки сканирования. Данные переданы на проверку QA"
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
	case errors.Is(err, e.ErrIncorrectImage):
		res.Code = http.StatusBadRequest
		res.Message = "Изображение повреждено или не поддерживается"
	case errors.Is(err, e.ErrCvScanNotFound):
		res.Code = http.StatusNotFound
		res.Message = "Скан не найден"
	default:
		res.Code = http.StatusInternalServerError
		res.Message = "Внутренняя ошибка сервера"
	}

	c.JSON(res.Code, res)
}
