// Code generated by go-swagger; DO NOT EDIT.

package payments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/volmedo/pAPI/pkg/models"
)

// NewUpdatePaymentParams creates a new UpdatePaymentParams object
// with the default values initialized.
func NewUpdatePaymentParams() *UpdatePaymentParams {
	var ()
	return &UpdatePaymentParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewUpdatePaymentParamsWithTimeout creates a new UpdatePaymentParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewUpdatePaymentParamsWithTimeout(timeout time.Duration) *UpdatePaymentParams {
	var ()
	return &UpdatePaymentParams{

		timeout: timeout,
	}
}

// NewUpdatePaymentParamsWithContext creates a new UpdatePaymentParams object
// with the default values initialized, and the ability to set a context for a request
func NewUpdatePaymentParamsWithContext(ctx context.Context) *UpdatePaymentParams {
	var ()
	return &UpdatePaymentParams{

		Context: ctx,
	}
}

// NewUpdatePaymentParamsWithHTTPClient creates a new UpdatePaymentParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewUpdatePaymentParamsWithHTTPClient(client *http.Client) *UpdatePaymentParams {
	var ()
	return &UpdatePaymentParams{
		HTTPClient: client,
	}
}

/*UpdatePaymentParams contains all the parameters to send to the API endpoint
for the update payment operation typically these are written to a http.Request
*/
type UpdatePaymentParams struct {

	/*PaymentUpdateRequest
	  New payment details

	*/
	PaymentUpdateRequest *models.PaymentUpdateRequest
	/*ID
	  ID of payment to update

	*/
	ID strfmt.UUID

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the update payment params
func (o *UpdatePaymentParams) WithTimeout(timeout time.Duration) *UpdatePaymentParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the update payment params
func (o *UpdatePaymentParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the update payment params
func (o *UpdatePaymentParams) WithContext(ctx context.Context) *UpdatePaymentParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the update payment params
func (o *UpdatePaymentParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the update payment params
func (o *UpdatePaymentParams) WithHTTPClient(client *http.Client) *UpdatePaymentParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the update payment params
func (o *UpdatePaymentParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithPaymentUpdateRequest adds the paymentUpdateRequest to the update payment params
func (o *UpdatePaymentParams) WithPaymentUpdateRequest(paymentUpdateRequest *models.PaymentUpdateRequest) *UpdatePaymentParams {
	o.SetPaymentUpdateRequest(paymentUpdateRequest)
	return o
}

// SetPaymentUpdateRequest adds the paymentUpdateRequest to the update payment params
func (o *UpdatePaymentParams) SetPaymentUpdateRequest(paymentUpdateRequest *models.PaymentUpdateRequest) {
	o.PaymentUpdateRequest = paymentUpdateRequest
}

// WithID adds the id to the update payment params
func (o *UpdatePaymentParams) WithID(id strfmt.UUID) *UpdatePaymentParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the update payment params
func (o *UpdatePaymentParams) SetID(id strfmt.UUID) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *UpdatePaymentParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.PaymentUpdateRequest != nil {
		if err := r.SetBodyParam(o.PaymentUpdateRequest); err != nil {
			return err
		}
	}

	// path param id
	if err := r.SetPathParam("id", o.ID.String()); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
