package impl

import (
	"context"

	"github.com/go-openapi/runtime/middleware"

	"github.com/volmedo/pAPI/pkg/models"
	"github.com/volmedo/pAPI/pkg/restapi/operations/payments"
)

// PaymentsAPI implements the business logic needed to fulfill the API's requirements
type PaymentsAPI struct {
}

// CreatePayment Adds a new payment with the data included in params
func (papi *PaymentsAPI) CreatePayment(ctx context.Context, params payments.CreatePaymentParams) middleware.Responder {
	payment := *params.PaymentCreationRequest.Data
	resp := &models.PaymentCreationResponse{Data: &payment}
	return payments.NewCreatePaymentCreated().WithPayload(resp)
}
