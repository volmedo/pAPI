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

// GetPaymentReader is a Reader for the GetPayment structure.
type GetPaymentReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetPaymentReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetPaymentOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 404:
		result := NewGetPaymentNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 500:
		result := NewGetPaymentInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetPaymentOK creates a GetPaymentOK with default headers values
func NewGetPaymentOK() *GetPaymentOK {
	return &GetPaymentOK{}
}

/*GetPaymentOK handles this case with default header values.

Payment details
*/
type GetPaymentOK struct {
	Payload *models.PaymentDetailsResponse
}

func (o *GetPaymentOK) Error() string {
	return fmt.Sprintf("[GET /payments/{id}][%d] getPaymentOK  %+v", 200, o.Payload)
}

func (o *GetPaymentOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.PaymentDetailsResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetPaymentNotFound creates a GetPaymentNotFound with default headers values
func NewGetPaymentNotFound() *GetPaymentNotFound {
	return &GetPaymentNotFound{}
}

/*GetPaymentNotFound handles this case with default header values.

Payment Not Found
*/
type GetPaymentNotFound struct {
	Payload *models.APIError
}

func (o *GetPaymentNotFound) Error() string {
	return fmt.Sprintf("[GET /payments/{id}][%d] getPaymentNotFound  %+v", 404, o.Payload)
}

func (o *GetPaymentNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetPaymentInternalServerError creates a GetPaymentInternalServerError with default headers values
func NewGetPaymentInternalServerError() *GetPaymentInternalServerError {
	return &GetPaymentInternalServerError{}
}

/*GetPaymentInternalServerError handles this case with default header values.

Internal Server Error
*/
type GetPaymentInternalServerError struct {
	Payload *models.APIError
}

func (o *GetPaymentInternalServerError) Error() string {
	return fmt.Sprintf("[GET /payments/{id}][%d] getPaymentInternalServerError  %+v", 500, o.Payload)
}

func (o *GetPaymentInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
