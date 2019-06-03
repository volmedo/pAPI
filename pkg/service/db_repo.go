package service

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/csv"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
	"github.com/mitchellh/copystructure"

	"github.com/volmedo/pAPI/pkg/models"
)

// Type amount mimics the composite type amount in the DB schema.
// It is declared here so that it can implement driver.Valuer and sql.Scanner
// interfaces, letting the driver to handle marshalling and unmarshalling
// of the type in a more natural manner.
// Other types are not mirrored here because they can be accesed directly by
// their fields. However, ChargesInformation contains an array of
// sender charges, which in turn is mapped as an amount[] in the DB.
// Being an array, it is more maintainable and less error-prone to make type
// amount implement the Valuer and Scanner interfaces and let the driver handle
// the specifics of array syntax through pq.Array
type amount struct {
	amount   string
	currency string
}

// Value satisfies driver.Valuer and allows the conversion of an
// amount to a driver.Value
func (a amount) Value() (driver.Value, error) {
	return fmt.Sprintf("(%s,%s)", a.amount, a.currency), nil
}

// Scan makes amount implement sql.Scanner. It assigns a value from a db driver
func (a *amount) Scan(raw interface{}) error {
	var s string
	switch v := raw.(type) {
	case []byte:
		s = string(v)
	case string:
		s = v
	default:
		return fmt.Errorf("Cannot sql.Scan() amount from: %#v", v)
	}

	// PostgreSQL syntax for composite types is "(field1, field2, ...)" and
	// field values can be quoted. We will strip parentheses and get rid of
	// single and double quotes and then split on commas
	s = strings.Trim(s, "()")
	s = strings.Replace(s, "\"", "", -1)
	s = strings.Replace(s, "'", "", -1)
	r := csv.NewReader(bytes.NewBuffer([]byte(s)))
	fields, err := r.Read()
	if err != nil {
		return err
	}
	if len(fields) != 2 {
		return fmt.Errorf("Expected 2 elements but got %d", len(fields))
	}

	a.amount = fields[0]
	a.currency = fields[1]

	return nil
}

// senderChargesToAmounts returns a slice of amount from a slice of SenderCharges
func senderChargesToAmounts(charges []*models.ChargesInformationSenderChargesItems0) []amount {
	amounts := make([]amount, 0)
	for _, charge := range charges {
		a := amount{amount: string(charge.Amount), currency: string(charge.Currency)}
		amounts = append(amounts, a)
	}

	return amounts
}

// amountsToSenderCharges returns a slice of SenderCharges from the data of a slice of amounts
func amountsToSenderCharges(amounts []amount) []*models.ChargesInformationSenderChargesItems0 {
	senderCharges := make([]*models.ChargesInformationSenderChargesItems0, 0)
	for _, amount := range amounts {
		charge := models.ChargesInformationSenderChargesItems0{
			Amount:   models.Amount(amount.amount),
			Currency: models.Currency(amount.currency),
		}
		senderCharges = append(senderCharges, &charge)
	}

	return senderCharges
}

// DBConfig contains the parameters needed to connect to a DB
type DBConfig struct {
	// Address of the server that hosts the DB
	Host string

	// Port where the DB server is listening for connections
	Port int

	// User to use when accessing the DB
	User string

	// Password to use when accessing the DB
	Pass string

	// Name of the DB to connect to
	Name string

	// Path to the folder that contains the migration files
	MigrationsPath string
}

// NewDB initializes a new DB connection object using the configuration paraemters provided
func NewDB(cfg *DBConfig) (*sql.DB, error) {
	if cfg.Host == "" || cfg.Port == 0 || cfg.User == "" || cfg.Name == "" {
		return nil, fmt.Errorf("Missing required parameters in config: %v", cfg)
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.Name)
	return sql.Open("postgres", connStr)
}

// DBPaymentRepository stores a collection of payment resources using
// an external database as data backend
type DBPaymentRepository struct {
	db *sql.DB
}

// NewDBPaymentRepository creates a new DBPaymentRepository that uses a previously
// configured sql.DB to connect to the DB
func NewDBPaymentRepository(db *sql.DB, dbName, migrationsPath string) (*DBPaymentRepository, error) {
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db: pinging the DB didn't work: %v", err)
	}

	if dbName != "" && migrationsPath != "" {
		if err := migrateDB(db, dbName, migrationsPath); err != nil {
			return nil, fmt.Errorf("db: migration failed: %v", err)
		}
	}

	return &DBPaymentRepository{db: db}, nil
}

// migrateDB updates the DB's schema to the latest version
func migrateDB(db *sql.DB, dbName, migrationsPath string) error {
	dbDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsPath, dbName, dbDriver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

// Close closes the underlying db instance and frees its associated resources
func (dbpr *DBPaymentRepository) Close() error {
	if dbpr.db != nil {
		if err := dbpr.db.Close(); err != nil {
			return fmt.Errorf("db: error closing underlying DB connection: %v", err)
		}
	}

	return nil
}

// Add adds a new payment resource to the repository
//
// Add returns an error if a payment with the same ID as the one
// to be added already exists
func (dbpr *DBPaymentRepository) Add(payment *models.Payment) (*models.Payment, error) {
	insertStmt := `
	INSERT INTO payments (
		id,
		organisation,
		version,
		amount,
		beneficiary_party.name,
		beneficiary_party.number,
		beneficiary_party.number_code,
		beneficiary_party.type,
		beneficiary_party.address,
		beneficiary_party.bank_id,
		beneficiary_party.bank_id_code,
		beneficiary_party.client_name,
		charges_info.bearer_code,
		charges_info.receiver_charges.amount,
		charges_info.receiver_charges.currency,
		charges_info.sender_charges,
		currency,
		debtor_party.name,
		debtor_party.number,
		debtor_party.number_code,
		debtor_party.type,
		debtor_party.address,
		debtor_party.bank_id,
		debtor_party.bank_id_code,
		debtor_party.client_name,
		e2e_reference,
		fx.contract_ref,
		fx.rate,
		fx.original_amount.amount,
		fx.original_amount.currency,
		numeric_reference,
		payment_id,
		payment_type,
		processing_date,
		purpose,
		reference,
		scheme,
		scheme_payment_subtype,
		scheme_payment_type,
		sponsor_party.account_number,
		sponsor_party.bank_id,
		sponsor_party.bank_id_code
	)
	VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
		$11, $12, $13, $14, $15, $16::amount[], $17, $18, $19, $20,
		$21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
		$31, $32, $33, $34, $35, $36, $37, $38, $39, $40,
		$41, $42
	)`

	attrs := payment.Attributes
	amounts := senderChargesToAmounts(attrs.ChargesInformation.SenderCharges)
	_, err := dbpr.db.Exec(insertStmt,
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

	if err != nil {
		if e, ok := err.(*pq.Error); ok && e.Code == "23505" {
			return nil, newErrConflict(fmt.Sprintf("db: a payment with ID %s already exists", *payment.ID))
		}

		return nil, fmt.Errorf("db: error executing insert: %v", err)
	}

	// Ignore the original type attribute and fix it to TYPE_PAYMENT
	added := copyPayment(payment)
	added.Type = TYPE_PAYMENT
	return added, nil
}

// Delete deletes the payment resource associated to the given paymentID
//
// Delete returns an error if the paymentID is not present in the respository
func (dbpr *DBPaymentRepository) Delete(paymentID strfmt.UUID) error {
	deleteStmt := `DELETE FROM payments WHERE id = $1`
	res, err := dbpr.db.Exec(deleteStmt, paymentID.String())
	if err != nil {
		return fmt.Errorf("db: error executing delete: %v", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: error getting rows affected by delete: %v", err)
	}
	if count == 0 {
		return newErrNoResults(fmt.Sprintf("db: payment with ID %s not found", paymentID))
	}

	return nil
}

// DeleteAll deletes every payment in the DB
func (dbpr *DBPaymentRepository) DeleteAll() error {
	_, err := dbpr.db.Exec(`DELETE FROM payments`)
	if err != nil {
		return fmt.Errorf("db: error executing delete: %v", err)
	}

	return nil
}

// Get returns the payment resource associated with the given paymentID
//
// Get returns an error if the paymentID does not exist in the collection
func (dbpr *DBPaymentRepository) Get(paymentID strfmt.UUID) (*models.Payment, error) {
	selectStmt := `
	SELECT
		id,
		organisation,
		version,
		amount,
		(beneficiary_party).name,
		(beneficiary_party).number,
		(beneficiary_party).number_code,
		(beneficiary_party).type,
		(beneficiary_party).address,
		(beneficiary_party).bank_id,
		(beneficiary_party).bank_id_code,
		(beneficiary_party).client_name,
		(charges_info).bearer_code,
		(charges_info).receiver_charges.amount,
		(charges_info).receiver_charges.currency,
		(charges_info).sender_charges,
		currency,
		(debtor_party).name,
		(debtor_party).number,
		(debtor_party).number_code,
		(debtor_party).type,
		(debtor_party).address,
		(debtor_party).bank_id,
		(debtor_party).bank_id_code,
		(debtor_party).client_name,
		e2e_reference,
		(fx).contract_ref,
		(fx).rate,
		(fx).original_amount.amount,
		(fx).original_amount.currency,
		numeric_reference,
		payment_id,
		payment_type,
		processing_date,
		purpose,
		reference,
		scheme,
		scheme_payment_subtype,
		scheme_payment_type,
		(sponsor_party).account_number,
		(sponsor_party).bank_id,
		(sponsor_party).bank_id_code
	FROM payments
	WHERE id = $1`

	payment := models.Payment{
		ID:             new(strfmt.UUID),
		OrganisationID: new(strfmt.UUID),
		Type:           TYPE_PAYMENT,
		Version:        new(int64),
	}
	attrs := models.PaymentAttributes{
		BeneficiaryParty:   &models.PaymentParty{},
		ChargesInformation: &models.ChargesInformation{},
		DebtorParty:        &models.PaymentParty{},
		Fx:                 &models.PaymentAttributesFx{},
		SponsorParty:       &models.PaymentAttributesSponsorParty{},
	}
	var amounts []amount

	row := dbpr.db.QueryRow(selectStmt, paymentID.String())
	err := row.Scan(
		payment.ID,                                        // id,
		payment.OrganisationID,                            // organisation,
		payment.Version,                                   // version,
		&attrs.Amount,                                     // amount,
		&attrs.BeneficiaryParty.AccountName,               // beneficiary_party.name,
		&attrs.BeneficiaryParty.AccountNumber,             // beneficiary_party.number,
		&attrs.BeneficiaryParty.AccountNumberCode,         // beneficiary_party.number_code,
		&attrs.BeneficiaryParty.AccountType,               // beneficiary_party.type,
		&attrs.BeneficiaryParty.Address,                   // beneficiary_party.address ,
		&attrs.BeneficiaryParty.BankID,                    // beneficiary_party.bank_id,
		&attrs.BeneficiaryParty.BankIDCode,                // beneficiary_party.bank_id_code,
		&attrs.BeneficiaryParty.Name,                      // beneficiary_party.client_name,
		&attrs.ChargesInformation.BearerCode,              // charges_info.bearer_code,
		&attrs.ChargesInformation.ReceiverChargesAmount,   // charges_info.receiver_charges.amount,
		&attrs.ChargesInformation.ReceiverChargesCurrency, // charges_info.receiver_charges.currency,
		pq.Array(&amounts),                                // charges_info.sender_charges,
		&attrs.Currency,                                   // currency,
		&attrs.DebtorParty.AccountName,                    // debtor_party.name,
		&attrs.DebtorParty.AccountNumber,                  // debtor_party.number,
		&attrs.DebtorParty.AccountNumberCode,              // debtor_party.number_code,
		&attrs.DebtorParty.AccountType,                    // debtor_party.type,
		&attrs.DebtorParty.Address,                        // debtor_party.address ,
		&attrs.DebtorParty.BankID,                         // debtor_party.bank_id,
		&attrs.DebtorParty.BankIDCode,                     // debtor_party.bank_id_code,
		&attrs.DebtorParty.Name,                           // debtor_party.client_name,
		&attrs.EndToEndReference,                          // e2e_reference,
		&attrs.Fx.ContractReference,                       // fx.contract_ref,
		&attrs.Fx.ExchangeRate,                            // fx.rate,
		&attrs.Fx.OriginalAmount,                          // fx.original_amount.amount,
		&attrs.Fx.OriginalCurrency,                        // fx.original_amount.currency,
		&attrs.NumericReference,                           // numeric_reference,
		&attrs.PaymentID,                                  // payment_id,
		&attrs.PaymentType,                                // payment_type,
		&attrs.ProcessingDate,                             // processing_date,
		&attrs.PaymentPurpose,                             // purpose,
		&attrs.Reference,                                  // reference,
		&attrs.PaymentScheme,                              // scheme,
		&attrs.SchemePaymentSubType,                       // scheme_payment_subtype,
		&attrs.SchemePaymentType,                          // scheme_payment_type,
		&attrs.SponsorParty.AccountNumber,                 // sponsor_party.account_number,
		&attrs.SponsorParty.BankID,                        // sponsor_party.bank_id,
		&attrs.SponsorParty.BankIDCode,                    // sponsor_party.bank_id_code
	)

	attrs.ChargesInformation.SenderCharges = amountsToSenderCharges(amounts)
	payment.Attributes = &attrs

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, newErrNoResults(fmt.Sprintf("db: payment with ID %s not found", paymentID))
		}

		return nil, fmt.Errorf("db: error executing select: %v", err)
	}

	return &payment, nil
}

// List returns a slice of payment resources. An empty slice will be returned
// if no payment exists.
//
// List implements basic pagination by means of offset and limit parameters.
// List will return an error if offset is beyond the number of elements available.
// Limit must be between 1 and 100.
func (dbpr *DBPaymentRepository) List(offset, limit int64) ([]*models.Payment, error) {
	// Check params before anything else
	if limit <= 0 || limit > 100 {
		return nil, newErrBadOffsetLimit(fmt.Sprintf("db: list limit %d is outside allowed range (0, 100]", limit))
	}

	if offset < 0 {
		return nil, newErrBadOffsetLimit(fmt.Sprintf("db: list offset %d negative", offset))
	}

	numRecords := 0
	row := dbpr.db.QueryRow(`SELECT COUNT(*) FROM payments`)
	err := row.Scan(&numRecords)
	if err != nil {
		return nil, fmt.Errorf("db: error counting records: %v", err)
	}
	if offset >= int64(numRecords) {
		return nil, newErrBadOffsetLimit(fmt.Sprintf("db: list offset is %d but only %d records exist", offset, numRecords))
	}

	listStmt := `
	SELECT
		id,
		organisation,
		version,
		amount,
		(beneficiary_party).name,
		(beneficiary_party).number,
		(beneficiary_party).number_code,
		(beneficiary_party).type,
		(beneficiary_party).address,
		(beneficiary_party).bank_id,
		(beneficiary_party).bank_id_code,
		(beneficiary_party).client_name,
		(charges_info).bearer_code,
		(charges_info).receiver_charges.amount,
		(charges_info).receiver_charges.currency,
		(charges_info).sender_charges,
		currency,
		(debtor_party).name,
		(debtor_party).number,
		(debtor_party).number_code,
		(debtor_party).type,
		(debtor_party).address,
		(debtor_party).bank_id,
		(debtor_party).bank_id_code,
		(debtor_party).client_name,
		e2e_reference,
		(fx).contract_ref,
		(fx).rate,
		(fx).original_amount.amount,
		(fx).original_amount.currency,
		numeric_reference,
		payment_id,
		payment_type,
		processing_date,
		purpose,
		reference,
		scheme,
		scheme_payment_subtype,
		scheme_payment_type,
		(sponsor_party).account_number,
		(sponsor_party).bank_id,
		(sponsor_party).bank_id_code
	FROM payments
	ORDER BY id ASC
	LIMIT $1
	OFFSET $2`

	rows, err := dbpr.db.Query(listStmt, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("db: error executing list query: %v", err)
	}
	defer rows.Close()

	payments := make([]*models.Payment, 0, limit)
	for rows.Next() {
		payment := models.Payment{
			ID:             new(strfmt.UUID),
			OrganisationID: new(strfmt.UUID),
			Type:           TYPE_PAYMENT,
			Version:        new(int64),
		}
		attrs := models.PaymentAttributes{
			BeneficiaryParty:   &models.PaymentParty{},
			ChargesInformation: &models.ChargesInformation{},
			DebtorParty:        &models.PaymentParty{},
			Fx:                 &models.PaymentAttributesFx{},
			SponsorParty:       &models.PaymentAttributesSponsorParty{},
		}
		var amounts []amount

		err := rows.Scan(
			payment.ID,                                        // id,
			payment.OrganisationID,                            // organisation,
			payment.Version,                                   // version,
			&attrs.Amount,                                     // amount,
			&attrs.BeneficiaryParty.AccountName,               // beneficiary_party.name,
			&attrs.BeneficiaryParty.AccountNumber,             // beneficiary_party.number,
			&attrs.BeneficiaryParty.AccountNumberCode,         // beneficiary_party.number_code,
			&attrs.BeneficiaryParty.AccountType,               // beneficiary_party.type,
			&attrs.BeneficiaryParty.Address,                   // beneficiary_party.address ,
			&attrs.BeneficiaryParty.BankID,                    // beneficiary_party.bank_id,
			&attrs.BeneficiaryParty.BankIDCode,                // beneficiary_party.bank_id_code,
			&attrs.BeneficiaryParty.Name,                      // beneficiary_party.client_name,
			&attrs.ChargesInformation.BearerCode,              // charges_info.bearer_code,
			&attrs.ChargesInformation.ReceiverChargesAmount,   // charges_info.receiver_charges.amount,
			&attrs.ChargesInformation.ReceiverChargesCurrency, // charges_info.receiver_charges.currency,
			pq.Array(&amounts),                                // charges_info.sender_charges,
			&attrs.Currency,                                   // currency,
			&attrs.DebtorParty.AccountName,                    // debtor_party.name,
			&attrs.DebtorParty.AccountNumber,                  // debtor_party.number,
			&attrs.DebtorParty.AccountNumberCode,              // debtor_party.number_code,
			&attrs.DebtorParty.AccountType,                    // debtor_party.type,
			&attrs.DebtorParty.Address,                        // debtor_party.address ,
			&attrs.DebtorParty.BankID,                         // debtor_party.bank_id,
			&attrs.DebtorParty.BankIDCode,                     // debtor_party.bank_id_code,
			&attrs.DebtorParty.Name,                           // debtor_party.client_name,
			&attrs.EndToEndReference,                          // e2e_reference,
			&attrs.Fx.ContractReference,                       // fx.contract_ref,
			&attrs.Fx.ExchangeRate,                            // fx.rate,
			&attrs.Fx.OriginalAmount,                          // fx.original_amount.amount,
			&attrs.Fx.OriginalCurrency,                        // fx.original_amount.currency,
			&attrs.NumericReference,                           // numeric_reference,
			&attrs.PaymentID,                                  // payment_id,
			&attrs.PaymentType,                                // payment_type,
			&attrs.ProcessingDate,                             // processing_date,
			&attrs.PaymentPurpose,                             // purpose,
			&attrs.Reference,                                  // reference,
			&attrs.PaymentScheme,                              // scheme,
			&attrs.SchemePaymentSubType,                       // scheme_payment_subtype,
			&attrs.SchemePaymentType,                          // scheme_payment_type,
			&attrs.SponsorParty.AccountNumber,                 // sponsor_party.account_number,
			&attrs.SponsorParty.BankID,                        // sponsor_party.bank_id,
			&attrs.SponsorParty.BankIDCode,                    // sponsor_party.bank_id_code
		)

		if err != nil {
			return nil, fmt.Errorf("db: error scanning row: %v", err)
		}

		attrs.ChargesInformation.SenderCharges = amountsToSenderCharges(amounts)
		payment.Attributes = &attrs
		payments = append(payments, &payment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db: error scanning rows: %v", err)
	}

	return payments, nil
}

// Update updates the details associated with the given paymentID. The current
// implementation is a basic one that doesn't support updating fields selectively.
//
// Update returns an error if the paymentID does not exist in the collection
func (dbpr *DBPaymentRepository) Update(paymentID strfmt.UUID, payment *models.Payment) (*models.Payment, error) {
	updateStmt := `
	UPDATE payments
	SET
		organisation = $2,
		version = $3,
		amount = $4,
		beneficiary_party.name = $5,
		beneficiary_party.number = $6,
		beneficiary_party.number_code = $7,
		beneficiary_party.type = $8,
		beneficiary_party.address = $9,
		beneficiary_party.bank_id = $10,
		beneficiary_party.bank_id_code = $11,
		beneficiary_party.client_name = $12,
		charges_info.bearer_code = $13,
		charges_info.receiver_charges.amount = $14,
		charges_info.receiver_charges.currency = $15,
		charges_info.sender_charges = $16,
		currency = $17,
		debtor_party.name = $18,
		debtor_party.number = $19,
		debtor_party.number_code = $20,
		debtor_party.type = $21,
		debtor_party.address = $22,
		debtor_party.bank_id = $23,
		debtor_party.bank_id_code = $24,
		debtor_party.client_name = $25,
		e2e_reference = $26,
		fx.contract_ref = $27,
		fx.rate = $28,
		fx.original_amount.amount = $29,
		fx.original_amount.currency = $30,
		numeric_reference = $31,
		payment_id = $32,
		payment_type = $33,
		processing_date = $34,
		purpose = $35,
		reference = $36,
		scheme = $37,
		scheme_payment_subtype = $38,
		scheme_payment_type = $39,
		sponsor_party.account_number = $40,
		sponsor_party.bank_id = $41,
		sponsor_party.bank_id_code = $42
	WHERE id = $1`

	version := *payment.Version + 1
	attrs := payment.Attributes
	amounts := senderChargesToAmounts(attrs.ChargesInformation.SenderCharges)
	res, err := dbpr.db.Exec(updateStmt,
		payment.ID,                                       // id,
		payment.OrganisationID,                           // organisation,
		version,                                          // version,
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
	if err != nil {
		return nil, fmt.Errorf("db: error executing update: %v", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("db: error getting rows affected by update: %v", err)
	}
	if count == 0 {
		return nil, newErrNoResults(fmt.Sprintf("db: payment with ID %s not found", paymentID))
	}

	// Ignore the original type attribute and fix it to TYPE_PAYMENT
	updated := copyPayment(payment)
	updated.Type = TYPE_PAYMENT
	updated.Version = &version
	return updated, nil
}

// copyPayment performs a deep copy of a models.Payment structure
func copyPayment(payment *models.Payment) *models.Payment {
	// Configuration for copystructure package to correctly copy strfmt.Date
	// Copy operation on this type fails if a custom copier function is not provided
	// because of the strfmt.RFC3339FullDate custom format.
	// This only needs to be done once
	if _, ok := copystructure.Copiers[reflect.TypeOf(strfmt.Date{})]; !ok {
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

	dup, _ := copystructure.Copy(*payment)
	paymentDup := dup.(models.Payment)
	return &paymentDup
}
