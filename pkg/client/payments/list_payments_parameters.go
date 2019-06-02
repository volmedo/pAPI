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
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewListPaymentsParams creates a new ListPaymentsParams object
// with the default values initialized.
func NewListPaymentsParams() *ListPaymentsParams {
	var (
		pageNumberDefault = int64(0)
		pageSizeDefault   = int64(10)
	)
	return &ListPaymentsParams{
		PageNumber: &pageNumberDefault,
		PageSize:   &pageSizeDefault,

		timeout: cr.DefaultTimeout,
	}
}

// NewListPaymentsParamsWithTimeout creates a new ListPaymentsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewListPaymentsParamsWithTimeout(timeout time.Duration) *ListPaymentsParams {
	var (
		pageNumberDefault = int64(0)
		pageSizeDefault   = int64(10)
	)
	return &ListPaymentsParams{
		PageNumber: &pageNumberDefault,
		PageSize:   &pageSizeDefault,

		timeout: timeout,
	}
}

// NewListPaymentsParamsWithContext creates a new ListPaymentsParams object
// with the default values initialized, and the ability to set a context for a request
func NewListPaymentsParamsWithContext(ctx context.Context) *ListPaymentsParams {
	var (
		pageNumberDefault = int64(0)
		pageSizeDefault   = int64(10)
	)
	return &ListPaymentsParams{
		PageNumber: &pageNumberDefault,
		PageSize:   &pageSizeDefault,

		Context: ctx,
	}
}

// NewListPaymentsParamsWithHTTPClient creates a new ListPaymentsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewListPaymentsParamsWithHTTPClient(client *http.Client) *ListPaymentsParams {
	var (
		pageNumberDefault = int64(0)
		pageSizeDefault   = int64(10)
	)
	return &ListPaymentsParams{
		PageNumber: &pageNumberDefault,
		PageSize:   &pageSizeDefault,
		HTTPClient: client,
	}
}

/*ListPaymentsParams contains all the parameters to send to the API endpoint
for the list payments operation typically these are written to a http.Request
*/
type ListPaymentsParams struct {

	/*PageNumber
	  Which page to select

	*/
	PageNumber *int64
	/*PageSize
	  Number of items per page

	*/
	PageSize *int64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the list payments params
func (o *ListPaymentsParams) WithTimeout(timeout time.Duration) *ListPaymentsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the list payments params
func (o *ListPaymentsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the list payments params
func (o *ListPaymentsParams) WithContext(ctx context.Context) *ListPaymentsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the list payments params
func (o *ListPaymentsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the list payments params
func (o *ListPaymentsParams) WithHTTPClient(client *http.Client) *ListPaymentsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the list payments params
func (o *ListPaymentsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithPageNumber adds the pageNumber to the list payments params
func (o *ListPaymentsParams) WithPageNumber(pageNumber *int64) *ListPaymentsParams {
	o.SetPageNumber(pageNumber)
	return o
}

// SetPageNumber adds the pageNumber to the list payments params
func (o *ListPaymentsParams) SetPageNumber(pageNumber *int64) {
	o.PageNumber = pageNumber
}

// WithPageSize adds the pageSize to the list payments params
func (o *ListPaymentsParams) WithPageSize(pageSize *int64) *ListPaymentsParams {
	o.SetPageSize(pageSize)
	return o
}

// SetPageSize adds the pageSize to the list payments params
func (o *ListPaymentsParams) SetPageSize(pageSize *int64) {
	o.PageSize = pageSize
}

// WriteToRequest writes these params to a swagger request
func (o *ListPaymentsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.PageNumber != nil {

		// query param page[number]
		var qrPageNumber int64
		if o.PageNumber != nil {
			qrPageNumber = *o.PageNumber
		}
		qPageNumber := swag.FormatInt64(qrPageNumber)
		if qPageNumber != "" {
			if err := r.SetQueryParam("page[number]", qPageNumber); err != nil {
				return err
			}
		}

	}

	if o.PageSize != nil {

		// query param page[size]
		var qrPageSize int64
		if o.PageSize != nil {
			qrPageSize = *o.PageSize
		}
		qPageSize := swag.FormatInt64(qrPageSize)
		if qPageSize != "" {
			if err := r.SetQueryParam("page[size]", qPageSize); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}