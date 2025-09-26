package delivery

import (
	"airport-tools-backend/internal/usecase"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *usecase.Service
}

func NewHandler(service *usecase.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	transaction := api.Group("/transaction")
	{
		transaction.POST("/checkout", h.checkout) // выдача инструментов инженеру
		transaction.POST("/checkin", h.checkin)   // сдача инструментов инженером
	}
}

func (h *Handler) checkout(c *gin.Context) {
	var req CheckReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	res, err := h.service.Checkout(c.Request.Context(), ToUseCaseCheckReq(&req))
	if err != nil {
		ErrorToHttpRes(err, c)
	}

	c.JSON(http.StatusOK, ToDeliveryCheckRes(res))
}

func (h *Handler) checkin(c *gin.Context) {
	var req CheckReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	res, err := h.service.Checkin(c.Request.Context(), ToUseCaseCheckReq(&req))
	if err != nil {
		ErrorToHttpRes(err, c)
	}

	c.JSON(http.StatusOK, ToDeliveryCheckRes(res))
}
