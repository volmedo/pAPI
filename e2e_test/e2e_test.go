package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/colors"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/go-openapi/strfmt"

	"github.com/volmedo/pAPI/pkg/client"
	"github.com/volmedo/pAPI/pkg/client/payments"
	"github.com/volmedo/pAPI/pkg/models"
)

var (
	scheme     string
	host       string
	port       int
	apiPath    string
	healthPath string

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
	lastResponse  interface{}
	lastError     error
	registeredIDs map[strfmt.UUID]struct{} // This property allows for cleaning after each scenario
	healthURL     *url.URL
}

func newClient(apiURL, healthURL *url.URL) *Client {
	conf := client.Config{URL: apiURL}
	payments := client.New(conf)
	registeredIDs := make(map[strfmt.UUID]struct{})
	return &Client{
		Client:        payments.Payments,
		lastResponse:  nil,
		lastError:     nil,
		registeredIDs: registeredIDs,
		healthURL:     healthURL}
}

func (c *Client) thereArePaymentsWithIDs(ids *gherkin.DataTable) error {
	// For simplicity, the same test payment will be created with different IDs
	version := int64(0)
	orgID := strfmt.UUID("743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb")
	procDate, _ := time.Parse(strfmt.RFC3339FullDate, "2017-01-18")

	testPayment := models.Payment{
		Type:           "Payment",
		ID:             nil,
		Version:        &version,
		OrganisationID: &orgID,
		Attributes: &models.PaymentAttributes{
			Amount: "100.21",
			BeneficiaryParty: &models.PaymentParty{
				AccountName:       "W Owens",
				AccountNumber:     models.AccountNumber("31926819"),
				AccountNumberCode: "BBAN",
				AccountType:       0,
				Address:           "1 The Beneficiary Localtown SE2",
				BankID:            "403000",
				BankIDCode:        "GBDSC",
				Name:              "Wilfred Jeremiah Owens",
			},
			ChargesInformation: &models.ChargesInformation{
				BearerCode:              "SHAR",
				ReceiverChargesAmount:   "1.00",
				ReceiverChargesCurrency: "USD",
				SenderCharges: []*models.ChargesInformationSenderChargesItems0{
					{Amount: "5.00", Currency: "GBP"},
					{Amount: "10.00", Currency: "USD"},
				},
			},
			Currency: "GBP",
			DebtorParty: &models.PaymentParty{
				AccountName:       "EJ Brown Black",
				AccountNumber:     "GB29XABC10161234567801",
				AccountNumberCode: "IBAN",
				Address:           "10 Debtor Crescent Sourcetown NE1",
				BankID:            "203301",
				BankIDCode:        "GBDSC",
				Name:              "Emelia Jane Brown",
			},
			EndToEndReference: "Wil piano Jan",
			Fx: &models.PaymentAttributesFx{
				ContractReference: "FX123",
				ExchangeRate:      "2.00000",
				OriginalAmount:    "200.42",
				OriginalCurrency:  "USD",
			},
			NumericReference:     "1002001",
			PaymentID:            "123456789012345678",
			PaymentPurpose:       "Paying for goods/services",
			PaymentScheme:        "FPS",
			PaymentType:          "Credit",
			ProcessingDate:       strfmt.Date(procDate),
			Reference:            "Payment for Em's piano lessons",
			SchemePaymentSubType: "InternetBanking",
			SchemePaymentType:    "ImmediatePayment",
			SponsorParty: &models.PaymentAttributesSponsorParty{
				AccountNumber: "56781234",
				BankID:        "123123",
				BankIDCode:    "GBDSC",
			},
		},
	}

	for _, row := range ids.Rows {
		idString := row.Cells[0].Value
		id := strfmt.UUID(idString)
		testPayment.ID = &id

		ctx := context.Background()
		req := &models.PaymentCreationRequest{Data: &testPayment}
		params := payments.NewCreatePaymentParams().WithPaymentCreationRequest(req)

		_, err := c.CreatePayment(ctx, params)
		if err != nil {
			return fmt.Errorf("Error creating payment: %s", err)
		}

		c.registeredIDs[id] = struct{}{}
	}

	return nil
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
	if c.lastError == nil {
		c.registeredIDs[*payment.ID] = struct{}{}
	}
	return nil
}

func (c *Client) iDeleteThePaymentWithID(paymentID string) error {
	pID := strfmt.UUID(paymentID)
	ctx := context.Background()
	params := payments.NewDeletePaymentParams().WithID(pID)

	c.lastResponse, c.lastError = c.DeletePayment(ctx, params)
	if c.lastError == nil {
		delete(c.registeredIDs, pID)
	}
	return nil
}

func (c *Client) iRequestThePaymentWithID(paymentID string) error {
	ctx := context.Background()
	params := payments.NewGetPaymentParams().WithID(strfmt.UUID(paymentID))

	c.lastResponse, c.lastError = c.GetPayment(ctx, params)
	return nil
}

func (c *Client) iRequestAListOfPayments() error {
	ctx := context.Background()
	params := payments.NewListPaymentsParams()

	c.lastResponse, c.lastError = c.ListPayments(ctx, params)
	return nil
}

func (c *Client) iRequestAListOfPaymentsPageWithPaymentsPerPage(pageNumber, pageSize int) error {
	ctx := context.Background()
	pNumber := int64(pageNumber)
	pSize := int64(pageSize)
	params := payments.NewListPaymentsParams().WithPageNumber(&pNumber).WithPageSize(&pSize)

	c.lastResponse, c.lastError = c.ListPayments(ctx, params)
	return nil
}

func (c *Client) iUpdateThePaymentWithIDWithNewDetailsInJSON(paymentID string, jsonPayment *gherkin.DocString) error {
	var payment models.Payment
	decoder := json.NewDecoder(bytes.NewBuffer([]byte(jsonPayment.Content)))
	err := decoder.Decode(&payment)
	if err != nil {
		return fmt.Errorf("Invalid JSON string in test specification: %s", err)
	}

	ctx := context.Background()
	pID := strfmt.UUID(paymentID)
	req := &models.PaymentUpdateRequest{Data: &payment}
	params := payments.NewUpdatePaymentParams().WithID(pID).WithPaymentUpdateRequest(req)

	c.lastResponse, c.lastError = c.UpdatePayment(ctx, params)
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

	case *payments.UpdatePaymentOK:
		respData := resp.Payload.Data
		if respData == nil {
			return errors.New("Empty response")
		}
		gotPayment = *respData

	default:
		return fmt.Errorf("Unknown response type %T", resp)
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

func (c *Client) theResponseContainsAListOfPaymentsWithIDs(ids *gherkin.DataTable) error {
	resp, ok := c.lastResponse.(*payments.ListPaymentsOK)
	if !ok {
		return errors.New("Wrong response type")
	}

	gotPayments := resp.Payload.Data
	if len(gotPayments) == 0 {
		return errors.New("Empty data in response")
	}

	// Use maps for IDs so that they can be compared directly using reflect.DeepEqual
	gotIDs := make(map[strfmt.UUID]struct{})
	for _, payment := range gotPayments {
		if payment == nil {
			return errors.New("Empty payment in data")
		}
		gotIDs[*payment.ID] = struct{}{}
	}

	wantIDs := make(map[strfmt.UUID]struct{})
	for _, row := range ids.Rows {
		idString := row.Cells[0].Value
		id := strfmt.UUID(idString)
		wantIDs[id] = struct{}{}
	}

	if !reflect.DeepEqual(gotIDs, wantIDs) {
		return errors.New("Payment data in the response don't match expected payment data")
	}

	return nil
}

// ping checks if the API is ready by sending a request to its health endpoint
func (c *Client) ping() error {
	resp, err := http.Get(c.healthURL.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("API not available")
	}

	return nil
}

func TestMain(m *testing.M) {
	flag.StringVar(&scheme, "scheme", "http", "Scheme to use to communicate with the server ('http' or 'https')")
	flag.StringVar(&host, "host", client.DefaultHost, "Address or URL of the server serving the Payments API (such as 'localhost' or 'api.example.com')")
	flag.IntVar(&port, "port", 8080, "Port where the server is listening for connections")
	flag.StringVar(&apiPath, "api-path", client.DefaultBasePath, "Base path for API endpoints")
	flag.StringVar(&healthPath, "health-path", "/health", "Path to the API's health endpoint")

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
		Path:   apiPath,
	}
	healthURL := &url.URL{
		Scheme: scheme,
		Host:   fmt.Sprintf("%s:%d", host, port),
		Path:   healthPath,
	}
	client := newClient(apiURL, healthURL)

	s.Step(`^there are payments with IDs:$`, client.thereArePaymentsWithIDs)
	s.Step(`^the response contains a list of payments with the following IDs:$`, client.theResponseContainsAListOfPaymentsWithIDs)
	s.Step(`^I create a new payment described in JSON as:$`, client.iCreateANewPaymentDescribedInJSONAs)
	s.Step(`^there is a payment described in JSON as:$`, client.iCreateANewPaymentDescribedInJSONAs)
	s.Step(`^I delete the payment with ID "([^"]*)"$`, client.iDeleteThePaymentWithID)
	s.Step(`^I request the payment with ID "([^"]*)"$`, client.iRequestThePaymentWithID)
	s.Step(`^I update the payment with ID "([^"]*)" with new details in JSON:$`, client.iUpdateThePaymentWithIDWithNewDetailsInJSON)
	s.Step(`^I request a list of payments$`, client.iRequestAListOfPayments)
	s.Step(`^I request a list of payments, page (\d+) with (\d+) payments per page$`, client.iRequestAListOfPaymentsPageWithPaymentsPerPage)
	s.Step(`^I get a[n]? "([^"]*)" response$`, client.iGetAResponse)
	s.Step(`^the response contains a payment described in JSON as:$`, client.theResponseContainsAPaymentDescribedInJSONAs)

	// Wait for the application to settle before running the test suite
	s.BeforeSuite(func() {
		max_retries := 20
		err := client.ping()
		for i := 0; i < max_retries && err != nil; i++ {
			time.Sleep(1 * time.Second)
			err = client.ping()
		}

		if err != nil {
			panic(fmt.Sprintf("API is not ready after %d seconds: %v", max_retries, err))
		}
	})

	// Ensure there are no payments in the server before each scenario
	s.BeforeScenario(func(interface{}) {
		for id := range client.registeredIDs {
			err := client.iDeleteThePaymentWithID(id.String())
			if err != nil {
				panic("Error deleting registered payments")
			}
		}
	})
}
