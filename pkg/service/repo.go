package service

import (
	"github.com/go-openapi/strfmt"
	"github.com/volmedo/pAPI/pkg/models"
)

// TYPE_PAYMENT is a constant string that contains the value of the Type
// attribute of every payment resource
const TYPE_PAYMENT = "Payment"

// PaymentRepository stores a collection of payment resources that
// is safe for concurrent use
type PaymentRepository interface {
	// Add adds a new payment resource to the repository
	//
	// Add returns an error if a payment with the same ID as the one
	// to be added already exists
	Add(payment *models.Payment) (*models.Payment, error)

	// Delete deletes the payment resource associated to the given paymentID
	//
	// Delete returns an error if the paymentID is not present in the respository
	Delete(paymentID strfmt.UUID) error

	// Get returns the payment resource associated with the given paymentID
	//
	// Get returns an error if the paymentID does not exist in the collection
	Get(paymentID strfmt.UUID) (*models.Payment, error)

	// List returns a slice of payment resources. An empty slice will be returned
	// if no payment exists.
	//
	// List implements basic pagination by means of offset and limit parameters.
	// List will return an error if offset is beyond the number of elements available.
	// A limit of 0 will return all elements available. Both parameters default to 0.
	List(offset, limit int64) ([]*models.Payment, error)

	// Update updates the details associated with the given paymentID
	//
	// Update returns an error if the paymentID does not exist in the collection
	Update(paymentID strfmt.UUID, payment *models.Payment) (*models.Payment, error)
}

// ErrConflict signals an attempt to add a new payment with the same
// id as one already present
type ErrConflict string

func newErrConflict(msg string) ErrConflict {
	return ErrConflict(msg)
}

// Error satisfies stdlib's error interface
func (e ErrConflict) Error() string {
	return string(e)
}

// ErrNoResults is returned when a get, delete or update is attempted
// for a non-existent id
type ErrNoResults string

func newErrNoResults(msg string) ErrNoResults {
	return ErrNoResults(msg)
}

// Error satisfies stdlib's error interface
func (e ErrNoResults) Error() string {
	return string(e)
}

// ErrBadOffsetLimit is returned when invalid pagination parameters are
// passed when listing payments
type ErrBadOffsetLimit string

func newErrBadOffsetLimit(msg string) ErrBadOffsetLimit {
	return ErrBadOffsetLimit(msg)
}

// Error satisfies stdlib's error interface
func (e ErrBadOffsetLimit) Error() string {
	return string(e)
}
