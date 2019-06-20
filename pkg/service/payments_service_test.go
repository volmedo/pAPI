// +build integration

package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sort"
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

var (
	ps          *service.PaymentsService
	testRepo    *service.DBPaymentRepository
	testPayment models.Payment
)

func init() {
	id := strfmt.UUID("4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43")
	version := int64(0)
	orgID := strfmt.UUID("743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb")
	procDate, _ := time.Parse(strfmt.RFC3339FullDate, "2017-01-18")

	testPayment = models.Payment{
		Type:           "",
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

func TestMain(m *testing.M) {
	var dbHost, dbUser, dbPass, dbName, migrationsPath string
	var dbPort int
	flag.StringVar(&dbHost, "dbhost", "localhost", "Address of the server that hosts the DB")
	flag.IntVar(&dbPort, "dbport", 5432, "Port where the DB server is listening for connections")
	flag.StringVar(&dbUser, "dbuser", "postgres", "User to use when accessing the DB")
	flag.StringVar(&dbPass, "dbpass", "postgres", "Password to use when accessing the DB")
	flag.StringVar(&dbName, "dbname", "postgres", "Name of the DB to connect to")
	flag.StringVar(&migrationsPath, "migrations", "./migrations", "Path to the folder that contains the migration files")

	flag.Parse()

	// Setup DB
	dbConf := &service.DBConfig{
		Host:           dbHost,
		Port:           dbPort,
		User:           dbUser,
		Pass:           dbPass,
		Name:           dbName,
		MigrationsPath: migrationsPath,
	}
	db, err := service.NewDB(dbConf)
	if err != nil {
		panic(fmt.Sprintf("Unable to configure DB connection: %v", err))
	}

	testRepo, err = service.NewDBPaymentRepository(db, dbName, migrationsPath)
	if err != nil {
		panic(fmt.Sprintf("Unable to create test DB repo: %v", err))
	}

	ps = &service.PaymentsService{Repo: testRepo}

	// Run tests
	exitCode := m.Run()

	// Teardown DB
	if err := testRepo.Close(); err != nil {
		panic(fmt.Sprintf("Error closing test DB repo: %v", err))
	}

	os.Exit(exitCode)
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
			// Arrange
			if err := testRepo.DeleteAll(); err != nil {
				t.Fatalf("Error cleaning test repository: %v", err)
			}

			if tc.setupData != nil {
				for _, payment := range tc.setupData {
					_, err := testRepo.Add(payment)
					if err != nil {
						t.Fatalf("Error populating test repository: %v", err)
					}
				}
			}

			// Act
			rr, err := doRequest(ps, tc.params)
			if err != nil {
				t.Fatalf(err.Error())
			}

			// Assert
			if rr.Code != tc.wantCode {
				t.Fatalf("Wrong status code: got %v, want %v", rr.Code, tc.wantCode)
			}

			if tc.wantResp == nil {
				return
			}

			dataDiff, linksDiff, err := compareResponses(rr.Body, tc.wantResp)
			if err != nil {
				t.Fatal(err.Error())
			}
			if dataDiff != "" {
				t.Fatalf("Payment data mismatch:\n%s", dataDiff)
			}
			if linksDiff != "" {
				t.Fatalf("Link objects mismatch:\n%s", linksDiff)
			}
		})
	}
}

// copyPayment performs a deep copy of a models.Payment structure
func copyPayment(payment *models.Payment) *models.Payment {
	dup, _ := copystructure.Copy(*payment)
	paymentDup := dup.(models.Payment)
	return &paymentDup
}

func doRequest(ps *service.PaymentsService, params interface{}) (*httptest.ResponseRecorder, error) {
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

func compareResponses(gotBody io.Reader, wantResp interface{}) (dataDiff, linksDiff string, e error) {
	decoder := json.NewDecoder(gotBody)
	// Use maps to allow direct comparison, independent of element order
	wantData := make(map[strfmt.UUID]*models.Payment)
	var wantLinks *models.Links
	gotData := make(map[strfmt.UUID]*models.Payment)
	var gotLinks *models.Links
	var err error
	switch resp := wantResp.(type) {
	case *models.PaymentCreationResponse:
		wantData[*resp.Data.ID] = resp.Data
		wantLinks = resp.Links
		var gotResp models.PaymentCreationResponse
		err = decoder.Decode(&gotResp)
		if err != nil {
			return "", "", fmt.Errorf("Malformed JSON in response: %v", err)
		}
		gotData[*gotResp.Data.ID] = gotResp.Data
		gotLinks = gotResp.Links

	case *models.PaymentDetailsResponse:
		wantData[*resp.Data.ID] = resp.Data
		wantLinks = resp.Links
		var gotResp models.PaymentCreationResponse
		err = decoder.Decode(&gotResp)
		if err != nil {
			return "", "", fmt.Errorf("Malformed JSON in response: %v", err)
		}
		gotData[*gotResp.Data.ID] = gotResp.Data
		gotLinks = gotResp.Links

	case *models.PaymentUpdateResponse:
		wantData[*resp.Data.ID] = resp.Data
		wantLinks = resp.Links
		var gotResp models.PaymentCreationResponse
		err = decoder.Decode(&gotResp)
		if err != nil {
			return "", "", fmt.Errorf("Malformed JSON in response: %v", err)
		}
		gotData[*gotResp.Data.ID] = gotResp.Data
		gotLinks = gotResp.Links

	case *models.PaymentDetailsListResponse:
		for _, payment := range resp.Data {
			wantData[*payment.ID] = payment
		}
		wantLinks = resp.Links
		var gotResp models.PaymentDetailsListResponse
		err = decoder.Decode(&gotResp)
		if err != nil {
			return "", "", fmt.Errorf("Malformed JSON in response: %v", err)
		}
		for _, payment := range gotResp.Data {
			gotData[*payment.ID] = payment
		}
		gotLinks = gotResp.Links

	default:
		return "", "", fmt.Errorf("Unable to decode response, unkwnown type: %T", resp)
	}

	// go-cmp requires a custom comparer for strfmt.Date because it has unexported fields
	// see https://godoc.org/github.com/google/go-cmp/cmp/cmpopts#IgnoreUnexported
	dateComparer := cmp.Comparer(func(d1, d2 strfmt.Date) bool {
		return d1.String() == d2.String()
	})
	dataDiff = cmp.Diff(gotData, wantData, dateComparer)
	linksDiff = cmp.Diff(gotLinks, wantLinks)

	return dataDiff, linksDiff, nil
}

func createTests() []TestCase {
	req := httptest.NewRequest("POST", "/payments", nil)
	setupData := []*models.Payment{&testPayment}
	params := payments.CreatePaymentParams{
		HTTPRequest: 			req,
		PaymentCreationRequest: &models.PaymentCreationRequest{Data: &testPayment},
	}
	wantPayment := copyPayment(&testPayment)
	wantPayment.Type = service.TYPE_PAYMENT
	wantLinks := &models.Links{
		Self: fmt.Sprintf("/payments/%s", testPayment.ID),
	}
	wantResp := &models.PaymentCreationResponse{
		Data:  wantPayment,
		Links: wantLinks,
	}
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
	req := httptest.NewRequest("GET", fmt.Sprintf("/payments/%s", testPayment.ID), nil)
	setupData := []*models.Payment{&testPayment}
	params := payments.GetPaymentParams{
		HTTPRequest: req,
		ID: 		 *testPayment.ID,
	}
	wantPayment := copyPayment(&testPayment)
	wantPayment.Type = service.TYPE_PAYMENT
	wantLinks := &models.Links{
		Self: fmt.Sprintf("/payments/%s", testPayment.ID),
	}
	wantResp := &models.PaymentDetailsResponse{
		Data:  wantPayment,
		Links: wantLinks,
	}
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
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/payments/%s", testPayment.ID), nil)
	setupData := []*models.Payment{&testPayment}
	params := payments.DeletePaymentParams{
		HTTPRequest: req,
		ID: 		 *testPayment.ID,
	}
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
	req := httptest.NewRequest("PUT", fmt.Sprintf("/payments/%s", testPayment.ID), nil)
	setupData := []*models.Payment{&testPayment}

	updatedPayment := copyPayment(&testPayment)
	updatedPayment.Attributes.Amount = models.Amount("150.00")
	updatedPayment.Attributes.Fx.OriginalAmount = models.Amount("300.00")
	updatedPayment.Attributes.PaymentID = "123456789012345679"
	params := payments.UpdatePaymentParams{
		HTTPRequest:		  req,
		ID:                   *updatedPayment.ID,
		PaymentUpdateRequest: &models.PaymentUpdateRequest{Data: updatedPayment},
	}

	wantPayment := copyPayment(updatedPayment)
	wantPayment.Type = service.TYPE_PAYMENT
	wantVersion := *wantPayment.Version
	wantVersion += 1
	wantPayment.Version = &wantVersion
	wantLinks := &models.Links{
		Self: fmt.Sprintf("/payments/%s", testPayment.ID),
	}
	wantResp := &models.PaymentUpdateResponse{
		Data:  wantPayment,
		Links: wantLinks,
	}
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
	setupPaymentNum := 20

	// Sorting setup data slice simplifies expressing wanted responses
	setupIDs := make([]string, 0, setupPaymentNum)
	for i := 0; i < setupPaymentNum; i++ {
		newID, _ := uuid.NewV4()
		setupIDs = append(setupIDs, newID.String())
	}
	sort.Strings(setupIDs)

	setupData := make([]*models.Payment, 0, len(setupIDs))
	for _, id := range setupIDs {
		payment := copyPayment(&testPayment)
		pID := strfmt.UUID(id)
		payment.ID = &pID
		payment.Type = service.TYPE_PAYMENT
		setupData = append(setupData, payment)
	}

	wantLinks := &models.Links{
		Self:  "/payments",
	}

	newParams := func(pNum, pSize *int64) payments.ListPaymentsParams {
		params := payments.NewListPaymentsParams()
		if pNum != nil {
			params.PageNumber = pNum
		}
		if pSize != nil {
			params.PageSize = pSize
		}

		return params
	}

	// noParams
	// pageNumber defaults to 0 and pageSize defaults to 10 according to spec
	params := newParams(nil, nil)
	noParams := TestCase{
		name:      "list",
		setupData: setupData,
		params:    params,
		wantCode:  http.StatusOK,
		wantResp:  &models.PaymentDetailsListResponse{
			Data: setupData[:10],
			Links: wantLinks,
		},
	}

	pageSize := new(int64)
	*pageSize = 5
	// pageNumber defaults to 0 according to spec
	params = newParams(nil, pageSize)
	firstFive := TestCase{
		name:      "list first five results",
		setupData: setupData,
		params:    params,
		wantCode:  http.StatusOK,
		wantResp:  &models.PaymentDetailsListResponse{
			Data: setupData[:5],
			Links: wantLinks,
		},
	}

	pageNumber := new(int64)
	*pageNumber = 3
	pageSize = new(int64)
	*pageSize = 3
	params = newParams(pageNumber, pageSize)
	from9To11 := TestCase{
		name:      "list results from 9 to 11",
		setupData: setupData,
		params:    params,
		wantCode:  http.StatusOK,
		wantResp:  &models.PaymentDetailsListResponse{
			Data: setupData[9:12],
			Links: wantLinks,
		},
	}

	pageNumber = new(int64)
	*pageNumber = 3
	pageSize = new(int64)
	*pageSize = 6
	params = newParams(pageNumber, pageSize)
	lastPage := TestCase{
		name:      "list last page with remaining elements",
		setupData: setupData,
		params:    params,
		wantCode:  http.StatusOK,
		wantResp:  &models.PaymentDetailsListResponse{
			Data: setupData[18:],
			Links: wantLinks,
		},
	}

	pageNumber = new(int64)
	*pageNumber = 1
	// pageSize defaults to 10 according to spec
	params = newParams(pageNumber, nil)
	pageNumberButNoPageSize := TestCase{
		name:      "list with page number 2 and no page size",
		setupData: setupData,
		params:    params,
		wantCode:  http.StatusOK,
		wantResp:  &models.PaymentDetailsListResponse{
			Data: setupData[10:],
			Links: wantLinks,
		},
	}

	pageNumber = new(int64)
	*pageNumber = 10
	pageSize = new(int64)
	*pageSize = 5
	params = newParams(pageNumber, pageSize)
	paginationOffLimits := TestCase{
		name:      "list a resource beyond the limit",
		setupData: setupData,
		params:    params,
		wantCode:  http.StatusNotFound,
		wantResp:  nil,
	}

	return []TestCase{
		noParams,
		firstFive,
		from9To11,
		lastPage,
		pageNumberButNoPageSize,
		paginationOffLimits,
	}
}
