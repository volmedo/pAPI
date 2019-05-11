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
	scheme   string
	host     string
	port     int
	basePath string

	testPayment models.Payment

	opt = godog.Options{
		Output: colors.Colored(os.Stdout),
		Format: "progress",
	}
)

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)

	id := strfmt.UUID("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43")
	version := int64(0)
	orgID := strfmt.UUID("743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb")
	procDate, _ := time.Parse(strfmt.RFC3339FullDate, "2017-01-18")

	testPayment = models.Payment{
		Type:           "Payment",
		ID:             &id,
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
}

type Client struct {
	*payments.Client
	lastResponse *payments.CreatePaymentCreated
	lastError    error
}

func newClient(apiURL *url.URL) *Client {
	conf := client.Config{URL: apiURL}
	payments := client.New(conf)
	return &Client{payments.Payments, nil, nil}
}

func (c *Client) iCreateANewPaymentDescribedInJSONAs(jsonPayment *gherkin.DocString) error {
	ctx := context.Background()
	req := &models.PaymentCreationRequest{Data: &testPayment}
	params := payments.NewCreatePaymentParams().WithPaymentCreationRequest(req)

	c.lastResponse, c.lastError = c.CreatePayment(ctx, params)
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
	resp := c.lastResponse
	gotPayment := *resp.Payload.Data

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
	s.Step(`^I get a "([^"]*)" response$`, client.iGetAResponse)
	s.Step(`^the response contains a payment described in JSON as:$`, client.theResponseContainsAPaymentDescribedInJSONAs)
}
