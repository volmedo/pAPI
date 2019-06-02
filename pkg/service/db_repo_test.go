// +build !integration

package service

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"

	"github.com/volmedo/pAPI/pkg/models"
)

var dbColumns = []string{
	"id",
	"organisation",
	"version",
	"amount",
	"beneficiary_party.name",
	"beneficiary_party.number",
	"beneficiary_party.number_code",
	"beneficiary_party.type",
	"beneficiary_party.address",
	"beneficiary_party.bank_id",
	"beneficiary_party.bank_id_code",
	"beneficiary_party.client_name",
	"charges_info.bearer_code",
	"charges_info.receiver_charges.amount",
	"charges_info.receiver_charges.currency",
	"charges_info.sender_charges",
	"currency",
	"debtor_party.name",
	"debtor_party.number",
	"debtor_party.number_code",
	"debtor_party.type",
	"debtor_party.address",
	"debtor_party.bank_id",
	"debtor_party.bank_id_code",
	"debtor_party.client_name",
	"e2e_reference",
	"fx.contract_ref",
	"fx.rate",
	"fx.original_amount.amount",
	"fx.original_amount.currency",
	"numeric_reference",
	"payment_id",
	"payment_type",
	"processing_date",
	"purpose",
	"reference",
	"scheme",
	"scheme_payment_subtype",
	"scheme_payment_type",
	"sponsor_party.account_number",
	"sponsor_party.bank_id",
	"sponsor_party.bank_id_code",
}

func setupRepo() (*DBPaymentRepository, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, mock, fmt.Errorf("Error creating DB mock: %v", err)
	}

	testRepo, err := NewDBPaymentRepository(db, "", "")
	if err != nil {
		return nil, mock, fmt.Errorf("Unable to create test DB repo: %v", err)
	}

	return testRepo, mock, nil
}

func generateDummyPayments(howMany int) []*models.Payment {
	payments := []*models.Payment{}
	for i := 0; i < howMany; i++ {
		id, _ := uuid.NewV4()
		uuid := strfmt.UUID(id.String())
		payment := &models.Payment{
			ID:             &uuid,
			OrganisationID: new(strfmt.UUID),
			Type:           "Payment",
			Version:        new(int64),
		}
		senderCharges := []*models.ChargesInformationSenderChargesItems0{&models.ChargesInformationSenderChargesItems0{}}
		attrs := &models.PaymentAttributes{
			BeneficiaryParty:   &models.PaymentParty{},
			ChargesInformation: &models.ChargesInformation{SenderCharges: senderCharges},
			DebtorParty:        &models.PaymentParty{},
			Fx:                 &models.PaymentAttributesFx{},
			SponsorParty:       &models.PaymentAttributesSponsorParty{},
		}

		payment.Attributes = attrs

		payments = append(payments, payment)
	}

	return payments
}

func paymentsToRows(payments []*models.Payment) *sqlmock.Rows {
	rows := sqlmock.NewRows(dbColumns)
	for _, payment := range payments {
		attrs := payment.Attributes
		amounts := senderChargesToAmounts(attrs.ChargesInformation.SenderCharges)
		rows.AddRow(
			payment.ID,                                       // id,
			payment.OrganisationID,                           // organisation,
			*payment.Version,                                 // version,
			attrs.Amount,                                     // amount,
			attrs.BeneficiaryParty.AccountName,               // beneficiary_party.name,
			attrs.BeneficiaryParty.AccountNumber,             // beneficiary_party.number,
			attrs.BeneficiaryParty.AccountNumberCode,         // beneficiary_party.number_code,
			attrs.BeneficiaryParty.AccountType,               // beneficiary_party.type,
			attrs.BeneficiaryParty.Address,                   // beneficiary_party.address ,
			attrs.BeneficiaryParty.BankID,                    // beneficiary_party.bank_id,
			attrs.BeneficiaryParty.BankIDCode,                // beneficiary_party.bank_id_code,
			attrs.BeneficiaryParty.Name,                      // beneficiary_party.client_name,
			attrs.ChargesInformation.BearerCode,              // charges_info.bearer_code,
			attrs.ChargesInformation.ReceiverChargesAmount,   // charges_info.receiver_charges.amount,
			attrs.ChargesInformation.ReceiverChargesCurrency, // charges_info.receiver_charges.currency,
			pq.Array(amounts),                                // charges_info.sender_charges,
			attrs.Currency,                                   // currency,
			attrs.DebtorParty.AccountName,                    // debtor_party.name,
			attrs.DebtorParty.AccountNumber,                  // debtor_party.number,
			attrs.DebtorParty.AccountNumberCode,              // debtor_party.number_code,
			attrs.DebtorParty.AccountType,                    // debtor_party.type,
			attrs.DebtorParty.Address,                        // debtor_party.address ,
			attrs.DebtorParty.BankID,                         // debtor_party.bank_id,
			attrs.DebtorParty.BankIDCode,                     // debtor_party.bank_id_code,
			attrs.DebtorParty.Name,                           // debtor_party.client_name,
			attrs.EndToEndReference,                          // e2e_reference,
			attrs.Fx.ContractReference,                       // fx.contract_ref,
			attrs.Fx.ExchangeRate,                            // fx.rate,
			attrs.Fx.OriginalAmount,                          // fx.original_amount.amount,
			attrs.Fx.OriginalCurrency,                        // fx.original_amount.currency,
			attrs.NumericReference,                           // numeric_reference,
			attrs.PaymentID,                                  // payment_id,
			attrs.PaymentType,                                // payment_type,
			attrs.ProcessingDate,                             // processing_date,
			attrs.PaymentPurpose,                             // purpose,
			attrs.Reference,                                  // reference,
			attrs.PaymentScheme,                              // scheme,
			attrs.SchemePaymentSubType,                       // scheme_payment_subtype,
			attrs.SchemePaymentType,                          // scheme_payment_type,
			attrs.SponsorParty.AccountNumber,                 // sponsor_party.account_number,
			attrs.SponsorParty.BankID,                        // sponsor_party.bank_id,
			attrs.SponsorParty.BankIDCode,                    // sponsor_party.bank_id_code
		)
	}

	return rows
}

func TestAdd(t *testing.T) {
	testRepo, mock, err := setupRepo()
	if err != nil {
		t.Fatalf("Error setting up test repo")
	}
	defer testRepo.Close()

	mock.ExpectExec(`^INSERT INTO payments`).WillReturnResult(sqlmock.NewResult(0, 1))

	testPayment := generateDummyPayments(1)[0]
	if err := testRepo.Add(testPayment); err != nil {
		t.Errorf("Unexpected error adding payment: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestAddConflict(t *testing.T) {
	testRepo, mock, err := setupRepo()
	if err != nil {
		t.Fatalf("Error setting up test repo")
	}
	defer testRepo.Close()

	mock.ExpectExec(`^INSERT INTO payments`).WillReturnResult(sqlmock.NewResult(0, 1)).WillReturnError(&pq.Error{Code: pq.ErrorCode("23505")})

	testPayment := generateDummyPayments(1)[0]
	err = testRepo.Add(testPayment)
	e, ok := err.(ErrConflict)
	if err == nil || !ok {
		t.Errorf("Expected ErrConflict but got %v", e)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestDelete(t *testing.T) {
	testRepo, mock, err := setupRepo()
	if err != nil {
		t.Fatalf("Error setting up test repo")
	}
	defer testRepo.Close()

	testPayment := generateDummyPayments(1)[0]
	mock.ExpectExec(`DELETE FROM payments WHERE id = \$1`).
		WithArgs(*testPayment.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := testRepo.Delete(*testPayment.ID); err != nil {
		t.Errorf("Unexpected error deleting payment: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestDeleteNonExistent(t *testing.T) {
	testRepo, mock, err := setupRepo()
	if err != nil {
		t.Fatalf("Error setting up test repo")
	}
	defer testRepo.Close()

	testPayment := generateDummyPayments(1)[0]
	mock.ExpectExec(`DELETE FROM payments WHERE id = \$1`).
		WithArgs(*testPayment.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = testRepo.Delete(*testPayment.ID)
	e, ok := err.(ErrNoResults)
	if err == nil || !ok {
		t.Errorf("Expected ErrNoResults but got %v", e)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestGet(t *testing.T) {
	testRepo, mock, err := setupRepo()
	if err != nil {
		t.Fatalf("Error setting up test repo")
	}
	defer testRepo.Close()

	testPayment := generateDummyPayments(1)[0]
	rows := paymentsToRows([]*models.Payment{testPayment})
	mock.ExpectQuery(`^SELECT (.+) FROM payments WHERE id = \$1$`).
		WithArgs(*testPayment.ID).
		WillReturnRows(rows)

	if _, err := testRepo.Get(*testPayment.ID); err != nil {
		t.Errorf("Error getting payment: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestGetNonExistent(t *testing.T) {
	testRepo, mock, err := setupRepo()
	if err != nil {
		t.Fatalf("Error setting up test repo")
	}
	defer testRepo.Close()

	testPayment := generateDummyPayments(1)[0]
	mock.ExpectQuery(`^SELECT (.+) FROM payments WHERE id = \$1$`).
		WithArgs(*testPayment.ID).
		WillReturnError(sql.ErrNoRows)

	_, err = testRepo.Get(*testPayment.ID)
	e, ok := err.(ErrNoResults)
	if err == nil || !ok {
		t.Errorf("Expected ErrNoResults but got %v", e)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestList(t *testing.T) {
	testPayments := generateDummyPayments(200)

	tests := map[string]struct {
		offset      int
		limit       int
		shouldFail  bool
		expectedLen int
	}{
		"first 5": {
			offset:      0,
			limit:       5,
			shouldFail:  false,
			expectedLen: 5,
		},
		"from 25 to 32": {
			offset:      25,
			limit:       6,
			shouldFail:  false,
			expectedLen: 6,
		},
		"less than limit available": {
			offset:      len(testPayments) - 10,
			limit:       20,
			shouldFail:  false,
			expectedLen: 10,
		},
		"offset negative": {
			offset:      -1,
			limit:       5,
			shouldFail:  true,
			expectedLen: 0,
		},
		"offset too high": {
			offset:      len(testPayments),
			limit:       5,
			shouldFail:  true,
			expectedLen: 0,
		},
		"limit 0": {
			offset:      0,
			limit:       0,
			shouldFail:  true,
			expectedLen: 0,
		},
		"limit negative": {
			offset:      0,
			limit:       -1,
			shouldFail:  true,
			expectedLen: 0,
		},
		"limit too high": {
			offset:      0,
			limit:       101,
			shouldFail:  true,
			expectedLen: 0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			testRepo, mock, err := setupRepo()
			if err != nil {
				t.Fatalf("Error setting up test repo")
			}
			defer testRepo.Close()

			mock.ExpectQuery(`SELECT COUNT\(\*\) FROM payments`).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(testPayments)))

			if tc.shouldFail {
				_, err := testRepo.List(int64(tc.offset), int64(tc.limit))
				if err == nil {
					t.Errorf("Test should've failed but no error was produced")
				} else if _, ok := err.(ErrBadOffsetLimit); !ok {
					t.Errorf("Expected ErrBadOffsetLimit but got %#v", err)
				}
				return
			}

			from := tc.offset
			to := tc.offset + tc.limit
			if to > len(testPayments) {
				to = len(testPayments)
			}
			rows := paymentsToRows(testPayments[from:to])
			mock.ExpectQuery(`^SELECT (.+) FROM payments ORDER BY id ASC LIMIT \$1 OFFSET \$2$`).
				WithArgs(tc.limit, tc.offset).
				WillReturnRows(rows)

			payments, err := testRepo.List(int64(tc.offset), int64(tc.limit))
			if err != nil {
				t.Fatalf("Unexpected error: %#v", err)
			}

			if len(payments) != tc.expectedLen {
				t.Errorf("Want %d items but got %d", tc.expectedLen, len(payments))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Expectations were not met: %s", err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	testRepo, mock, err := setupRepo()
	if err != nil {
		t.Fatalf("Error setting up test repo")
	}
	defer testRepo.Close()

	testPayment := generateDummyPayments(1)[0]
	args := make([]driver.Value, 42)
	for i := range args {
		args[i] = sqlmock.AnyArg()
	}
	mock.ExpectExec(`^UPDATE payments SET (.+) WHERE id = \$1$`).
		WithArgs(args...).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := testRepo.Update(*testPayment.ID, testPayment); err != nil {
		t.Errorf("Unexpected error updating payment: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestUpdateNonExistent(t *testing.T) {
	testRepo, mock, err := setupRepo()
	if err != nil {
		t.Fatalf("Error setting up test repo")
	}
	defer testRepo.Close()

	testPayment := generateDummyPayments(1)[0]
	args := make([]driver.Value, 42)
	for i := range args {
		args[i] = sqlmock.AnyArg()
	}
	mock.ExpectExec(`^UPDATE payments SET (.+) WHERE id = \$1$`).
		WithArgs(args...).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = testRepo.Update(*testPayment.ID, testPayment)
	e, ok := err.(ErrNoResults)
	if err == nil || !ok {
		t.Errorf("Expected ErrNoResults but got %v", e)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestAmountScan(t *testing.T) {
	tests := map[string]struct {
		input      string
		shouldFail bool
		want       amount
	}{
		"basic": {
			input:      "(5.00,USD)",
			shouldFail: false,
			want:       amount{"5.00", "USD"},
		},
		"no parens": {
			input:      "5.00,USD",
			shouldFail: false,
			want:       amount{"5.00", "USD"},
		},
		"no Lparen": {
			input:      "5.00,USD)",
			shouldFail: false,
			want:       amount{"5.00", "USD"},
		},
		"no Rparen": {
			input:      "(5.00,USD",
			shouldFail: false,
			want:       amount{"5.00", "USD"},
		},
		"quotes": {
			input:      "(\"5.00\",\"USD\")",
			shouldFail: false,
			want:       amount{"5.00", "USD"},
		},
		"quotes full": {
			input:      "(\"5.00,USD\")",
			shouldFail: false,
			want:       amount{"5.00", "USD"},
		},
		"single quotes": {
			input:      "('5.00','USD')",
			shouldFail: false,
			want:       amount{"5.00", "USD"},
		},
		"single quotes full": {
			input:      "('5.00,USD')",
			shouldFail: false,
			want:       amount{"5.00", "USD"},
		},
		"not enough elements": {
			input:      "(5.00)",
			shouldFail: true,
			want:       amount{},
		},
		"too much elements": {
			input:      "(5.00,USD,GBP)",
			shouldFail: true,
			want:       amount{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			a := &amount{}
			err := a.Scan(tc.input)

			if tc.shouldFail {
				if err == nil {
					t.Fatalf("Test should've failed but no error was produced")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %#v", err)
			}

			if *a != tc.want {
				t.Fatalf("got: %v, want: %v", a, tc.want)
			}
		})
	}
}
