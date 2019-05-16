package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/mitchellh/copystructure"

	"github.com/volmedo/pAPI/pkg/models"
	"github.com/volmedo/pAPI/pkg/restapi/operations/payments"
	"github.com/volmedo/pAPI/pkg/service"
)

var testPayment models.Payment

func init() {
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

	// Configuration for copystructure package to correctly copy strfmt.Date
	// Copy operation on this type fails if a custom copier function is not provided
	// because of the strfmt.RFC3339FullDate custom format
	dateCopier := func(d interface{}) (interface{}, error) {
		date, ok := d.(strfmt.Date)
		if !ok {
			return nil, fmt.Errorf("Wrong type: %T", d)
		}

		dup := date.DeepCopy()
		return *dup, nil
	}
	copystructure.Copiers[reflect.TypeOf(strfmt.Date{})] = dateCopier
}

func TestCreatePayment(t *testing.T) {
	testRepo := service.NewPaymentRepository()
	ps := service.PaymentsService{Repo: testRepo}

	params := payments.CreatePaymentParams{
		PaymentCreationRequest: &models.PaymentCreationRequest{Data: &testPayment},
	}

	rr, err := doRequest(ps, params)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if rr.Code != http.StatusCreated {
		t.Errorf("Wrong status code: got %v, expected %v", rr.Code, http.StatusCreated)
	}

	var respBody models.PaymentCreationResponse
	decoder := json.NewDecoder(rr.Body)
	err = decoder.Decode(&respBody)
	if err != nil {
		t.Errorf("Malformed JSON in response: %v", err)
	}

	if !reflect.DeepEqual(testPayment, *respBody.Data) {
		t.Fatal("Payment data in request and response don't match")
	}
}

func TestCreateConflictingPayment(t *testing.T) {
	payment, _ := copyPayment(&testPayment)
	testRepo := service.NewPaymentRepository()
	err := testRepo.Add(payment)
	if err != nil {
		t.Fatal("Error populating test repository")
	}
	ps := service.PaymentsService{Repo: testRepo}

	params := payments.CreatePaymentParams{
		PaymentCreationRequest: &models.PaymentCreationRequest{Data: &testPayment},
	}
	rr, err := doRequest(ps, params)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if rr.Code != http.StatusConflict {
		t.Errorf("Wrong status code: got %v, expected %v", rr.Code, http.StatusConflict)
	}
}

func TestGetPayment(t *testing.T) {
	payment, _ := copyPayment(&testPayment)
	testRepo := service.NewPaymentRepository()
	err := testRepo.Add(payment)
	if err != nil {
		t.Fatal("Error populating test repository")
	}
	ps := service.PaymentsService{Repo: testRepo}

	pID := testPayment.ID.DeepCopy()
	params := payments.GetPaymentParams{ID: *pID}
	rr, err := doRequest(ps, params)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if rr.Code != http.StatusOK {
		t.Errorf("Wrong status code: got %v, expected %v", rr.Code, http.StatusOK)
	}

	var respBody models.PaymentDetailsResponse
	decoder := json.NewDecoder(rr.Body)
	err = decoder.Decode(&respBody)
	if err != nil {
		t.Errorf("Malformed JSON in response: %v", err)
	}

	if !reflect.DeepEqual(testPayment, *respBody.Data) {
		t.Fatal("Payment data in request and response don't match")
	}
}

func TestGetNonExistentPayment(t *testing.T) {
	testRepo := service.NewPaymentRepository()
	ps := service.PaymentsService{Repo: testRepo}

	pID := testPayment.ID.DeepCopy()
	params := payments.GetPaymentParams{ID: *pID}
	rr, err := doRequest(ps, params)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if rr.Code != http.StatusNotFound {
		t.Errorf("Wrong status code: got %v, expected %v", rr.Code, http.StatusNotFound)
	}
}

func TestDeletePayment(t *testing.T) {
	testRepo := service.NewPaymentRepository()
	err := testRepo.Add(&testPayment)
	if err != nil {
		t.Fatal("Error populating test repository")
	}
	ps := service.PaymentsService{Repo: testRepo}

	pID := testPayment.ID.DeepCopy()
	params := payments.DeletePaymentParams{ID: *pID}
	rr, err := doRequest(ps, params)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if rr.Code != http.StatusNoContent {
		t.Errorf("Wrong status code: got %v, expected %v", rr.Code, http.StatusNoContent)
	}
}

func TestDeleteNonExistentPayment(t *testing.T) {
	testRepo := service.NewPaymentRepository()
	ps := service.PaymentsService{Repo: testRepo}

	pID := testPayment.ID.DeepCopy()
	params := payments.DeletePaymentParams{ID: *pID}
	rr, err := doRequest(ps, params)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if rr.Code != http.StatusNotFound {
		t.Errorf("Wrong status code: got %v, expected %v", rr.Code, http.StatusNotFound)
	}
}

func TestUpdatePayment(t *testing.T) {
	payment, _ := copyPayment(&testPayment)
	testRepo := service.NewPaymentRepository()
	err := testRepo.Add(payment)
	if err != nil {
		t.Fatal("Error populating test repository")
	}
	ps := service.PaymentsService{Repo: testRepo}

	updatedPayment, _ := copyPayment(payment)
	updatedPayment.Attributes.Amount = models.Amount("150.00")
	updatedPayment.Attributes.Fx.OriginalAmount = models.Amount("300.00")
	updatedPayment.Attributes.PaymentID = "123456789012345679"
	params := payments.UpdatePaymentParams{
		ID:                   *updatedPayment.ID,
		PaymentUpdateRequest: &models.PaymentUpdateRequest{Data: updatedPayment},
	}
	rr, err := doRequest(ps, params)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if rr.Code != http.StatusOK {
		t.Errorf("Wrong status code: got %v, expected %v", rr.Code, http.StatusOK)
	}

	var respBody models.PaymentUpdateResponse
	decoder := json.NewDecoder(rr.Body)
	err = decoder.Decode(&respBody)
	if err != nil {
		t.Errorf("Malformed JSON in response: %v", err)
	}

	if !reflect.DeepEqual(*updatedPayment, *respBody.Data) {
		t.Fatal("Payment data in request and response don't match")
	}
}

func TestUpdateNonExistentPayment(t *testing.T) {
	testRepo := service.NewPaymentRepository()
	ps := service.PaymentsService{Repo: testRepo}

	updatedPayment, _ := copyPayment(&testPayment)
	updatedPayment.Attributes.Amount = models.Amount("150.00")
	updatedPayment.Attributes.Fx.OriginalAmount = models.Amount("300.00")
	updatedPayment.Attributes.PaymentID = "123456789012345679"
	params := payments.UpdatePaymentParams{
		ID:                   *updatedPayment.ID,
		PaymentUpdateRequest: &models.PaymentUpdateRequest{Data: updatedPayment},
	}
	rr, err := doRequest(ps, params)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if rr.Code != http.StatusNotFound {
		t.Errorf("Wrong status code: got %v, expected %v", rr.Code, http.StatusNotFound)
	}
}

func copyPayment(payment *models.Payment) (*models.Payment, error) {
	dup, err := copystructure.Copy(*payment)
	if err != nil {
		return nil, err
	}
	paymentDup, ok := dup.(models.Payment)
	if !ok {
		return nil, errors.New("Error copying payment")
	}

	return &paymentDup, nil
}

func doRequest(ps service.PaymentsService, params interface{}) (*httptest.ResponseRecorder, error) {
	ctx := context.Background()

	var responder middleware.Responder
	switch p := params.(type) {
	case payments.CreatePaymentParams:
		responder = ps.CreatePayment(ctx, p)
	case payments.GetPaymentParams:
		responder = ps.GetPayment(ctx, p)
	case payments.DeletePaymentParams:
		responder = ps.DeletePayment(ctx, p)
	case payments.UpdatePaymentParams:
		responder = ps.UpdatePayment(ctx, p)
	default:
		return nil, fmt.Errorf("Unknown params type: %T", p)
	}

	if responder == nil {
		return nil, errors.New("The returned responder should not be nil")
	}

	rr := httptest.NewRecorder()
	responder.WriteResponse(rr, runtime.JSONProducer())

	return rr, nil
}
