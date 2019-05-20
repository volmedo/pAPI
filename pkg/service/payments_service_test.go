package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/google/go-cmp/cmp"
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

type TestCase struct {
	name      string
	setupData []*models.Payment
	params    interface{}
	wantCode  int
	wantResp  interface{}
}

func TestPaymentsService(t *testing.T) {
	tests := []TestCase{}

	tests = append(tests, createTests()...)
	tests = append(tests, getTests()...)
	tests = append(tests, deleteTests()...)
	tests = append(tests, updateTests()...)
	tests = append(tests, listTests()...)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testRepo := service.NewPaymentRepository()
			if tc.setupData != nil {
				for _, payment := range tc.setupData {
					err := testRepo.Add(payment)
					if err != nil {
						t.Fatalf("Error populating test repository: %v", err)
					}
				}
			}
			ps := service.PaymentsService{Repo: testRepo}

			rr, err := doRequest(ps, tc.params)
			if err != nil {
				t.Fatalf(err.Error())
			}

			if rr.Code != tc.wantCode {
				t.Errorf("Wrong status code: got %v, want %v", rr.Code, tc.wantCode)
			}

			if tc.wantResp == nil {
				return
			}

			diff, err := compareResponses(rr.Body, tc.wantResp)
			if err != nil {
				t.Fatal(err.Error())
			}
			if diff != "" {
				t.Fatalf("Payment data mismatch:\n%s", diff)
			}
		})
	}
}

func createTests() []TestCase {
	setupData := []*models.Payment{&testPayment}
	params := payments.CreatePaymentParams{
		PaymentCreationRequest: &models.PaymentCreationRequest{Data: &testPayment},
	}
	wantResp := &models.PaymentCreationResponse{Data: &testPayment}
	return []TestCase{
		{
			name:      "create",
			setupData: nil,
			params:    params,
			wantCode:  http.StatusCreated,
			wantResp:  wantResp,
		}, {
			name:      "create conflict",
			setupData: setupData,
			params:    params,
			wantCode:  http.StatusConflict,
			wantResp:  nil,
		},
	}
}

func getTests() []TestCase {
	setupData := []*models.Payment{&testPayment}
	params := payments.GetPaymentParams{ID: *testPayment.ID}
	wantResp := &models.PaymentDetailsResponse{Data: &testPayment}
	return []TestCase{
		{
			name:      "get",
			setupData: setupData,
			params:    params,
			wantCode:  http.StatusOK,
			wantResp:  wantResp,
		}, {
			name:      "get non-existent",
			setupData: nil,
			params:    params,
			wantCode:  http.StatusNotFound,
			wantResp:  nil,
		},
	}
}

func deleteTests() []TestCase {
	setupData := []*models.Payment{&testPayment}
	params := payments.DeletePaymentParams{ID: *testPayment.ID}
	return []TestCase{
		{
			name:      "delete",
			setupData: setupData,
			params:    params,
			wantCode:  http.StatusNoContent,
			wantResp:  nil,
		}, {
			name:      "delete non-existent",
			setupData: nil,
			params:    params,
			wantCode:  http.StatusNotFound,
			wantResp:  nil,
		},
	}
}

func updateTests() []TestCase {
	setupData := []*models.Payment{&testPayment}
	updatedPayment := copyPayment(&testPayment)
	updatedPayment.Attributes.Amount = models.Amount("150.00")
	updatedPayment.Attributes.Fx.OriginalAmount = models.Amount("300.00")
	updatedPayment.Attributes.PaymentID = "123456789012345679"
	params := payments.UpdatePaymentParams{
		ID:                   *updatedPayment.ID,
		PaymentUpdateRequest: &models.PaymentUpdateRequest{Data: updatedPayment},
	}
	wantResp := &models.PaymentUpdateResponse{Data: updatedPayment}
	return []TestCase{
		{
			name:      "update",
			setupData: setupData,
			params:    params,
			wantCode:  http.StatusOK,
			wantResp:  wantResp,
		}, {
			name:      "update non-existent",
			setupData: nil,
			params:    params,
			wantCode:  http.StatusNotFound,
			wantResp:  nil,
		},
	}
}

func listTests() []TestCase {
	setupData := make([]*models.Payment, 0)
	for i := 0; i < 20; i++ {
		payment := copyPayment(&testPayment)
		newID, _ := uuid.NewV4()
		pID := strfmt.UUID(newID.String())
		payment.ID = &pID
		setupData = append(setupData, payment)
	}
	params := payments.ListPaymentsParams{}
	wantResp := &models.PaymentDetailsListResponse{Data: setupData}
	return []TestCase{
		{
			name:      "list",
			setupData: setupData,
			params:    params,
			wantCode:  http.StatusOK,
			wantResp:  wantResp,
		},
	}
}

func copyPayment(payment *models.Payment) *models.Payment {
	dup, _ := copystructure.Copy(*payment)
	paymentDup := dup.(models.Payment)
	return &paymentDup
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

	case payments.ListPaymentsParams:
		responder = ps.ListPayments(ctx, p)

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

func compareResponses(body io.Reader, wantResp interface{}) (string, error) {
	decoder := json.NewDecoder(body)
	// Use maps to allow direct comparison, independent of element order
	want := make(map[strfmt.UUID]*models.Payment)
	got := make(map[strfmt.UUID]*models.Payment)
	var err error
	switch resp := wantResp.(type) {
	case *models.PaymentCreationResponse:
		want[*resp.Data.ID] = resp.Data
		var gotResp models.PaymentCreationResponse
		err = decoder.Decode(&gotResp)
		if err != nil {
			return "", fmt.Errorf("Malformed JSON in response: %v", err)
		}
		got[*gotResp.Data.ID] = gotResp.Data

	case *models.PaymentDetailsResponse:
		want[*resp.Data.ID] = resp.Data
		var gotResp models.PaymentCreationResponse
		err = decoder.Decode(&gotResp)
		if err != nil {
			return "", fmt.Errorf("Malformed JSON in response: %v", err)
		}
		got[*gotResp.Data.ID] = gotResp.Data

	case *models.PaymentUpdateResponse:
		want[*resp.Data.ID] = resp.Data
		var gotResp models.PaymentCreationResponse
		err = decoder.Decode(&gotResp)
		if err != nil {
			return "", fmt.Errorf("Malformed JSON in response: %v", err)
		}
		got[*gotResp.Data.ID] = gotResp.Data

	case *models.PaymentDetailsListResponse:
		for _, payment := range resp.Data {
			want[*payment.ID] = payment
		}
		var gotResp models.PaymentDetailsListResponse
		err = decoder.Decode(&gotResp)
		if err != nil {
			return "", fmt.Errorf("Malformed JSON in response: %v", err)
		}
		for _, payment := range gotResp.Data {
			got[*payment.ID] = payment
		}

	default:
		return "", fmt.Errorf("Unable to decode response, unkwnown type: %T", resp)
	}

	// go-cmp requires a custom comparer for strfmt.Date because it has unexported fields
	// see https://godoc.org/github.com/google/go-cmp/cmp/cmpopts#IgnoreUnexported
	dateComparer := cmp.Comparer(func(d1, d2 strfmt.Date) bool {
		return d1.String() == d2.String()
	})
	diff := cmp.Diff(got, want, dateComparer)

	return diff, nil
}
