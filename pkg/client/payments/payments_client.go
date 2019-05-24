// Code generated by go-swagger; DO NOT EDIT.

package payments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

//go:generate mockery -name API -inpkg

// API is the interface of the payments client
type API interface {
	// CreatePayment creates payment
	CreatePayment(ctx context.Context, params *CreatePaymentParams) (*CreatePaymentCreated, error)
	// DeletePayment deletes a payment resource
	DeletePayment(ctx context.Context, params *DeletePaymentParams) (*DeletePaymentNoContent, error)
	// GetPayment fetches payment
	GetPayment(ctx context.Context, params *GetPaymentParams) (*GetPaymentOK, error)
	// ListPayments lists payments
	ListPayments(ctx context.Context, params *ListPaymentsParams) (*ListPaymentsOK, error)
	// UpdatePayment updates payment details
	UpdatePayment(ctx context.Context, params *UpdatePaymentParams) (*UpdatePaymentOK, error)
}

// New creates a new payments API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry, authInfo runtime.ClientAuthInfoWriter) *Client {
	return &Client{
		transport: transport,
		formats:   formats,
		authInfo:  authInfo,
	}
}

/*
Client for payments API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
	authInfo  runtime.ClientAuthInfoWriter
}

/*
CreatePayment creates payment
*/
func (a *Client) CreatePayment(ctx context.Context, params *CreatePaymentParams) (*CreatePaymentCreated, error) {

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "createPayment",
		Method:             "POST",
		PathPattern:        "/payments",
		ProducesMediaTypes: []string{"application/vnd.api+json"},
		ConsumesMediaTypes: []string{"application/vnd.api+json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &CreatePaymentReader{formats: a.formats},
		Context:            ctx,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*CreatePaymentCreated), nil

}

/*
DeletePayment deletes a payment resource
*/
func (a *Client) DeletePayment(ctx context.Context, params *DeletePaymentParams) (*DeletePaymentNoContent, error) {

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "deletePayment",
		Method:             "DELETE",
		PathPattern:        "/payments/{id}",
		ProducesMediaTypes: []string{"application/vnd.api+json"},
		ConsumesMediaTypes: []string{"application/vnd.api+json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &DeletePaymentReader{formats: a.formats},
		Context:            ctx,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*DeletePaymentNoContent), nil

}

/*
GetPayment fetches payment
*/
func (a *Client) GetPayment(ctx context.Context, params *GetPaymentParams) (*GetPaymentOK, error) {

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getPayment",
		Method:             "GET",
		PathPattern:        "/payments/{id}",
		ProducesMediaTypes: []string{"application/vnd.api+json"},
		ConsumesMediaTypes: []string{"application/vnd.api+json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetPaymentReader{formats: a.formats},
		Context:            ctx,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetPaymentOK), nil

}

/*
ListPayments lists payments
*/
func (a *Client) ListPayments(ctx context.Context, params *ListPaymentsParams) (*ListPaymentsOK, error) {

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "listPayments",
		Method:             "GET",
		PathPattern:        "/payments",
		ProducesMediaTypes: []string{"application/vnd.api+json"},
		ConsumesMediaTypes: []string{"application/vnd.api+json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &ListPaymentsReader{formats: a.formats},
		Context:            ctx,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*ListPaymentsOK), nil

}

/*
UpdatePayment updates payment details
*/
func (a *Client) UpdatePayment(ctx context.Context, params *UpdatePaymentParams) (*UpdatePaymentOK, error) {

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "updatePayment",
		Method:             "PUT",
		PathPattern:        "/payments/{id}",
		ProducesMediaTypes: []string{"application/vnd.api+json"},
		ConsumesMediaTypes: []string{"application/vnd.api+json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &UpdatePaymentReader{formats: a.formats},
		Context:            ctx,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*UpdatePaymentOK), nil

}
