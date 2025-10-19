package v1

import (
	_ "airport-tools-backend/docs"
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/internal/usecase"
	"airport-tools-backend/pkg/e"
	"airport-tools-backend/pkg/parse"
	"net/http"
	"strconv"

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

		// AUTH
		auth := v1.Group("/auth")
		{
			auth.POST("/login", h.login)
			auth.POST("/register", h.register)
		}

		//  USER
		user := v1.Group("/users")
		{
			user.GET("/roles", h.getRoles)
			user.POST("/check", h.check) // выдача/сдача инструментов пользователем
		}

		// QA
		qa := v1.Group("/qa")
		{
			transactions := qa.Group("/transactions")
			{
				transactions.GET("/", h.list)                                          // список всех проблемных транзакций
				transactions.GET("/:transaction_id", h.getVerification)                // получение данных для QA
				transactions.POST("/:transaction_id/verification", h.postVerification) // отправка QA результата
			}

			// Аналитика QA
			statisticsGroup := qa.Group("/statistics")
			{
				statisticsGroup.GET("/users", h.getUserStatistics)               // Для ?type=users
				statisticsGroup.GET("/errors", h.getErrorStatistics)             // Для ?type=errors
				statisticsGroup.GET("/qa", h.getQaStatistics)                    // Для ?type=qa
				statisticsGroup.GET("/transactions", h.getTransactionStatistics) // Для ?type=transactions
			}
		}
	}
}

// getUserStatistics
//
//	@Summary		Получить статистику пользователей (инженеров)
//
//	@Description	Возвращает статистику по всем инженерам или конкретному сотруднику. Поддерживает:<br/>- `employee_id` — список транзакций конкретного пользователя (можно фильтровать по дате, лимиту транзакций, добавить среднее время работы);<br/>- `avg_work_duration=true` — среднее время работы каждого инженера;<br/>- `start_date/end_date` — начало и конец периода транзакций;<br/>- `limit` — кол-во транзакций на вывод;<br/>- Без параметров — список всех транзакций всех инженеров.
//	@Tags			statistics
//	@Produce		json
//
//	@Param			employee_id			query		string			false	"Табельный номер инженера"
//	@Param			start_date			query		string			false	"Начало периода (формат DD-MM-YYYY)"
//	@Param			end_date			query		string			false	"Конец периода (формат DD-MM-YYYY)"
//	@Param			limit				query		int				false	"Максимальное количество записей для вывода"
//	@Param			avg_work_duration	query		bool			false	"true — получить среднее время работы каждого инженера"
//	@Success		200					{object}	StatisticsRes	"Успешный ответ"
//	@Failure		400					{object}	HTTPError		"Неверные параметры"
//	@Failure		500					{object}	HTTPError		"Ошибка сервера"
//	@Router			/api/v1/qa/statistics/users [get]
func (h *Handler) getUserStatistics(c *gin.Context) {
	flags, err := parse.ParseCommonFilters(c)
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	var res interface{}
	if flags.EmployeeId != nil && *flags.EmployeeId != "" {
		userReq := usecase.NewUserTransactionsReq(*flags.EmployeeId, flags.StartDate, flags.EndDate, flags.Limit, flags.AvgWorkDuration)
		result, err := h.service.UserTransactions(c.Request.Context(), userReq)
		if err != nil {
			ErrorToHttpRes(err, c)
			return
		}

		res = toDeliveryGetUsersListTransactionsRes(result)
	} else if (flags.EmployeeId == nil || *flags.EmployeeId == "") && flags.AvgWorkDuration == true {
		result, err := h.service.GetAvgWorkDuration(c.Request.Context())
		if err != nil {
			ErrorToHttpRes(err, c)
			return
		}

		res = toDeliveryGetAvgWorkDurationRes(result)
	} else {
		result, err := h.service.GetAllTransactions(c.Request.Context())
		if err != nil {
			ErrorToHttpRes(err, c)
			return
		}

		resultArr := make([]*GetAllTransactions, len(result))
		for i, item := range result {
			resultArr[i] = toDeliveryGetAllTransactions(item)
		}

		res = resultArr
	}

	c.JSON(http.StatusOK, res)
}

// getErrorStatistics
//
//	@Summary		Получить статистику ошибок
//
//	@Description	Возвращает статистику ошибок системы и QA. Поддерживает:<br/>- `error_type=MODEL_ERR` — список транзакций, где ошиблась ML-модель;<br/>- `error_type=HUMAN_ERR` — статистика ошибок QA-инженеров;<br/>- Без параметров — общее сравнение ML vs Human ошибок.
//
//	@Tags			statistics
//	@Produce		json
//
//	@Param			error_type	query		string			false	"Тип ошибки: MODEL_ERR или HUMAN_ERR"
//	@Success		200			{object}	StatisticsRes	"Успешный ответ"
//	@Failure		400			{object}	HTTPError		"Неверные параметры"
//	@Failure		500			{object}	HTTPError		"Ошибка сервера"
//	@Router			/api/v1/qa/statistics/errors [get]
func (h *Handler) getErrorStatistics(c *gin.Context) {
	flags, err := parse.ParseCommonFilters(c)
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	var res interface{}
	if flags.ErrorType != nil && *flags.ErrorType == string(domain.ModelError) {
		result, err := h.service.GetMlErrorTransactions(c.Request.Context())
		if err != nil {
			ErrorToHttpRes(err, c)
			return
		}

		res = toArrDeliveryMlErrorTransaction(result)
	} else if flags.ErrorType != nil && *flags.ErrorType == string(domain.HumanError) {
		result, err := h.service.GetUsersQAStats(c.Request.Context())
		if err != nil {
			ErrorToHttpRes(err, c)
			return
		}

		res = toArrDeliveryHumanErrorStats(result)
	} else {
		result, err := h.service.GetMlVsHuman(c.Request.Context())
		if err != nil {
			ErrorToHttpRes(err, c)
			return
		}

		res = toDeliveryModelOrHumanStatsRes(result)
	}

	c.JSON(http.StatusOK, res)
}

// getQaStatistics
//
//	@Summary		Получить статистику QA
//
//	@Description	Возвращает список QA-сотрудников или статистику конкретного QA-инженера.<br/>Поддерживает:<br/>- `employee_id` — статистика проверок конкретного QA-инженера;<br/>- Без параметров — список всех QA-сотрудников, выполняющих проверки.
//
//	@Tags			statistics
//	@Produce		json
//
//	@Param			employee_id	query		string			false	"Табельный номер QA-инженера"
//	@Success		200			{object}	StatisticsRes	"Успешный ответ"
//	@Failure		400			{object}	HTTPError		"Неверные параметры"
//	@Failure		500			{object}	HTTPError		"Ошибка сервера"
//	@Router			/api/v1/qa/statistics/qa [get]
func (h *Handler) getQaStatistics(c *gin.Context) {
	flags, err := parse.ParseCommonFilters(c)
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	var res interface{}
	if flags.EmployeeId != nil && *flags.EmployeeId != "" {
		result, err := h.service.GetQAChecks(c.Request.Context(), *flags.EmployeeId)
		if err != nil {
			ErrorToHttpRes(err, c)
			return
		}

		res = toDeliveryQaTransactionsRes(result)
	} else {
		result, err := h.service.GetAllQaEmployers(c.Request.Context())
		if err != nil {
			ErrorToHttpRes(err, c)
			return
		}

		res = toArrDeliveryUserDto(result)
	}
	c.JSON(http.StatusOK, res)
}

// getTransactionStatistics
//
//	@Summary		Получить общую статистику транзакций
//	@Description	Возвращает агрегированную статистику по всем транзакциям:<br/>- общее количество;<br/>- количество QA-транзакций;<br/>- количество открытых/закрытых транзакций;<br/>- количество неудачных транзакций.
//	@Tags			statistics
//	@Produce		json
//
//	@Success		200	{object}	StatisticsRes	"Успешный ответ"
//	@Failure		400	{object}	HTTPError		"Неверные параметры"
//	@Failure		500	{object}	HTTPError		"Ошибка сервера"
//	@Router			/api/v1/qa/statistics/transactions [get]
func (h *Handler) getTransactionStatistics(c *gin.Context) {
	var res interface{}
	result, err := h.service.GetTransactionStatistics(c.Request.Context())
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	res = toDeliveryGetTransactionStatisticsRes(*result)

	c.JSON(http.StatusOK, res)
}

// check
//
//	@Summary		Операция выдачи/сдачи инструментов
//	@Description	Принимает табельный номер инженера и фотографию инструментов в формате base64.<br> Сервис анализирует изображение, сопоставляет инструменты с ожидаемым набором и возвращает: <br><br>• URL обработанного изображения <br>• четыре массива: <br>1) access_tools — инструменты, прошедшие автоматическую проверку<br>1) manual_check_tools — инструменты, требующие ручной проверки <br>2) unknown_tools — инструменты, отсутствующие в ожидаемом наборе <br>3) missing_tools — инструменты, отсутствующие на фотографии, но ожидаемые<br>• transaction_type - тип транзакции(Checkin - Сдача/Checkout - Выдача)<br>• status - статус транзакции(OPEN - открыта, CLOSED - закрыта, QA VERIFICATION - QA проверка)<br><br> Если 4 или более инструментов не попали в access_tools или за 3 попытки сканирования транзакция не закрылась, устанавливается флаг "QA ПРОВЕРКА" (QA VERIFICATION). <br><br>Эндпоинт используется как для выдачи инструментов инженеру, так и для их последующей сдачи.
//
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CheckReq	true	"Запрос на выдачу или сдачу инструментов"
//	@Success		200		{object}	CheckRes	"Успешная проверка"
//	@Failure		400		{object}	HTTPError	"Неверное тело запроса"
//	@Failure		500		{object}	HTTPError	"Внутренняя ошибка сервера"
//	@Router			/api/v1/users/check [post]
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
//	@Description	После авторизации сотрудника QA выбирает из списка транзакцию.<br>Открывается экран сверки:<br><br> • Фотография инструментов (полноразмерное изображение)<br> • access_tools — инструменты, прошедшие автоматическую проверку<br> • Список проблемных инструментов с пояснениями, сгруппированных по категориям:<br> &nbsp;&nbsp;1) manual_check_tools — инструменты, требующие ручной проверки<br> &nbsp;&nbsp;2) unknown_tools — инструменты, не входящие в ожидаемый набор<br> &nbsp;&nbsp;3) missing_tools — инструменты, отсутствующие на фото, но ожидаемые<br><br>
//
//	@Tags			QA
//	@Accept			json
//	@Produce		json
//	@Param			transaction_id	path		string			true	"Идентификатор транзакции"
//	@Param			request			body		VerificationReq	true	"Данные завершения QA-проверки"
//	@Success		200				{object}	VerificationRes	"Успешное закрытие транзакции"
//	@Failure		400				{object}	HTTPError		"Неверное тело запроса"
//	@Failure		404				{object}	HTTPError		"Транзакция не найдена"
//	@Failure		500				{object}	HTTPError		"Внутренняя ошибка сервера"
//	@Router			/api/v1/qa/transactions/:transaction_id/verification [post]
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

	res, err := h.service.Verification(c.Request.Context(), usecase.NewVerification(int64(transactionId), req.QAEmployeeId, req.Reason, req.Notes))
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, toDeliveryVerificationRes(res))
}

// getVerification
//
//	@Summary		Получение информации о транзакции
//	@Description	Получить информацию о проблемной транзакции.<br>Открывается экран сверки:<br><br> • Фотография инструментов (полноразмерное изображение)<br> • access_tools — инструменты, прошедшие автоматическую проверку<br> • Список проблемных инструментов с пояснениями, сгруппированных по категориям:<br> &nbsp;&nbsp;2) manual_check_tools — инструменты, требующие ручной проверки<br> &nbsp;&nbsp;3) unknown_tools — инструменты, не входящие в ожидаемый набор<br> &nbsp;&nbsp;4) missing_tools — инструменты, отсутствующие на фото, но ожидаемые<br><br>
//
//	@Tags			QA
//	@Accept			json
//	@Produce		json
//	@Param			transaction_id	path		string					true	"Идентификатор транзакции"
//	@Param			request			body		GetQAVerificationRes	true	"Данные о транзакции"
//	@Success		200				{object}	VerificationRes			"Информация о транзакции"
//	@Failure		400				{object}	HTTPError				"Неверное тело запроса"
//	@Failure		404				{object}	HTTPError				"Транзакция не найдена"
//	@Failure		500				{object}	HTTPError				"Внутренняя ошибка сервера"
//	@Router			/api/v1/qa/transactions/:transaction_id [get]
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
//	@Description	Возвращает список транзакций QA.<br> Можно фильтровать по статусу с помощью query-параметра `status`.<br> Допустимое значение: 'qa' вернёт только транзакции, требующие проверки QA.<br> Каждая транзакция содержит минимальные данные: ID, инженера, номер набора инструментов, дату создания транзакции, текущий статус.
//
//	@Tags			QA
//	@Accept			json
//	@Produce		json
//	@Param			status	query		string				false	"Фильтр по статусу транзакции"
//	@Success		200		{object}	ListTransactionsRes	"Список транзакций"
//	@Failure		400		{object}	HTTPError			"Неверное тело запроса"
//	@Failure		500		{object}	HTTPError			"Внутренняя ошибка сервера"
//	@Router			/api/v1/qa/transactions/ [get]
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
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoginReq	true	"Данные для входа"
//	@Success		200		{object}	LoginRes	"Успешная авторизация"
//	@Failure		400		{object}	HTTPError	"Неверное тело запроса"
//	@Failure		404		{object}	HTTPError	"Пользователь не найден"
//	@Failure		500		{object}	HTTPError	"Внутренняя ошибка сервера"
//	@Router			/api/v1/auth/login [post]
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
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RegisterReq	true	"Данные для регистрации"
//	@Success		201		{object}	RegisterRes	"Регистрация успешна"
//	@Failure		400		{object}	HTTPError	"Неверное тело запроса"
//	@Failure		409		{object}	HTTPError	"Пользователь с таким табельным номером уже существует"
//	@Failure		500		{object}	HTTPError	"Внутренняя ошибка сервера"
//	@Router			/api/v1/auth/register [post]
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
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		GetRolesRes	"Список ролей"
//	@Failure		500	{object}	HTTPError	"Внутренняя ошибка сервера"
//	@Router			/api/v1/users/roles [get]
func (h *Handler) getRoles(c *gin.Context) {
	res, err := h.service.GetRoles(c.Request.Context())
	if err != nil {
		ErrorToHttpRes(err, c)
		return
	}

	c.JSON(http.StatusOK, toDeliveryGetRolesRes(res))
}

func (h *Handler) addToolSet(c *gin.Context) {
	var req AddToolSetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorToHttpRes(e.ErrInvalidRequestBody, c)
		return
	}

	res, err := h.service.AddToolSet(c.Request.Context(), req)
}
