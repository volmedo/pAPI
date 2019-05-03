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
