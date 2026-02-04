package handlers

import (
	"agrigation_api/pkg/models"
	"agrigation_api/pkg/tools"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

// GetSubscription - GET конкретной подписки: GET /subscriptions?user_id=xxx&service=yyy
func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	// TODO: logger надо
	if r.Method != "GET" {
		tools.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query()
	userIDStr := query.Get("user_id")
	serviceName := query.Get("service")

	if userIDStr == "" || serviceName == "" {
		tools.WriteError(w, http.StatusBadRequest, "user_id and service parameters are required")
		return
	}

	userID, err := tools.ParseUUID(userIDStr)
	if err != nil {
		tools.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	subscription, err := h.db.GetSubscription(r.Context(), userID, serviceName)
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if subscription == nil {
		tools.WriteError(w, http.StatusNotFound, "Subscription not found")
		return
	}

	tools.WriteJSON(w, http.StatusOK, subscription)
}

// UpsertSubscription - CREATE/UPDATE: POST /subscriptions
// Создает новую подписку или обновляет существующую
func (h *Handler) UpsertSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		tools.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.CreateOrUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		tools.WriteError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	// Валидация
	if req.ServiceName == "" {
		tools.WriteError(w, http.StatusBadRequest, "service_name is required")
		return
	}
	if req.Price <= 0 {
		tools.WriteError(w, http.StatusBadRequest, "price must be positive")
		return
	}
	if req.StartDate == "" {
		tools.WriteError(w, http.StatusBadRequest, "start_date is required")
		return
	}

	subscription, err := h.db.UpsertSubscription(r.Context(), req)
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tools.WriteJSON(w, http.StatusOK, subscription)
}

func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		tools.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.DeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		tools.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Валидация
	if req.ServiceName == "" {
		tools.WriteError(w, http.StatusBadRequest, "service_name is required")
		return
	}

	err := h.db.DeleteSubscription(r.Context(), req.UserID, req.ServiceName)
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tools.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Subscription deleted successfully",
	})
}

func (h *Handler) ListUserSubscriptions(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		tools.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Извлекаем user_id из пути
	path := r.PathValue("id")
	userID, err := tools.ParseUUID(path)
	if err != nil {
		tools.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	subscriptions, err := h.db.ListUserSubscriptions(r.Context(), userID)
	if err != nil {
		tools.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"user_id":       userID,
		"subscriptions": subscriptions,
		"currency":      "RUB",
	}

	tools.WriteJSON(w, http.StatusOK, response)
}

// CalculateTotalHandler - GET /subscriptions/total
func (h *Handler) CalculateTotalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
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
			tools.WriteError(w, http.StatusBadRequest, "Invalid user_id")
			return
		}
		req.UserID = userID
	}

	// Валидация обязательных полей
	if req.StartMonth == "" {
		tools.WriteError(w, http.StatusBadRequest, "start_month is required")
		return
	}

	if req.EndMonth == "" {
		tools.WriteError(w, http.StatusBadRequest, "end_month is required")
		return
	}

	// Подсчет суммы
	total, err := h.db.CalculateTotal(r.Context(), req)
	if err != nil {
		tools.WriteError(w, http.StatusBadRequest, err.Error())
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

	tools.WriteJSON(w, http.StatusOK, response)
}
