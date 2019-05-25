package service

import (
	"fmt"
	"sort"
	"sync"

	"github.com/go-openapi/strfmt"
	"github.com/volmedo/pAPI/pkg/models"
)

// MapPaymentRepository stores a collection of payment resources using
// a standard map in memory as data backend
type MapPaymentRepository struct {
	sync.RWMutex
	m map[strfmt.UUID]*models.Payment
}

// NewMapPaymentRepository creates a freshly brewed MapPaymentRepository
func NewMapPaymentRepository() *MapPaymentRepository {
	return &MapPaymentRepository{
		m: make(map[strfmt.UUID]*models.Payment),
	}
}

// Add adds a new payment resource to the repository
//
// Add returns an error if a payment with the same ID as the one
// to be added already exists
func (mpr *MapPaymentRepository) Add(payment *models.Payment) error {
	paymentID := payment.ID.DeepCopy()
	mpr.RLock()
	_, ok := mpr.m[*paymentID]
	mpr.RUnlock()
	if ok {
		return fmt.Errorf("Payment ID %s already exists", paymentID)
	}

	mpr.Lock()
	mpr.m[*paymentID] = payment
	mpr.Unlock()
	return nil
}

// Delete deletes the payment resource associated to the given paymentID
//
// Delete returns an error if the paymentID is not present in the respository
func (mpr *MapPaymentRepository) Delete(paymentID strfmt.UUID) error {
	mpr.RLock()
	_, ok := mpr.m[paymentID]
	mpr.RUnlock()
	if !ok {
		return fmt.Errorf("Payment with ID %s not found", paymentID)
	}

	mpr.Lock()
	delete(mpr.m, paymentID)
	mpr.Unlock()
	return nil
}

// Get returns the payment resource associated with the given paymentID
//
// Get returns an error if the paymentID does not exist in the collection
func (mpr *MapPaymentRepository) Get(paymentID strfmt.UUID) (*models.Payment, error) {
	mpr.RLock()
	payment, ok := mpr.m[paymentID]
	mpr.RUnlock()
	if !ok {
		return nil, fmt.Errorf("Payment with ID %s not found", paymentID)
	}

	return payment, nil
}

// List returns a slice of payment resources. An empty slice will be returned
// if no payment exists.
//
// List implements basic pagination by means of offset and limit parameters.
// List will return an error if offset is beyond the number of elements available.
// A limit of 0 will return all elements available. Both parameters default to 0.
func (mpr *MapPaymentRepository) List(offset, limit int64) ([]*models.Payment, error) {
	// Check params before anything else
	from := offset
	to := offset + limit
	if from >= int64(len(mpr.m)) {
		return nil, fmt.Errorf("Requested item at %d but only %d items exist", from, len(mpr.m))
	}
	if limit == 0 || to > int64(len(mpr.m)) {
		to = int64(len(mpr.m))
	}

	mpr.RLock()
	var ids []string
	for id := range mpr.m {
		ids = append(ids, id.String())
	}
	mpr.RUnlock()

	sort.Strings(ids)

	ids = ids[from:to]
	payments := make([]*models.Payment, 0, len(ids))
	mpr.RLock()
	for _, id := range ids {
		payments = append(payments, mpr.m[strfmt.UUID(id)])
	}
	mpr.RUnlock()

	return payments, nil
}

// Update updates the details associated with the given paymentID
//
// Update returns an error if the paymentID does not exist in the collection
func (mpr *MapPaymentRepository) Update(paymentID strfmt.UUID, payment *models.Payment) error {
	mpr.RLock()
	_, ok := mpr.m[paymentID]
	mpr.RUnlock()
	if !ok {
		return fmt.Errorf("Payment with ID %s not found", paymentID)
	}

	mpr.Lock()
	mpr.m[paymentID] = payment
	mpr.Unlock()
	return nil
}
