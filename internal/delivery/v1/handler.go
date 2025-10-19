package v1

import (
	"airport-tools-backend/internal/usecase"
	"airport-tools-backend/pkg/e"
	"net/http"
	"strconv"
	"time"

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
			qa.GET("/statistics", h.getStatistics)
		}
	}
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

// getStatistics
//
//	@Summary		Получить детализированную статистику QA
//	@Description	Возвращает гибкую статистику по качеству проверок и ошибкам QA-системы. Поддерживает несколько режимов работы, задаваемых параметром `type`.<br/>**Типы статистики (`type`):**<br/>- `users` — Рейтинг инженеров, чьи транзакции чаще всего попадали на QA по причине `HUMAN_ERR`.<br/>- `users&employee_id=...` — Список транзакций конкретного пользователя. Используйте `start_date`, `end_date` и `limit`, чтобы уточнить выборку.<br/>- `qa` — список всех QA-сотрудников, выполняющих проверки.<br/>- `qa&employee_id=...` — статистика проверок, проведённых конкретным QA-инженером.<br/>- `errors` — сводная статистика ошибок **Model vs Human**.<br/>- `transactions` — Выводит статистику по транзакциям (кол-во всех транзакций, а также кол-во открытых, закрытых, QA и неудачных) <br/>- `work_duration` - возвращает список всех закрытых транзакций с рассчитанной длительностью работы (Ставится значение true) <br/>- `avg_work_duration` - "Получить среднее время работы каждого инженера по всем его транзакциям (Ставится значение true)"<br/> avg_work_duration и work_duration **не совместимы*, выбирайте что то одно*.
//
//	@Tags			statistics
//	@Accept			json
//	@Produce		json
//
//	@Param			type				query		string			true	"Тип статистики. Варианты: users, qs, errors, transactions"
//	@Param			employee_id			query		string			false	"Табельный номер пользователя (для type=qa и для type=users)"
//	@Param			start_date			query		string			false	"Начало периода (формат DD-MM-YYYY, используется с type=users&employee_id=...)"
//	@Param			end_date			query		string			false	"Конец периода (формат DD-MM-YYYY, используется с type=users&employee_id=...)"
//	@Param			limit				query		int				false	"Максимальное количество записей в ответе (топ-N). По умолчанию — без ограничений (для type=users&employee_id=...)"
//	@Param			avg_work_duration	query		string			false	"Получить среднее время работы каждого инженера по всем его транзакциям (Ставится значение true)"
//	@Param			work_duration		query		string			false	"Получить список всех закрытых транзакций с рассчитанной длительностью работы (Ставится значение true)"
//
//	@Success		200					{object}	StatisticsRes	"Успешный ответ: структура зависит от значения параметра type"
//	@Failure		400					{object}	HTTPError		"Неверное тело запроса"
//	@Failure		500					{object}	HTTPError		"Внутренняя ошибка сервера"
//
//	@Router			/api/v1/qa/statistics [get]
func (h *Handler) getStatistics(c *gin.Context) {
	statisticsType := c.Query("type")
	if statisticsType == "" {
		ErrorToHttpRes(e.ErrRequestNoStatisticsType, c)
		return
	}

	employeeId := c.Query("employee_id") // табельный номер
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	limitStr := c.Query("limit")
	avgWorkDuration := c.Query("avg_work_duration")
	workDuration := c.Query("work_duration")

	valid := func(v string) bool {
		return v == "" || v == "true"
	}

	if !valid(avgWorkDuration) || !valid(workDuration) {
		ErrorToHttpRes(e.ErrInvalidRequestBody, c)
		return
	}

	var startDate, endDate *time.Time
	if startDateStr != "" {
		t, err := time.Parse("02-01-2006", startDateStr)
		if err != nil {
			ErrorToHttpRes(e.ErrInvalidRequestBody, c)
			return
		}
		startDate = &t
	}

	if endDateStr != "" {
		t, err := time.Parse("02-01-2006", endDateStr)
		if err != nil {
			ErrorToHttpRes(e.ErrInvalidRequestBody, c)
			return
		}
		endDate = &t
	}

	var limit *int
	if limitStr != "" {
		n, err := strconv.Atoi(limitStr)
		if err != nil || n <= 0 {
			ErrorToHttpRes(e.ErrInvalidRequestBody, c)
			return
		}
		limit = &n
	}

	var res interface{}
	switch statisticsType {
	case "users":
		if employeeId != "" {
			userReq := usecase.NewUserTransactionsReq(employeeId, startDate, endDate, limit)
			result, err := h.service.UserTransactions(c.Request.Context(), userReq)
			if err != nil {
				ErrorToHttpRes(err, c)
				return
			}

			res = toDeliveryListTransactionsRes(result.Transactions)
		} else {
			if workDuration != "" {
				result, err := h.service.GetWorkDuration(c.Request.Context())
				if err != nil {
					ErrorToHttpRes(err, c)
					return
				}

				res = toDeliveryGetAllWorkDurationRes(result)
			} else if avgWorkDuration != "" {
				result, err := h.service.GetAvgWorkDuration(c.Request.Context())
				if err != nil {
					ErrorToHttpRes(err, c)
					return
				}

				res = toDeliveryGetAvgWorkDurationRes(result)
			} else {
				result, err := h.service.GetUsersQAStats(c.Request.Context())
				if err != nil {
					ErrorToHttpRes(err, c)
					return
				}

				res = toArrDeliveryHumanErrorStats(result)
			}
		}
	case "errors":
		result, err := h.service.GetMlVsHuman(c.Request.Context())
		if err != nil {
			ErrorToHttpRes(err, c)
			return
		}

		res = toDeliveryModelOrHumanStatsRes(result)

	case "qa":
		if employeeId != "" {
			result, err := h.service.GetQAChecks(c.Request.Context(), employeeId)
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
	case "transactions":
		result, err := h.service.GetTransactionStatistics(c.Request.Context())
		if err != nil {
			ErrorToHttpRes(err, c)
			return
		}

		res = toDeliveryGetTransactionStatisticsRes(*result)
	default:
		ErrorToHttpRes(e.ErrRequestNoStatisticsType, c)
		return
	}

	c.JSON(http.StatusOK, NewStatisticsRes(statisticsType, res))
}
