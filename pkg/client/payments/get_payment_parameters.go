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
)

// NewGetPaymentParams creates a new GetPaymentParams object
// with the default values initialized.
func NewGetPaymentParams() *GetPaymentParams {
	var ()
	return &GetPaymentParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetPaymentParamsWithTimeout creates a new GetPaymentParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetPaymentParamsWithTimeout(timeout time.Duration) *GetPaymentParams {
	var ()
	return &GetPaymentParams{

		timeout: timeout,
	}
}

// NewGetPaymentParamsWithContext creates a new GetPaymentParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetPaymentParamsWithContext(ctx context.Context) *GetPaymentParams {
	var ()
	return &GetPaymentParams{

		Context: ctx,
	}
}

// NewGetPaymentParamsWithHTTPClient creates a new GetPaymentParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetPaymentParamsWithHTTPClient(client *http.Client) *GetPaymentParams {
	var ()
	return &GetPaymentParams{
		HTTPClient: client,
	}
}

/*GetPaymentParams contains all the parameters to send to the API endpoint
for the get payment operation typically these are written to a http.Request
*/
type GetPaymentParams struct {

	/*ID
	  ID of payment to fetch

	*/
	ID strfmt.UUID

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get payment params
func (o *GetPaymentParams) WithTimeout(timeout time.Duration) *GetPaymentParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get payment params
func (o *GetPaymentParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get payment params
func (o *GetPaymentParams) WithContext(ctx context.Context) *GetPaymentParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get payment params
func (o *GetPaymentParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get payment params
func (o *GetPaymentParams) WithHTTPClient(client *http.Client) *GetPaymentParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get payment params
func (o *GetPaymentParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the get payment params
func (o *GetPaymentParams) WithID(id strfmt.UUID) *GetPaymentParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get payment params
func (o *GetPaymentParams) SetID(id strfmt.UUID) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *GetPaymentParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", o.ID.String()); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
