package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/volmedo/pAPI/pkg/models"
	"github.com/volmedo/pAPI/pkg/restapi/operations/payments"
)

// PaymentRepository stores a collection of payment resources that
// is safe for concurrent use
type PaymentRepository struct {
	sync.RWMutex
	m map[strfmt.UUID]*models.Payment
}

// NewPaymentRepository creates a freshly brewed PaymentRepository
func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{
		m: make(map[strfmt.UUID]*models.Payment),
	}
}

// Add adds a new payment resource to the repository
//
// Add returns an error if a payment with the same ID as the one
// to be added already exists
func (pr *PaymentRepository) Add(payment *models.Payment) error {
	paymentID := payment.ID.DeepCopy()
	pr.RLock()
	_, ok := pr.m[*paymentID]
	pr.RUnlock()
	if ok {
		return fmt.Errorf("Payment ID %s already exists", paymentID)
	}

	pr.Lock()
	pr.m[*paymentID] = payment
	pr.Unlock()
	return nil
}

// Delete deletes the payment resource associated to the given paymentID
//
// Delete returns an error if the paymentID is not present in the respository
func (pr *PaymentRepository) Delete(paymentID strfmt.UUID) error {
	pr.RLock()
	_, ok := pr.m[paymentID]
	pr.RUnlock()
	if !ok {
		return fmt.Errorf("Payment with ID %s not found", paymentID)
	}

	pr.Lock()
	delete(pr.m, paymentID)
	pr.Unlock()
	return nil
}

// Get returns the payment resource associated with the given paymentID
//
// Get returns an error if the paymentID does not exist in the collection
func (pr *PaymentRepository) Get(paymentID strfmt.UUID) (*models.Payment, error) {
	pr.RLock()
	payment, ok := pr.m[paymentID]
	pr.RUnlock()
	if !ok {
		return nil, fmt.Errorf("Payment with ID %s not found", paymentID)
	}

	return payment, nil
}

// Update updates the details associated with the given paymentID
//
// Update returns an error if the paymentID does not exist in the collection
func (pr *PaymentRepository) Update(paymentID strfmt.UUID, payment *models.Payment) error {
	pr.RLock()
	_, ok := pr.m[paymentID]
	pr.RUnlock()
	if !ok {
		return fmt.Errorf("Payment with ID %s not found", paymentID)
	}

	pr.Lock()
	pr.m[paymentID] = payment
	pr.Unlock()
	return nil
}

// PaymentsService implements the business logic needed to fulfill the API's requirements
type PaymentsService struct {
	// Repo is a repository for payments
	Repo *PaymentRepository
}

// CreatePayment Adds a new payment with the data included in params
func (papi *PaymentsService) CreatePayment(ctx context.Context, params payments.CreatePaymentParams) middleware.Responder {
	payment := params.PaymentCreationRequest.Data
	err := papi.Repo.Add(payment)
	if err != nil {
		apiError := newAPIError(err.Error())
		return payments.NewCreatePaymentConflict().WithPayload(apiError)
	}

	respData := payment
	resp := &models.PaymentCreationResponse{Data: respData}
	return payments.NewCreatePaymentCreated().WithPayload(resp)
}

// DeletePayment Deletes a payment identified by its ID
func (papi *PaymentsService) DeletePayment(ctx context.Context, params payments.DeletePaymentParams) middleware.Responder {
	paymentID := params.ID
	err := papi.Repo.Delete(paymentID)
	if err != nil {
		apiError := newAPIError(err.Error())
		return payments.NewDeletePaymentNotFound().WithPayload(apiError)
	}

	return payments.NewDeletePaymentNoContent()
}

// GetPayment Returns details of a payment identified by its ID
func (papi *PaymentsService) GetPayment(ctx context.Context, params payments.GetPaymentParams) middleware.Responder {
	paymentID := params.ID.DeepCopy()
	payment, err := papi.Repo.Get(*paymentID)
	if err != nil {
		apiError := newAPIError(err.Error())
		return payments.NewGetPaymentNotFound().WithPayload(apiError)
	}

	resp := &models.PaymentDetailsResponse{Data: payment}
	return payments.NewGetPaymentOK().WithPayload(resp)
}

// UpdatePayment Adds a new payment with the data included in params
func (papi *PaymentsService) UpdatePayment(ctx context.Context, params payments.UpdatePaymentParams) middleware.Responder {
	paymentID := params.ID.DeepCopy()
	payment := params.PaymentUpdateRequest.Data
	err := papi.Repo.Update(*paymentID, payment)
	if err != nil {
		apiError := newAPIError(err.Error())
		return payments.NewUpdatePaymentNotFound().WithPayload(apiError)
	}

	resp := &models.PaymentUpdateResponse{Data: payment}
	return payments.NewUpdatePaymentOK().WithPayload(resp)
}

func newAPIError(msg string) *models.APIError {
	errorCode, _ := uuid.NewV4()
	return &models.APIError{
		ErrorCode:    strfmt.UUID(errorCode.String()),
		ErrorMessage: msg,
	}
}
