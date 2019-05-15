// Code generated by go-swagger; DO NOT EDIT.

package payments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/volmedo/pAPI/pkg/models"
)

// UpdatePaymentReader is a Reader for the UpdatePayment structure.
type UpdatePaymentReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdatePaymentReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewUpdatePaymentOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 404:
		result := NewUpdatePaymentNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewUpdatePaymentOK creates a UpdatePaymentOK with default headers values
func NewUpdatePaymentOK() *UpdatePaymentOK {
	return &UpdatePaymentOK{}
}

/*UpdatePaymentOK handles this case with default header values.

Payment details
*/
type UpdatePaymentOK struct {
	Payload *models.PaymentUpdateResponse
}

func (o *UpdatePaymentOK) Error() string {
	return fmt.Sprintf("[PUT /payments/{id}][%d] updatePaymentOK  %+v", 200, o.Payload)
}

func (o *UpdatePaymentOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.PaymentUpdateResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdatePaymentNotFound creates a UpdatePaymentNotFound with default headers values
func NewUpdatePaymentNotFound() *UpdatePaymentNotFound {
	return &UpdatePaymentNotFound{}
}

/*UpdatePaymentNotFound handles this case with default header values.

Payment Not Found
*/
type UpdatePaymentNotFound struct {
	Payload *models.APIError
}

func (o *UpdatePaymentNotFound) Error() string {
	return fmt.Sprintf("[PUT /payments/{id}][%d] updatePaymentNotFound  %+v", 404, o.Payload)
}

func (o *UpdatePaymentNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}