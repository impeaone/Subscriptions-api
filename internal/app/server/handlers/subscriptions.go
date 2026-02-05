package handlers

import (
	"agrigation_api/pkg/logger/logger"
	"agrigation_api/pkg/models"
	"agrigation_api/pkg/tools"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

// GetSubscription - GET конкретной подписки: GET /subscriptions?user_id=xxx&service=yyy
// GetSubscription godoc
// @Summary Get a specific subscription
// @Description Get subscription by user ID and service name
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string true "User ID (UUID)" example(60601fee-2bf1-4721-ae6f-7636e79a0cba)
// @Param service_name query string true "Service name" example(Netflix)
// @Success 200 {object} models.Subscription
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/subscriptions [get]
func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	// TODO: logger надо
	if r.Method != "GET" {
		h.logs.Warning(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: user uses not allowed method",
			r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
		tools.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query()
	userIDStr := query.Get("user_id")
	serviceName := query.Get("service")

	if userIDStr == "" || serviceName == "" {
		h.logs.Warning(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: user request without needed params",
			r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
		tools.WriteError(w, http.StatusBadRequest, "user_id and service parameters are required")
		return
	}

	userID, err := tools.ParseUUID(userIDStr)
	if err != nil {
		h.logs.Warning(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: user request with invalid userID",
			r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
		tools.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	subscription, err := h.db.GetSubscription(r.Context(), userID, serviceName)
	if err != nil {
		h.logs.Warning(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: %v",
			r.RemoteAddr, r.URL, r.Method, logger.TimeFormat, err), logger.GetPlace())
		tools.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if subscription == nil {
		h.logs.Info(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: subscription not found",
			r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
		tools.WriteError(w, http.StatusNotFound, "Subscription not found")
		return
	}
	h.logs.Info(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: subscription found successfully",
		r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
	tools.WriteJSON(w, http.StatusOK, subscription)
}

// UpsertSubscription - CREATE/UPDATE: POST /subscriptions
// Создает новую подписку или обновляет существующую
// UpsertSubscription godoc
// @Summary Create or update a subscription
// @Description Create new subscription or update existing one (upsert operation)
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body models.CreateOrUpdateRequest true "Subscription data"
// @Success 200 {object} models.Subscription
// @Success 201 {object} models.Subscription
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/subscriptions [post]
func (h *Handler) UpsertSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.logs.Warning(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: user uses not allowed method",
			r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
		tools.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.CreateOrUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logs.Warning(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: user request with invalid json",
			r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
		tools.WriteError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	// Валидация
	if req.ServiceName == "" {
		h.logs.Warning(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: user request with invalid service-name",
			r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
		tools.WriteError(w, http.StatusBadRequest, "service_name is required")
		return
	}
	if req.Price <= 0 {
		h.logs.Warning(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: user request with invalid price",
			r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
		tools.WriteError(w, http.StatusBadRequest, "price must be positive")
		return
	}
	if req.StartDate == "" {
		h.logs.Warning(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: user request with invalid start-date",
			r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
		tools.WriteError(w, http.StatusBadRequest, "start_date is required")
		return
	}

	subscription, err := h.db.UpsertSubscription(r.Context(), req)
	if err != nil {
		h.logs.Error(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: upsert subscription error: %v",
			r.RemoteAddr, r.URL, r.Method, logger.TimeFormat, err), logger.GetPlace())
		tools.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.logs.Info(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: subscription upsert successfully",
		r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
	tools.WriteJSON(w, http.StatusOK, subscription)
}

// DeleteSubscription godoc
// @Summary Delete a subscription
// @Description Delete subscription by user ID and service name
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body models.DeleteRequest true "Subscription identification"
// @Success 204 "Subscription deleted successfully"
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/subscriptions [delete]
func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		h.logs.Warning(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: user uses not allowed method",
			r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
		tools.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.DeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logs.Warning(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: user request with invalid json",
			r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
		tools.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Валидация
	if req.ServiceName == "" {
		h.logs.Warning(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: user uses request with invalid service-name",
			r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
		tools.WriteError(w, http.StatusBadRequest, "service_name is required")
		return
	}

	err := h.db.DeleteSubscription(r.Context(), req.UserID, req.ServiceName)
	if err != nil {
		h.logs.Error(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: delete subscription error: %v",
			r.RemoteAddr, r.URL, r.Method, logger.TimeFormat, err), logger.GetPlace())
		tools.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.logs.Info(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: subscription delete successfully",
		r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
	tools.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Subscription deleted successfully",
	})
}

// ListUserSubscriptions godoc
// @Summary List all subscriptions for a user
// @Description Get all subscriptions for a specific user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)" example(60601fee-2bf1-4721-ae6f-7636e79a0cba)
// @Success 200 {array} models.Subscription
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/subscriptions/user/{id} [get]
func (h *Handler) ListUserSubscriptions(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		h.logs.Warning(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: user uses not allowed method",
			r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
		tools.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Извлекаем user_id из пути
	path := r.PathValue("id")
	userID, err := tools.ParseUUID(path)
	if err != nil {
		h.logs.Warning(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: user request with invalid userID",
			r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
		tools.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	subscriptions, err := h.db.ListUserSubscriptions(r.Context(), userID)
	if err != nil {
		h.logs.Error(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: list subscriptions error",
			r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
		tools.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response := map[string]interface{}{
		"user_id":       userID,
		"subscriptions": subscriptions,
		"currency":      "RUB",
	}
	h.logs.Info(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: User list subscription found successfully",
		r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
	tools.WriteJSON(w, http.StatusOK, response)
}

// CalculateTotalHandler - GET /subscriptions/total
// CalculateTotalHandler godoc
// @Summary Calculate total cost for a period
// @Description Calculate total cost of subscriptions for a given period with optional filters
// @Tags analytics
// @Accept json
// @Produce json
// @Param start_month query string true "Start month (MM-YYYY or YYYY-MM)" example(01-2024)
// @Param end_month query string true "End month (MM-YYYY or YYYY-MM)" example(12-2024)
// @Param user_id query string false "User ID for filtering (UUID)" example(60601fee-2bf1-4721-ae6f-7636e79a0cba)
// @Param service_name query string false "Service name for filtering" example(Netflix)
// @Success 200 {object} models.CalculateTotalResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/subscriptions/total [get]
func (h *Handler) CalculateTotalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		h.logs.Warning(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: user uses not allowed method",
			r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
		tools.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query()

	// Парсим параметры из query string
	req := models.CalculateTotalRequest{
		StartMonth:  query.Get("start_month"),
		EndMonth:    query.Get("end_month"),
		ServiceName: query.Get("service_name"),
	}

	// user_id из query параметра
	if userIDStr := query.Get("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			h.logs.Warning(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: user request with invalid userID",
				r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
			tools.WriteError(w, http.StatusBadRequest, "Invalid user_id")
			return
		}
		req.UserID = userID
	}

	// Валидация обязательных полей
	if req.StartMonth == "" {
		h.logs.Warning(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: user request with invalid start_month",
			r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
		tools.WriteError(w, http.StatusBadRequest, "start_month is required")
		return
	}

	if req.EndMonth == "" {
		h.logs.Warning(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: user request with invalid end_month",
			r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
		tools.WriteError(w, http.StatusBadRequest, "end_month is required")
		return
	}

	// Подсчет суммы
	total, err := h.db.CalculateTotal(r.Context(), req)
	if err != nil {
		h.logs.Error(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: calculate total error: %v",
			r.RemoteAddr, r.URL, r.Method, logger.TimeFormat, err), logger.GetPlace())
		tools.WriteError(w, http.StatusBadRequest, "calculate_total error")
		return
	}

	// Формируем ответ
	response := models.CalculateTotalResponse{
		Success:  true,
		Total:    total,
		Currency: "RUB",
		Period: models.PeriodInfo{
			StartMonth: req.StartMonth,
			EndMonth:   req.EndMonth,
		},
	}

	// Добавляем фильтры, если они были указаны
	if req.UserID != uuid.Nil {
		userIDStr := req.UserID.String()
		response.Filters.UserID = &userIDStr
	}

	if req.ServiceName != "" {
		response.Filters.ServiceName = &req.ServiceName
	}

	h.logs.Info(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v; Message: calculate total successfully",
		r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())
	tools.WriteJSON(w, http.StatusOK, response)
}
