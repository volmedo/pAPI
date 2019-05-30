CREATE TYPE account_number_code AS ENUM (
    'IBAN',
    'BBAN'
);

CREATE TYPE amount AS (
    amount      NUMERIC(10, 2),
    currency    VARCHAR(3)
);

CREATE TYPE bank_id_code AS ENUM (
    'GBDSC',
    'SWBIC'
);

CREATE TYPE bearer AS ENUM (
    'DEBT',
    'CRED',
    'SHAR',
    'SLEV'
);

CREATE TYPE charges_info AS (
    bearer_code         bearer,
    receiver_charges    amount,
    sender_charges      amount []
);

CREATE TYPE customer_account AS (
    name            TEXT,
    number          TEXT,
    number_code     account_number_code,
    type            INT,
    address         TEXT,
    bank_id         TEXT,
    bank_id_code    bank_id_code,
    client_name     TEXT
);

CREATE TYPE fx AS (
    contract_ref    TEXT,
    rate            NUMERIC(10, 5),
    original_amount amount
);

CREATE TYPE scheme AS ENUM (
    'BACS', 
    'CHAPS', 
    'FPS', 
    'SEPA-CT', 
    'SEPAINSTANT', 
    'SWIFT'
);

CREATE TYPE scheme_payment_type AS ENUM (
    'ImmediatePayment',
    'ForwardDatedPayment',
    'StandingOrder',
    'Credit',
    'Interest',
    'Dividend'
);

CREATE TYPE scheme_payment_subtype AS ENUM (
    'TelephoneBanking',
    'InternetBanking',
    'BranchInstruction',
    'Letter',
    'Email',
    'MobilePaymentsService'
);

CREATE TYPE sponsor_account AS (
    account_number  TEXT,
    bank_id         TEXT,
    bank_id_code    bank_id_code
);

CREATE TABLE IF NOT EXISTS payments (
    id                      UUID PRIMARY KEY,
    organisation            UUID,
    version                 INT,
    amount                  NUMERIC(8, 2),
    beneficiary_party       customer_account,
    charges_info            charges_info,
    currency                VARCHAR(3),
    debtor_party            customer_account,
    e2e_reference           TEXT,
    fx                      fx,
    numeric_reference       TEXT,
    payment_id              TEXT,
    payment_type            scheme_payment_type,
    processing_date         date,
    purpose                 TEXT,
    reference               TEXT,
    scheme                  scheme,
    scheme_payment_subtype  scheme_payment_subtype,
    scheme_payment_type     scheme_payment_type,
    sponsor_party           sponsor_account
);
