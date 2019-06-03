package service

import (
	"context"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/volmedo/pAPI/pkg/models"
	"github.com/volmedo/pAPI/pkg/restapi/operations/payments"
)

// PaymentsService implements the business logic needed to fulfill the API's requirements
type PaymentsService struct {
	// Repo is a repository for payments
	Repo PaymentRepository
}

// CreatePayment Adds a new payment with the data included in params
func (papi *PaymentsService) CreatePayment(ctx context.Context, params payments.CreatePaymentParams) middleware.Responder {
	payment := params.PaymentCreationRequest.Data
	created, err := papi.Repo.Add(payment)
	if err != nil {
		apiError := newAPIError(err.Error())
		if _, ok := err.(ErrConflict); ok {
			return payments.NewCreatePaymentConflict().WithPayload(apiError)
		}

		return payments.NewCreatePaymentInternalServerError().WithPayload(apiError)
	}

	resp := &models.PaymentCreationResponse{Data: created}
	return payments.NewCreatePaymentCreated().WithPayload(resp)
}

// DeletePayment Deletes a payment identified by its ID
func (papi *PaymentsService) DeletePayment(ctx context.Context, params payments.DeletePaymentParams) middleware.Responder {
	paymentID := params.ID
	err := papi.Repo.Delete(paymentID)
	if err != nil {
		apiError := newAPIError(err.Error())
		if _, ok := err.(ErrNoResults); ok {
			return payments.NewDeletePaymentNotFound().WithPayload(apiError)
		}

		return payments.NewDeletePaymentInternalServerError().WithPayload(apiError)
	}

	return payments.NewDeletePaymentNoContent()
}

// GetPayment Returns details of a payment identified by its ID
func (papi *PaymentsService) GetPayment(ctx context.Context, params payments.GetPaymentParams) middleware.Responder {
	paymentID := params.ID
	got, err := papi.Repo.Get(paymentID)
	if err != nil {
		apiError := newAPIError(err.Error())
		if _, ok := err.(ErrNoResults); ok {
			return payments.NewGetPaymentNotFound().WithPayload(apiError)
		}

		return payments.NewGetPaymentInternalServerError().WithPayload(apiError)
	}

	resp := &models.PaymentDetailsResponse{Data: got}
	return payments.NewGetPaymentOK().WithPayload(resp)
}

// ListPayments Returns details of a collection of payments
func (papi *PaymentsService) ListPayments(ctx context.Context, params payments.ListPaymentsParams) middleware.Responder {
	// Request params have already been validated by go-swagger generated code
	pageNumber := *params.PageNumber
	pageSize := *params.PageSize

	offset := pageNumber * pageSize
	limit := pageSize
	list, err := papi.Repo.List(offset, limit)
	if err != nil {
		apiError := newAPIError(err.Error())
		if _, ok := err.(ErrBadOffsetLimit); ok {
			return payments.NewListPaymentsBadRequest().WithPayload(apiError)
		}

		return payments.NewListPaymentsInternalServerError().WithPayload(apiError)
	}

	resp := &models.PaymentDetailsListResponse{Data: list}
	return payments.NewListPaymentsOK().WithPayload(resp)
}

// UpdatePayment Adds a new payment with the data included in params
func (papi *PaymentsService) UpdatePayment(ctx context.Context, params payments.UpdatePaymentParams) middleware.Responder {
	paymentID := params.ID
	payment := params.PaymentUpdateRequest.Data
	updated, err := papi.Repo.Update(paymentID, payment)
	if err != nil {
		apiError := newAPIError(err.Error())
		if _, ok := err.(ErrNoResults); ok {
			return payments.NewUpdatePaymentNotFound().WithPayload(apiError)
		}

		return payments.NewUpdatePaymentInternalServerError().WithPayload(apiError)
	}

	resp := &models.PaymentUpdateResponse{Data: updated}
	return payments.NewUpdatePaymentOK().WithPayload(resp)
}

func newAPIError(msg string) *models.APIError {
	errorCode, _ := uuid.NewV4()
	return &models.APIError{
		ErrorCode:    strfmt.UUID(errorCode.String()),
		ErrorMessage: msg,
	}
}
