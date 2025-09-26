package delivery

import (
	"airport-tools-backend/internal/usecase"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) mockCheckout(c *gin.Context) {
	var req MockCheckReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	res, err := h.service.MockCheckout(c.Request.Context(), &usecase.MockCheckReq{
		EmployeeId: req.EmployeeId,
		ImageId:    req.ImageId,
		ImageUrl:   req.ImageUrl,
	})

	if err != nil {
		ErrorToHttpRes(err, c)
	}

	c.JSON(http.StatusOK, ToDeliveryCheckRes(res))
}

func (h *Handler) mockCheckin(c *gin.Context) {
	var req MockCheckReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	res, err := h.service.MockCheckin(c.Request.Context(), &usecase.MockCheckReq{
		EmployeeId: req.EmployeeId,
		ImageId:    req.ImageId,
		ImageUrl:   req.ImageUrl,
	})
	if err != nil {
		ErrorToHttpRes(err, c)
	}

	c.JSON(http.StatusOK, ToDeliveryCheckRes(res))
}
