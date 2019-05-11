package impl

import (
	"context"
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/volmedo/pAPI/pkg/models"
	"github.com/volmedo/pAPI/pkg/restapi/operations/payments"
)

// PaymentsAPI implements the business logic needed to fulfill the API's requirements
type PaymentsAPI struct {
	store map[string]*models.Payment
}

func NewPaymentsAPI() *PaymentsAPI {
	newStore := map[string]*models.Payment{}
	return &PaymentsAPI{
		store: newStore,
	}
}

// CreatePayment Adds a new payment with the data included in params
func (papi *PaymentsAPI) CreatePayment(ctx context.Context, params payments.CreatePaymentParams) middleware.Responder {
	payment := *params.PaymentCreationRequest.Data
	paymentID := (*payment.ID).String()
	if _, ok := papi.store[paymentID]; ok {
		errorCode, _ := uuid.NewV4()
		apiError := models.APIError{
			ErrorCode:    strfmt.UUID(errorCode.String()),
			ErrorMessage: fmt.Sprintf("Payment ID %s already exists", paymentID),
		}
		return payments.NewCreatePaymentConflict().WithPayload(&apiError)
	}

	papi.store[paymentID] = &payment

	resp := &models.PaymentCreationResponse{Data: &payment}
	return payments.NewCreatePaymentCreated().WithPayload(resp)
}
