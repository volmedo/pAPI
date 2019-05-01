// Code generated by go-swagger; DO NOT EDIT.

package payments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/volmedo/pAPI/pkg/models"
)

// CreatePaymentCreatedCode is the HTTP code returned for type CreatePaymentCreated
const CreatePaymentCreatedCode int = 201

/*CreatePaymentCreated Payment creation response

swagger:response createPaymentCreated
*/
type CreatePaymentCreated struct {

	/*
	  In: Body
	*/
	Payload *models.PaymentCreationResponse `json:"body,omitempty"`
}

// NewCreatePaymentCreated creates CreatePaymentCreated with default headers values
func NewCreatePaymentCreated() *CreatePaymentCreated {

	return &CreatePaymentCreated{}
}

// WithPayload adds the payload to the create payment created response
func (o *CreatePaymentCreated) WithPayload(payload *models.PaymentCreationResponse) *CreatePaymentCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create payment created response
func (o *CreatePaymentCreated) SetPayload(payload *models.PaymentCreationResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreatePaymentCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
