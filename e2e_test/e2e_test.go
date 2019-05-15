package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/colors"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/go-openapi/strfmt"

	"github.com/volmedo/pAPI/pkg/client"
	"github.com/volmedo/pAPI/pkg/client/payments"
	"github.com/volmedo/pAPI/pkg/models"
)

var (
	scheme   string
	host     string
	port     int
	basePath string

	opt = godog.Options{
		Output: colors.Colored(os.Stdout),
		Format: "progress",
	}
)

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

type Client struct {
	*payments.Client
	lastResponse interface{}
	lastError    error
}

func newClient(apiURL *url.URL) *Client {
	conf := client.Config{URL: apiURL}
	payments := client.New(conf)
	return &Client{payments.Payments, nil, nil}
}

func (c *Client) iCreateANewPaymentDescribedInJSONAs(jsonPayment *gherkin.DocString) error {
	var payment models.Payment
	decoder := json.NewDecoder(bytes.NewBuffer([]byte(jsonPayment.Content)))
	err := decoder.Decode(&payment)
	if err != nil {
		return fmt.Errorf("Invalid JSON string in test specification: %s", err)
	}

	ctx := context.Background()
	req := &models.PaymentCreationRequest{Data: &payment}
	params := payments.NewCreatePaymentParams().WithPaymentCreationRequest(req)

	c.lastResponse, c.lastError = c.CreatePayment(ctx, params)
	return nil
}

func (c *Client) iDeleteThePaymentWithID(paymentID string) error {
	ctx := context.Background()
	params := payments.NewDeletePaymentParams().WithID(strfmt.UUID(paymentID))

	c.lastResponse, c.lastError = c.DeletePayment(ctx, params)
	return nil
}

func (c *Client) iRequestThePaymentWithID(paymentID string) error {
	ctx := context.Background()
	params := payments.NewGetPaymentParams().WithID(strfmt.UUID(paymentID))

	c.lastResponse, c.lastError = c.GetPayment(ctx, params)
	return nil
}

func (c *Client) iGetAResponse(expectedStatus string) error {
	// An error will be raised in case of error but also if the StatusCode in the response
	// doesn't match the expected status, the generated code already takes care of this
	if c.lastError != nil {
		return fmt.Errorf("Error processing response or unexpected status: %s", c.lastError)
	}

	return nil
}

func (c *Client) theResponseContainsAPaymentDescribedInJSONAs(jsonPayment *gherkin.DocString) error {
	var gotPayment models.Payment
	switch resp := c.lastResponse.(type) {
	case *payments.CreatePaymentCreated:
		respData := resp.Payload.Data
		if respData == nil {
			return errors.New("Empty response")
		}
		gotPayment = *respData
	case *payments.GetPaymentOK:
		respData := resp.Payload.Data
		if respData == nil {
			return errors.New("Empty response")
		}
		gotPayment = *respData
	}

	var expectedPayment models.Payment
	decoder := json.NewDecoder(bytes.NewBuffer([]byte(jsonPayment.Content)))
	err := decoder.Decode(&expectedPayment)
	if err != nil {
		return fmt.Errorf("Invalid JSON string in test specification: %s", err)
	}

	if !reflect.DeepEqual(gotPayment, expectedPayment) {
		return errors.New("Payment data in the response don't match expected payment data")
	}

	return nil
}

func TestMain(m *testing.M) {
	flag.StringVar(&scheme, "scheme", "http", "Scheme to use to communicate with the server ('http' or 'https')")
	flag.StringVar(&host, "host", client.DefaultHost, "Address or URL of the server serving the Payments API (such as 'localhost' or 'api.example.com')")
	flag.IntVar(&port, "port", 8080, "Port where the server is listening for connections")
	flag.StringVar(&basePath, "base-path", client.DefaultBasePath, "Base path for API endpoints")

	flag.Parse()
	opt.Paths = flag.Args()

	status := godog.RunWithOptions("papi-e2e", func(s *godog.Suite) {
		FeatureContext(s)
	}, opt)

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func FeatureContext(s *godog.Suite) {
	apiURL := &url.URL{
		Scheme: scheme,
		Host:   fmt.Sprintf("%s:%d", host, port),
		Path:   basePath,
	}
	client := newClient(apiURL)

	s.Step(`^I create a new payment described in JSON as:$`, client.iCreateANewPaymentDescribedInJSONAs)
	s.Step(`^there is a payment described in JSON as:$`, client.iCreateANewPaymentDescribedInJSONAs)
	s.Step(`^I delete the payment with ID "([^"]*)"$`, client.iDeleteThePaymentWithID)
	s.Step(`^I request the payment with ID "([^"]*)"$`, client.iRequestThePaymentWithID)
	s.Step(`^I get a "([^"]*)" response$`, client.iGetAResponse)
	s.Step(`^I get an "([^"]*)" response$`, client.iGetAResponse)
	s.Step(`^the response contains a payment described in JSON as:$`, client.theResponseContainsAPaymentDescribedInJSONAs)
}
