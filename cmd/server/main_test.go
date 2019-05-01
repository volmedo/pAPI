package server_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/volmedo/pAPI/pkg/restapi"
)

const (
	apiRoot         = "/v1"
	testPaymentJSON = `
	{
		"type": "Payment",
		"id": "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
		"version": 0,
		"organisation_id": "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
		"attributes": {
			"amount": "100.21",
			"beneficiary_party": {
				"account_name": "W Owens",
				"account_number": "31926819",
				"account_number_code": "BBAN",
				"account_type": 0,
				"address": "1 The Beneficiary Localtown SE2",
				"bank_id": "403000",
				"bank_id_code": "GBDSC",
				"name": "Wilfred Jeremiah Owens"
			},
			"charges_information": {
				"bearer_code": "SHAR",
				"sender_charges": [
					{ "amount": "5.00", "currency": "GBP" },
					{ "amount": "10.00", "currency": "USD" }
				],
				"receiver_charges_amount": "1.00",
				"receiver_charges_currency": "USD"
			},
			"currency": "GBP",
			"debtor_party": {
				"account_name": "EJ Brown Black",
				"account_number": "GB29XABC10161234567801",
				"account_number_code": "IBAN",
				"address": "10 Debtor Crescent Sourcetown NE1",
				"bank_id": "203301",
				"bank_id_code": "GBDSC",
				"name": "Emelia Jane Brown"
			},
			"end_to_end_reference": "Wil piano Jan",
			"fx": {
				"contract_reference": "FX123",
				"exchange_rate": "2.00000",
				"original_amount": "200.42",
				"original_currency": "USD"
			},
			"numeric_reference": "1002001",
			"payment_id": "123456789012345678",
			"payment_purpose": "Paying for goods/services",
			"payment_scheme": "FPS",
			"payment_type": "Credit",
			"processing_date": "2017-01-18",
			"reference": "Payment for Em's piano lessons",
			"scheme_payment_sub_type": "InternetBanking",
			"scheme_payment_type": "ImmediatePayment",
			"sponsor_party": {
				"account_number": "56781234",
				"bank_id": "123123",
				"bank_id_code": "GBDSC"
			}
		}
	}`
)

var handler http.Handler

func init() {
	papi := impl.PaymentsAPI{}

	// Initiate the http handler, with the objects that are implementing the business logic.
	handler, err := restapi.Handler(restapi.Config{
		PaymentsAPI: papi,
		Logger:      log.Printf,
	})
	_ = handler
	if err != nil {
		log.Fatal(err)
	}
}

func TestCreatePayment(t *testing.T) {
	req, err := http.NewRequest("POST", apiRoot+"/payments", bytes.NewBuffer([]byte(testPaymentJSON)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Wrong status code: got %v, expected %v", status, http.StatusCreated)
	}

	var gotPayment interface{}
	decoder := json.NewDecoder(rr.Body)
	err = decoder.Decode(&gotPayment)
	if err != nil {
		t.Errorf("Invalid JSON response: %s", err)
	}

	var expectedPayment interface{}
	decoder = json.NewDecoder(bytes.NewBuffer([]byte(testPaymentJSON)))
	err = decoder.Decode(&expectedPayment)
	if err != nil {
		t.Errorf("Invalid JSON string in test specification: %s", err)
	}

	if !reflect.DeepEqual(gotPayment, expectedPayment) {
		t.Fatal("Payment data in the response don't match expected payment data")
	}
}
