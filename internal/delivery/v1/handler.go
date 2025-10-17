package v1

import (
	"airport-tools-backend/internal/usecase"
	"airport-tools-backend/pkg/e"
	"net/http"
	"strconv"

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
			transaction.POST("/check", h.check)                                  // выдача/сдача инструментов
			transaction.POST("/:trasaction_id/verification", h.postVerification) // отправка qa результата
			transaction.GET("/:trasaction_id/verification", h.getVerification)   // получение данных для qa
			transaction.GET("/", h.list)                                         // list проблемных проверок
		}

		user := v1.Group("/user")
		{
			user.POST("/login", h.login)       // вход в систему по табельному номеру
			user.GET("/roles", h.getRoles)     // получить список ролей
			user.POST("/register", h.register) // регистрация в системе
		}
	}
}

// check
//
//	@Summary		Операция выдачи/сдачи инструментов
//	@Description	Принимает табельный номер инженера и фотографию инструментов в формате base64.<br> Сервис анализирует изображение, сопоставляет инструменты с ожидаемым набором и возвращает: <br><br>• URL обработанного изображения <br>• четыре массива: <br>1) AccessTools — инструменты, прошедшие автоматическую проверку <br>2) ManualCheckTools — инструменты, требующие ручной проверки <br>3) UnknownTools — инструменты, отсутствующие в ожидаемом наборе <br>4) MissingTools — инструменты, отсутствующие на фотографии, но ожидаемые<br><br> Если не все инструменты попали в AccessTools, устанавливается флаг "ТРЕБУЕТСЯ РУЧНАЯ ПРОВЕРКА" (MANUAL VERIFICATION). <br><br>Эндпоинт используется как для выдачи инструментов инженеру, так и для их последующей сдачи.
//
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CheckReq	true	"Запрос на выдачу или сдачу инструментов"
//	@Success		200		{object}	CheckRes
//	@Failure		400		{object}	HTTPError
//	@Failure		404		{object}	HTTPError
//	@Failure		500		{object}	HTTPError
//	@Router			/api/v1/transaction/check [post]
func (h *Handler) check(c *gin.Context) {
	var req CheckReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorToHttpRes(e.ErrInvalidRequestBody, c)
		return
	}

	res, err := h.service.Check(c.Request.Context(), ToUseCaseCheckReq(&req))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, ToDeliveryCheckRes(res))
}

// postVerification
//
//	@Summary		QA-проверка и завершение транзакции
//	@Description	После авторизации сотрудника QA отображается список всех незавершённых транзакций.<br> При выборе конкретной транзакции открывается экран сверки:<br><br> • Фотография инструментов (полноразмерное изображение)<br> • Список проблемных инструментов с пояснениями, сгруппированных по категориям:<br> &nbsp;&nbsp;1) AccessTools — инструменты, прошедшие автоматическую проверку<br> &nbsp;&nbsp;2) ManualCheckTools — инструменты, требующие ручной проверки<br> &nbsp;&nbsp;3) UnknownTools — инструменты, не входящие в ожидаемый набор<br> &nbsp;&nbsp;4) MissingTools — инструменты, отсутствующие на фото, но ожидаемые<br><br>
//
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			transaction_id	path		string		true	"Идентификатор транзакции"
//	@Param			request			body		VerificationReq	true	"Данные завершения QA-проверки"
//	@Success		200				{object}	VerificationRes
//	@Failure		400				{object}	HTTPError
//	@Failure		404				{object}	HTTPError
//	@Failure		500				{object}	HTTPError
//	@Router			/api/v1/transaction/{transaction_id}/verification [post]
func (h *Handler) postVerification(c *gin.Context) {
	strTransactionId := c.Param("transaction_id")
	transactionId, err := strconv.Atoi(strTransactionId)
	if err != nil {
		ErrorToHttpRes(e.ErrInvalidRequestBody, c)
		return
	}

	var req VerificationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorToHttpRes(e.ErrInvalidRequestBody, c)
		return
	}

	res, err := h.service.Verification(c.Request.Context(), usecase.NewVerification(int64(transactionId), req.QAEmployeeId, req.QAEmployeeId))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, toDeliveryVerificationRes(res))
}

// получение инфы по конкретной транзакции
func (h *Handler) getVerification(c *gin.Context) {
	strTransactionId := c.Param("transaction_id")
	transactionId, err := strconv.Atoi(strTransactionId)
	if err != nil {
		ErrorToHttpRes(e.ErrInvalidRequestBody, c)
		return
	}

	res, err := h.service.GetQATransaction(c.Request.Context(), int64(transactionId))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, toDeliveryGetQAVerificationRes(res))
}

// list
//
//	@Summary		Список транзакций
//	@Description	Возвращает список транзакций QA.<br> Можно фильтровать по статусу с помощью query-параметра `status`.<br> Допустимое значение: 'qa' вернёт только транзакции, требующие проверки QA.<br> Каждая транзакция содержит минимальные данные: ID, инженера, дату создания, текущий статус.
//
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			status	query		string	false	"Фильтр по статусу транзакции"
//	@Success		200		{object}	ListTransactionsRes
//	@Failure		400		{object}	HTTPError
//	@Failure		500		{object}	HTTPError
//	@Router			/api/v1/transaction [get]
func (h *Handler) list(c *gin.Context) {
	status := c.Query("status")
	res, err := h.service.List(c.Request.Context(), status)
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, toDeliveryListTransactionsRes(res.Transactions))
}

// login
//
//	@Summary		Вход в систему
//	@Description	Вход в систему по табельному номеру сотрудника.<br> После успешного входа пользователь перенаправляется:<br> • инженеру — на экран загрузки фотографии инструментов;<br> • QA — на экран проверки незавершённых транзакций.
//
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoginReq	true	"Данные для входа"
//	@Success		200		{object}	LoginRes
//	@Failure		400		{object}	HTTPError
//	@Failure		401		{object}	HTTPError
//	@Failure		500		{object}	HTTPError
//	@Router			/api/v1/user/login [post]
func (h *Handler) login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorToHttpRes(e.ErrInvalidRequestBody, c)
		return
	}

	res, err := h.service.Login(c.Request.Context(), toUseCaseLoginReq(req))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, toDeliveryLoginRes(res))
}

// register
//
//	@Summary		Регистрация сотрудника в системе
//	@Description	Регистрация сотрудника в системе.<br> Необходимые данные: табельный номер, ФИО и роль (например, "Инженер" или "QA").<br>
//
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RegisterReq	true	"Данные для регистрации"
//	@Success		201		{object}	RegisterRes
//	@Failure		400		{object}	HTTPError
//	@Failure		409		{object}	HTTPError
//	@Failure		500		{object}	HTTPError
//	@Router			/api/v1/user/register [post]
func (h *Handler) register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorToHttpRes(e.ErrInvalidRequestBody, c)
		return
	}

	res, err := h.service.Register(c.Request.Context(), toUseCaseRegisterReq(req))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusCreated, toDeliveryRegisterRes(res))
}

// getRoles
//
//	@Summary		Получить список ролей
//	@Description	Возвращает список всех возможных ролей пользователей в системе.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		GetRolesRes	"Список ролей"
//	@Failure		500	{object}	HTTPError	"Внутренняя ошибка сервера"
//	@Router			/api/v1/user/roles [get]
func (h *Handler) getRoles(c *gin.Context) {
	res, err := h.service.GetRoles(c.Request.Context())
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, toDeliveryGetRolesRes(res))
}
