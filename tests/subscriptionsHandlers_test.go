package tests

import (
	handlers2 "agrigation_api/internal/app/server/handlers"
	"agrigation_api/internal/service"
	logger2 "agrigation_api/pkg/logger"
	"agrigation_api/pkg/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubscriptionHandlers(t *testing.T) {
	testLoger := logger2.NewMyLogger("INFO")
	testRepository, err := NewTestRepository()
	testService := service.NewSubscriptionService(testRepository)
	if err != nil {
		t.Fatalf("Error initializing repository: %v", err)
	}
	handlers := handlers2.NewHandler(testService, testLoger)

	testRequest := &models.CreateOrUpdateRequest{
		ServiceName: "test_service",
		Price:       300,
		UserID:      uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba"),
		StartDate:   "07-2025",
		EndDate:     "08-2025",
	}
	jsonDataTestRequest, _ := json.Marshal(testRequest)
	///////// CreateSubscription Test //////////////////////////////////////////////////////////////////////////////////
	reqCrSub, errCrSub := http.NewRequest("POST",
		"/api/v1/subscriptions/",
		bytes.NewBuffer(jsonDataTestRequest))
	if errCrSub != nil {
		t.Fatal(errCrSub)
	}
	rCrSub := httptest.NewRecorder()

	handlers.CreateSubscription(rCrSub, reqCrSub)
	if status := rCrSub.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	///////// GetSubscription Test /////////////////////////////////////////////////////////////////////////////////////
	reqGetSub, errGetSub := http.NewRequest("GET",
		fmt.Sprintf("/api/v1/subscriptions/?user_id=%s&service_name=%s", testRequest.UserID, testRequest.ServiceName),
		nil)
	if errGetSub != nil {
		t.Fatal(errGetSub)
	}
	rGetSub := httptest.NewRecorder()

	handlers.GetSubscription(rGetSub, reqGetSub)
	if status := rGetSub.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	///////// UpdateSubscription Test //////////////////////////////////////////////////////////////////////////////////
	reqPutSub, errPutSub := http.NewRequest("PUT", "/api/v1/subscriptions/", bytes.NewBuffer(jsonDataTestRequest))
	if errPutSub != nil {
		t.Fatal(errPutSub)
	}
	rPutSub := httptest.NewRecorder()

	handlers.UpdateSubscription(rPutSub, reqPutSub)
	if status := rPutSub.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	///////// ListUserSubscriptions Test ///////////////////////////////////////////////////////////////////////////////
	reqLstSub, errLstSub := http.NewRequest("GET",
		fmt.Sprintf("/api/v1/subscriptions/user/%s", testRequest.UserID),
		nil)
	if errLstSub != nil {
		t.Fatal(errLstSub)
	}
	rLstSub := httptest.NewRecorder()

	handlers.ListUserSubscriptions(rLstSub, reqLstSub)
	if status := rLstSub.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	///////// TotalSubscriptions Test /////////////////////////////////////////////////////////////////////////////////
	reqCTSub, errCTSub := http.NewRequest("GET",
		fmt.Sprintf("/api/v1/subscriptions/total/?user_id=%s&service_name=%s&start_month=%s&end_month=%s",
			testRequest.UserID, testRequest.ServiceName, testRequest.StartDate, testRequest.EndDate),
		nil)
	if errCTSub != nil {
		t.Fatal(errCTSub)
	}
	rCTSub := httptest.NewRecorder()

	handlers.CalculateTotalHandler(rCTSub, reqCTSub)
	if status := rCTSub.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	///////// DELETESubscriptions Test /////////////////////////////////////////////////////////////////////////////////
	delReq := &models.DeleteRequest{
		ServiceName: testRequest.ServiceName,
		UserID:      testRequest.UserID,
	}
	jsonDelReq, _ := json.Marshal(delReq)
	reqDelSub, errDelSub := http.NewRequest("DELETE", "/api/v1/subscriptions/", bytes.NewBuffer(jsonDelReq))
	if errDelSub != nil {
		t.Fatal(errDelSub)
	}
	rDelSub := httptest.NewRecorder()

	handlers.DeleteSubscription(rDelSub, reqDelSub)
	if status := rDelSub.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
}
