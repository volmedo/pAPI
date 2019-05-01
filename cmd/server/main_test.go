package server_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/volmedo/pAPI/pkg/impl"
	"github.com/volmedo/pAPI/pkg/models"
	"github.com/volmedo/pAPI/pkg/restapi"
)

const apiRoot = "/v1"

var testPayment models.Payment

func init() {
	id := strfmt.UUID("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43")
	version := int64(0)
	orgID := strfmt.UUID("743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb")
	procDate, _ := time.Parse(time.RFC3339, "2017-01-18")

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

func configureHandler() http.Handler {
	papi := &impl.PaymentsAPI{}

	handler, err := restapi.Handler(restapi.Config{
		PaymentsAPI: papi,
		Logger:      log.Printf,
	})

	if err != nil {
		log.Fatal(err)
	}

	return handler
}

func TestCreatePayment(t *testing.T) {
	handler := configureHandler()
	reqBody, _ := json.Marshal(models.PaymentCreationRequest{Data: &testPayment})
	req, err := http.NewRequest("POST", apiRoot+"/payments", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/vnd.api+json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Fatalf("Wrong status code: got %v, expected %v", status, http.StatusCreated)
	}

	var paymentCreationResp models.PaymentCreationResponse
	decoder := json.NewDecoder(rr.Body)
	err = decoder.Decode(&paymentCreationResp)
	if err != nil {
		t.Fatalf("Invalid JSON response: %s", err)
	}
	gotPayment := *paymentCreationResp.Data

	if !reflect.DeepEqual(gotPayment, testPayment) {
		t.Error("Payment data in the response don't match that of the request")
	}
}
