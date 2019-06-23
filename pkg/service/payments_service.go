package service

import (
	"context"
	"fmt"
	"log"

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

	// Logger will be use to write logs. Only unexpected errors will be logged
	Logger *log.Logger
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

		papi.Logger.Printf("Error on CreatePayment: %v", err)
		return payments.NewCreatePaymentInternalServerError().WithPayload(apiError)
	}

	links := &models.Links{
		Self: fmt.Sprintf("%s/%s", params.HTTPRequest.URL.Path, created.ID),
	}
	resp := &models.PaymentCreationResponse{Data: created, Links: links}
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

		papi.Logger.Printf("Error on DeletePayment: %v", err)
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

		papi.Logger.Printf("Error on GetPayment: %v", err)
		return payments.NewGetPaymentInternalServerError().WithPayload(apiError)
	}

	links := &models.Links{
		Self: params.HTTPRequest.URL.Path,
	}
	resp := &models.PaymentDetailsResponse{Data: got, Links: links}
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
		if _, ok := err.(ErrNoResults); ok {
			return payments.NewListPaymentsNotFound().WithPayload(apiError)
		}

		papi.Logger.Printf("Error on ListPayments: %v", err)
		return payments.NewListPaymentsInternalServerError().WithPayload(apiError)
	}

	links := &models.Links{
		Self: "/payments",
	}
	resp := &models.PaymentDetailsListResponse{
		Data:  list,
		Links: links,
	}
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

		papi.Logger.Printf("Error on UpdatePayment: %v", err)
		return payments.NewUpdatePaymentInternalServerError().WithPayload(apiError)
	}

	links := &models.Links{
		Self: params.HTTPRequest.URL.Path,
	}
	resp := &models.PaymentUpdateResponse{Data: updated, Links: links}
	return payments.NewUpdatePaymentOK().WithPayload(resp)
}

func newAPIError(msg string) *models.APIError {
	errorCode, _ := uuid.NewV4()
	return &models.APIError{
		ErrorCode:    strfmt.UUID(errorCode.String()),
		ErrorMessage: msg,
	}
}
