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

// PaymentRepository stores a collection of payment resources
type PaymentRepository map[strfmt.UUID]*models.Payment

// PaymentsAPI implements the business logic needed to fulfill the API's requirements
type PaymentsAPI struct {
	// Repo is a repository for payments
	Repo PaymentRepository
}

// CreatePayment Adds a new payment with the data included in params
func (papi *PaymentsAPI) CreatePayment(ctx context.Context, params payments.CreatePaymentParams) middleware.Responder {
	payment := *params.PaymentCreationRequest.Data
	paymentID := payment.ID.DeepCopy()
	if _, ok := papi.Repo[*paymentID]; ok {
		apiError := newAPIError(fmt.Sprintf("Payment ID %s already exists", paymentID))
		return payments.NewCreatePaymentConflict().WithPayload(apiError)
	}

	papi.Repo[*paymentID] = &payment

	respData := payment
	resp := &models.PaymentCreationResponse{Data: &respData}
	return payments.NewCreatePaymentCreated().WithPayload(resp)
}

// GetPayment Returns details of a payment identified by its ID
func (papi *PaymentsAPI) GetPayment(ctx context.Context, params payments.GetPaymentParams) middleware.Responder {
	paymentID := params.ID.DeepCopy()
	if payment, ok := papi.Repo[*paymentID]; ok {
		respData := *payment
		resp := &models.PaymentDetailsResponse{Data: &respData}
		return payments.NewGetPaymentOK().WithPayload(resp)
	}

	apiError := newAPIError(fmt.Sprintf("Payment with ID %s not found", paymentID))
	return payments.NewGetPaymentNotFound().WithPayload(apiError)
}

func newAPIError(msg string) *models.APIError {
	errorCode, _ := uuid.NewV4()
	return &models.APIError{
		ErrorCode:    strfmt.UUID(errorCode.String()),
		ErrorMessage: msg,
	}
}
