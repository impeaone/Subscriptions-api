package tests

import (
	"agrigation_api/internal/service"
	"agrigation_api/pkg/models"
	"agrigation_api/pkg/tools"
	"context"
	"github.com/google/uuid"
	"testing"
)

func TestService(t *testing.T) {
	testRepo, _ := NewTestRepository()
	serv := service.NewSubscriptionService(testRepo)

	testRequest := models.CreateOrUpdateRequest{
		ServiceName: "test_service",
		Price:       300,
		UserID:      uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba"),
		StartDate:   "07-2025",
		EndDate:     "08-2025",
	}
	startDate, _ := tools.ParseMonthYear(testRequest.StartDate)
	testResponse := models.Subscription{
		ServiceName: testRequest.ServiceName,
		Price:       testRequest.Price,
		UserID:      testRequest.UserID,
		StartDate:   startDate,
	}
	ctx := context.Background()
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	sub, err := serv.CreateSubscription(ctx, testRequest)
	if err != nil {
		t.Fatal(err)
	}
	if sub == nil {
		t.Fatal("sub is nil")
	}
	if sub.ServiceName != testResponse.ServiceName || sub.Price != testResponse.Price || sub.UserID !=
		testResponse.UserID || sub.StartDate != testResponse.StartDate {
		t.Fatal("sub is wrong")
	}
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	testRequestUpdate := models.CreateOrUpdateRequest{
		ServiceName: testRequest.ServiceName,
		Price:       testRequest.Price + 100,
		UserID:      testRequest.UserID,
		StartDate:   testRequest.StartDate,
		EndDate:     testRequest.EndDate,
	}
	sub, errUp := serv.UpdateSubscription(context.Background(), testRequestUpdate)
	if errUp != nil {
		t.Fatal(errUp)
	}
	if sub == nil {
		t.Fatal("sub is nil")
	}
	if sub.ServiceName != testResponse.ServiceName || sub.Price != testResponse.Price+100 || sub.UserID !=
		testResponse.UserID || sub.StartDate != testResponse.StartDate {
		t.Fatal("sub is wrong")
	}
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	testRequestTotal := models.CalculateTotalRequest{
		ServiceName: testRequest.ServiceName,
		UserID:      testRequest.UserID,
		StartMonth:  startDate,
	}

	price, err := serv.CalculateTotal(context.Background(), testRequestTotal)
	if err != nil {
		t.Fatal(err)
	}
	if price != 0 {
		t.Fatal("price is wrong")
	}
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
}
