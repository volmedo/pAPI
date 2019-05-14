// Code generated by go-swagger; DO NOT EDIT.

package payments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetPaymentHandlerFunc turns a function with the right signature into a get payment handler
type GetPaymentHandlerFunc func(GetPaymentParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetPaymentHandlerFunc) Handle(params GetPaymentParams) middleware.Responder {
	return fn(params)
}

// GetPaymentHandler interface for that can handle valid get payment params
type GetPaymentHandler interface {
	Handle(GetPaymentParams) middleware.Responder
}

// NewGetPayment creates a new http.Handler for the get payment operation
func NewGetPayment(ctx *middleware.Context, handler GetPaymentHandler) *GetPayment {
	return &GetPayment{Context: ctx, Handler: handler}
}

/*GetPayment swagger:route GET /payments/{id} Payments getPayment

Fetch payment

*/
type GetPayment struct {
	Context *middleware.Context
	Handler GetPaymentHandler
}

func (o *GetPayment) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetPaymentParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}