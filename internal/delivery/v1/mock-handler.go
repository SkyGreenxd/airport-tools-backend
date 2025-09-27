package v1

import (
	"airport-tools-backend/internal/usecase"
	"airport-tools-backend/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) mockCheckout(c *gin.Context) {
	var req MockCheckReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorToHttpRes(e.ErrInvalidRequestBody, c)
		return
	}

	res, err := h.service.MockCheckout(c.Request.Context(), &usecase.MockCheckReq{
		EmployeeId: req.EmployeeId,
		ImageId:    req.ImageId,
		ImageUrl:   req.ImageUrl,
	})

	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, ToDeliveryCheckRes(res))
}

func (h *Handler) mockCheckin(c *gin.Context) {
	var req MockCheckReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorToHttpRes(e.ErrInvalidRequestBody, c)
		return
	}

	res, err := h.service.MockCheckin(c.Request.Context(), &usecase.MockCheckReq{
		EmployeeId: req.EmployeeId,
		ImageId:    req.ImageId,
		ImageUrl:   req.ImageUrl,
	})
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, ToDeliveryCheckRes(res))
}
