package v1

import (
	"airport-tools-backend/internal/usecase"
	"airport-tools-backend/pkg/e"
	"net/http"

	_ "airport-tools-backend/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	v1 := api.Group("/v1")
	{
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		transaction := v1.Group("/transaction")
		{
			transaction.POST("/checkout", h.checkout) // выдача инструментов инженеру
			transaction.POST("/checkin", h.checkin)   // сдача инструментов инженером
		}
	}
}

// checkout
//
//	@Summary		Instrument issuance operation
//	@Description	Receives the engineer's personnel number and a photo of the tools in base64 format,<br>verifies them against the expected set, and returns image URL and four arrays:<br>1) AccessTools: tools that passed automated checks<br>2) ManualCheckTools: tools requiring manual verification<br>3) UnknownTools: tools not in the expected set<br>4) MissingTools: tools missing from the expected set<br>If not all tools are in Access, a MANUAL VERIFICATION flag is set.
//
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CheckReq	true	"Checkout request"
//	@Success		200		{object}	CheckRes
//	@Failure		400		{object}	HTTPError
//	@Failure		404		{object}	HTTPError
//	@Failure		500		{object}	HTTPError
//	@Router			/api/v1/transaction/checkout [post]
func (h *Handler) checkout(c *gin.Context) {
	var req CheckReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorToHttpRes(e.ErrInvalidRequestBody, c)
		return
	}

	res, err := h.service.Checkout(c.Request.Context(), ToUseCaseCheckReq(&req))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, ToDeliveryCheckRes(res))
}

// checkin
//
//	@Summary		Instrument delivery operation
//	@Description 	Receives the engineer's personnel number and a photo of the tools in base64 format, verifies them against the expected set,<br>and returns image URL and four arrays:<br>1) AccessTools - tools that passed automated checks<br>2) ManualCheckTools - tools requiring manual verification<br>3) UnknownTools - tools not in the expected set<br>4) MissingTools - tools missing from the expected set<br>If not all tools are in Access, a MANUAL VERIFICATION flag is set.
//
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CheckReq	true	"Checkin request"
//	@Success		200		{object}	CheckRes
//	@Failure		400		{object}	HTTPError
//	@Failure		404		{object}	HTTPError
//	@Failure		500		{object}	HTTPError
//	@Router			/api/v1/transaction/checkin [post]
func (h *Handler) checkin(c *gin.Context) {
	var req CheckReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorToHttpRes(e.ErrInvalidRequestBody, c)
		return
	}

	res, err := h.service.Checkin(c.Request.Context(), ToUseCaseCheckReq(&req))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, ToDeliveryCheckRes(res))
}
