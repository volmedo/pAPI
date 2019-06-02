// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/vnd.api+json"
  ],
  "produces": [
    "application/vnd.api+json"
  ],
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Payments API as specified in Form3 take home test",
    "title": "Payments API",
    "version": "1"
  },
  "host": "api.example.com",
  "basePath": "/v1",
  "paths": {
    "/payments": {
      "get": {
        "tags": [
          "Payments"
        ],
        "summary": "List payments",
        "operationId": "listPayments",
        "parameters": [
          {
            "type": "integer",
            "default": 0,
            "description": "Which page to select",
            "name": "page[number]",
            "in": "query"
          },
          {
            "maximum": 100,
            "minimum": 1,
            "type": "integer",
            "default": 10,
            "description": "Number of items per page",
            "name": "page[size]",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "List of payment details",
            "schema": {
              "$ref": "#/definitions/PaymentDetailsListResponse"
            }
          },
          "400": {
            "description": "Bad Request",
            "schema": {
              "$ref": "#/definitions/ApiError"
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "$ref": "#/definitions/ApiError"
            }
          }
        }
      },
      "post": {
        "tags": [
          "Payments"
        ],
        "summary": "Create payment",
        "operationId": "createPayment",
        "parameters": [
          {
            "name": "Payment creation request",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/PaymentCreationRequest"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Payment created successfully",
            "schema": {
              "$ref": "#/definitions/PaymentCreationResponse"
            }
          },
          "409": {
            "description": "A payment with the given ID already exists",
            "schema": {
              "$ref": "#/definitions/ApiError"
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "$ref": "#/definitions/ApiError"
            }
          }
        }
      }
    },
    "/payments/{id}": {
      "get": {
        "tags": [
          "Payments"
        ],
        "summary": "Fetch payment",
        "operationId": "getPayment",
        "parameters": [
          {
            "type": "string",
            "format": "uuid",
            "description": "ID of payment to fetch",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Payment details",
            "schema": {
              "$ref": "#/definitions/PaymentDetailsResponse"
            }
          },
          "404": {
            "description": "Payment Not Found",
            "schema": {
              "$ref": "#/definitions/ApiError"
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "$ref": "#/definitions/ApiError"
            }
          }
        }
      },
      "put": {
        "tags": [
          "Payments"
        ],
        "summary": "Update payment details",
        "operationId": "updatePayment",
        "parameters": [
          {
            "type": "string",
            "format": "uuid",
            "description": "ID of payment to update",
            "name": "id",
            "in": "path",
            "required": true
          },
          {
            "description": "New payment details",
            "name": "Payment update request",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/PaymentUpdateRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Payment details",
            "schema": {
              "$ref": "#/definitions/PaymentUpdateResponse"
            }
          },
          "404": {
            "description": "Payment Not Found",
            "schema": {
              "$ref": "#/definitions/ApiError"
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "$ref": "#/definitions/ApiError"
            }
          }
        }
      },
      "delete": {
        "tags": [
          "Payments"
        ],
        "summary": "Deletes a payment resource",
        "operationId": "deletePayment",
        "parameters": [
          {
            "type": "string",
            "format": "uuid",
            "description": "ID of payment to delete",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "204": {
            "description": "Payment deleted OK. No body content will be returned"
          },
          "404": {
            "description": "Payment Not Found",
            "schema": {
              "$ref": "#/definitions/ApiError"
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "$ref": "#/definitions/ApiError"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "AccountNumber": {
      "description": "Account number",
      "type": "string",
      "example": "71268996"
    },
    "Amount": {
      "description": "Amount of money. Requires 1 to 2 decimal places.",
      "type": "string",
      "pattern": "^[0-9.]{0,20}$",
      "example": "10.00"
    },
    "ApiError": {
      "type": "object",
      "properties": {
        "error_code": {
          "type": "string",
          "format": "uuid"
        },
        "error_message": {
          "type": "string"
        }
      }
    },
    "BankId": {
      "description": "Financial institution identification",
      "type": "string",
      "example": "333333"
    },
    "BankIdCode": {
      "description": "The type of identification provided at ` + "`" + `bank_id` + "`" + ` attribute. Must be ISO code as listed in the [External Code Sets spreadsheet](https://www.iso20022.org/external_code_list.page)",
      "type": "string",
      "enum": [
        "SWBIC",
        "GBDSC",
        "BE",
        "FR",
        "DEBLZ",
        "GRBIC",
        "ITNCC",
        "PLKNR",
        "PTNCC",
        "ESNCC",
        "CHBCC"
      ],
      "example": "GBDSC"
    },
    "ChargesInformation": {
      "type": "object",
      "properties": {
        "bearer_code": {
          "description": "Specifies which party/parties will bear the charges associated with the processing of the payment transaction.",
          "type": "string",
          "enum": [
            "DEBT",
            "CRED",
            "SHAR",
            "SLEV"
          ],
          "example": "SLEV"
        },
        "receiver_charges_amount": {
          "description": "Transaction charges due to the receiver of the transaction.",
          "$ref": "#/definitions/Amount"
        },
        "receiver_charges_currency": {
          "$ref": "#/definitions/Currency"
        },
        "sender_charges": {
          "type": "array",
          "items": {
            "description": "List of transaction charges due to the sender of the transaction",
            "type": "object",
            "properties": {
              "amount": {
                "description": "Amount of each transaction charge due to the sender of the transaction.",
                "$ref": "#/definitions/Amount"
              },
              "currency": {
                "$ref": "#/definitions/Currency"
              }
            }
          }
        }
      }
    },
    "Currency": {
      "description": "Currency code as defined in [ISO 4217](http://www.iso.org/iso/home/standards/currency_codes.htm).",
      "type": "string",
      "pattern": "^[A-Z]{3}$",
      "example": "EUR"
    },
    "Links": {
      "type": "object",
      "properties": {
        "first": {
          "description": "Link to the first resource in the list",
          "type": "string",
          "example": "https://api.test.example.com/v1/api_name/resource_type"
        },
        "last": {
          "description": "Link to the last resource in the list",
          "type": "string",
          "example": "https://api.test.example.com/v1/api_name/resource_type"
        },
        "next": {
          "description": "Link to the next resource in the list",
          "type": "string",
          "example": "https://api.test.example.com/v1/api_name/resource_type"
        },
        "prev": {
          "description": "Link to the previous resource in the list",
          "type": "string",
          "example": "https://api.test.example.com/v1/api_name/resource_type"
        },
        "self": {
          "description": "Link to this resource type",
          "type": "string",
          "example": "https://api.test.example.com/v1/api_name/resource_type"
        }
      }
    },
    "Payment": {
      "type": "object",
      "required": [
        "id",
        "organisation_id",
        "attributes"
      ],
      "properties": {
        "attributes": {
          "type": "object",
          "properties": {
            "amount": {
              "description": "Amount of money moved between the instructing agent and instructed agent",
              "$ref": "#/definitions/Amount"
            },
            "beneficiary_party": {
              "$ref": "#/definitions/PaymentParty"
            },
            "charges_information": {
              "$ref": "#/definitions/ChargesInformation"
            },
            "currency": {
              "$ref": "#/definitions/Currency"
            },
            "debtor_party": {
              "$ref": "#/definitions/PaymentParty"
            },
            "end_to_end_reference": {
              "description": "Unique identification, as assigned by the initiating party, to unambiguously identify the transaction. This identification is passed on, unchanged, throughout the entire end-to-end chain.",
              "type": "string",
              "example": "PAYMENT REF: 20094"
            },
            "fx": {
              "type": "object",
              "properties": {
                "contract_reference": {
                  "description": "Reference to the foreign exchange contract associated with the transaction",
                  "type": "string",
                  "example": "FXCONTRACT/REF/123567"
                },
                "exchange_rate": {
                  "description": "Factor used to convert an amount from the instructed currency into the transaction currency. Decimal value, represented as a string, maximum length 12. Must be \u003e 0.",
                  "type": "string",
                  "example": "0.13343"
                },
                "original_amount": {
                  "description": "Amount of money to be moved between the debtor and creditor, before deduction of charges, expressed in the currency as instructed by the initiating party. Decimal value. Must be \u003e 0.",
                  "$ref": "#/definitions/Amount"
                },
                "original_currency": {
                  "description": "Currency of ` + "`" + `orginal_amount` + "`" + `.",
                  "$ref": "#/definitions/Currency"
                }
              }
            },
            "numeric_reference": {
              "description": "Numeric reference field, see scheme specific descriptions for usage",
              "type": "string",
              "example": "0001"
            },
            "payment_id": {
              "description": "Payment identification (legacy?)",
              "type": "string",
              "example": "123456789012345678"
            },
            "payment_purpose": {
              "description": "Purpose of the payment in a proprietary form",
              "type": "string",
              "example": "Paying for goods/services"
            },
            "payment_scheme": {
              "description": "Clearing infrastructure through which the payment instruction is to be processed. Default for given organisation ID is used if left empty. Currently only FPS is supported.",
              "type": "string",
              "enum": [
                "FPS"
              ],
              "example": "FPS"
            },
            "payment_type": {
              "type": "string",
              "enum": [
                "Credit"
              ]
            },
            "processing_date": {
              "description": "Date on which the payment is to be debited from the debtor account. Formatted according to ISO 8601 format YYYY-MM-DD.",
              "type": "string",
              "format": "date",
              "example": "2015-02-12"
            },
            "reference": {
              "description": "Payment reference for beneficiary use",
              "type": "string",
              "example": "rent for oct"
            },
            "scheme_payment_sub_type": {
              "description": "The scheme specific payment sub type",
              "type": "string",
              "enum": [
                "TelephoneBanking",
                "InternetBanking",
                "BranchInstruction",
                "Letter",
                "Email",
                "MobilePaymentsService"
              ],
              "example": "TelephoneBanking"
            },
            "scheme_payment_type": {
              "description": "The scheme-specific payment type",
              "type": "string",
              "enum": [
                "ImmediatePayment",
                "ForwardDatedPayment",
                "StandingOrder"
              ],
              "example": "ImmediatePayment"
            },
            "sponsor_party": {
              "description": "Sponsor party",
              "type": "object",
              "properties": {
                "account_number": {
                  "$ref": "#/definitions/AccountNumber"
                },
                "bank_id": {
                  "$ref": "#/definitions/BankId"
                },
                "bank_id_code": {
                  "$ref": "#/definitions/BankIdCode"
                }
              }
            }
          }
        },
        "id": {
          "description": "Unique resource ID",
          "type": "string",
          "format": "uuid",
          "example": "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43"
        },
        "organisation_id": {
          "description": "Unique ID of the organisation this resource is created by",
          "type": "string",
          "format": "uuid",
          "example": "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb"
        },
        "type": {
          "description": "Name of the resource type",
          "type": "string",
          "pattern": "^[A-Za-z_]*$",
          "example": "Payment"
        },
        "version": {
          "description": "Version number",
          "type": "integer",
          "example": 0
        }
      }
    },
    "PaymentCreationRequest": {
      "type": "object",
      "required": [
        "data"
      ],
      "properties": {
        "data": {
          "$ref": "#/definitions/Payment"
        }
      }
    },
    "PaymentCreationResponse": {
      "type": "object",
      "required": [
        "data"
      ],
      "properties": {
        "data": {
          "$ref": "#/definitions/Payment"
        },
        "links": {
          "$ref": "#/definitions/Links"
        }
      }
    },
    "PaymentDetailsListResponse": {
      "type": "object",
      "properties": {
        "data": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Payment"
          }
        },
        "links": {
          "$ref": "#/definitions/Links"
        }
      }
    },
    "PaymentDetailsResponse": {
      "type": "object",
      "properties": {
        "data": {
          "$ref": "#/definitions/Payment"
        },
        "links": {
          "$ref": "#/definitions/Links"
        }
      }
    },
    "PaymentParty": {
      "type": "object",
      "properties": {
        "account_name": {
          "description": "Name of beneficiary/debtor as given with account",
          "type": "string",
          "example": "James Bond"
        },
        "account_number": {
          "$ref": "#/definitions/AccountNumber"
        },
        "account_number_code": {
          "description": "The type of identification given at ` + "`" + `account_number` + "`" + ` attribute",
          "type": "string",
          "enum": [
            "IBAN",
            "BBAN"
          ],
          "example": "IBAN"
        },
        "account_type": {
          "description": "The type of the account given with account_number. Single digit number. Only required if requested by the beneficiary party. Defaults to 0.",
          "type": "integer",
          "example": 0
        },
        "address": {
          "description": "Beneficiary/debtor address",
          "type": "string",
          "example": "1 Clarence Mew, Horsforth, Leeds Ls18 4EP"
        },
        "bank_id": {
          "$ref": "#/definitions/BankId"
        },
        "bank_id_code": {
          "$ref": "#/definitions/BankIdCode"
        },
        "name": {
          "description": "Beneficiary/debtor name",
          "type": "string",
          "example": "Norman Smith"
        }
      }
    },
    "PaymentUpdateRequest": {
      "type": "object",
      "required": [
        "data"
      ],
      "properties": {
        "data": {
          "$ref": "#/definitions/Payment"
        }
      }
    },
    "PaymentUpdateResponse": {
      "type": "object",
      "required": [
        "data"
      ],
      "properties": {
        "data": {
          "$ref": "#/definitions/Payment"
        },
        "links": {
          "$ref": "#/definitions/Links"
        }
      }
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/vnd.api+json"
  ],
  "produces": [
    "application/vnd.api+json"
  ],
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Payments API as specified in Form3 take home test",
    "title": "Payments API",
    "version": "1"
  },
  "host": "api.example.com",
  "basePath": "/v1",
  "paths": {
    "/payments": {
      "get": {
        "tags": [
          "Payments"
        ],
        "summary": "List payments",
        "operationId": "listPayments",
        "parameters": [
          {
            "minimum": 0,
            "type": "integer",
            "default": 0,
            "description": "Which page to select",
            "name": "page[number]",
            "in": "query"
          },
          {
            "maximum": 100,
            "minimum": 1,
            "type": "integer",
            "default": 10,
            "description": "Number of items per page",
            "name": "page[size]",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "List of payment details",
            "schema": {
              "$ref": "#/definitions/PaymentDetailsListResponse"
            }
          },
          "400": {
            "description": "Bad Request",
            "schema": {
              "$ref": "#/definitions/ApiError"
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "$ref": "#/definitions/ApiError"
            }
          }
        }
      },
      "post": {
        "tags": [
          "Payments"
        ],
        "summary": "Create payment",
        "operationId": "createPayment",
        "parameters": [
          {
            "name": "Payment creation request",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/PaymentCreationRequest"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Payment created successfully",
            "schema": {
              "$ref": "#/definitions/PaymentCreationResponse"
            }
          },
          "409": {
            "description": "A payment with the given ID already exists",
            "schema": {
              "$ref": "#/definitions/ApiError"
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "$ref": "#/definitions/ApiError"
            }
          }
        }
      }
    },
    "/payments/{id}": {
      "get": {
        "tags": [
          "Payments"
        ],
        "summary": "Fetch payment",
        "operationId": "getPayment",
        "parameters": [
          {
            "type": "string",
            "format": "uuid",
            "description": "ID of payment to fetch",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Payment details",
            "schema": {
              "$ref": "#/definitions/PaymentDetailsResponse"
            }
          },
          "404": {
            "description": "Payment Not Found",
            "schema": {
              "$ref": "#/definitions/ApiError"
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "$ref": "#/definitions/ApiError"
            }
          }
        }
      },
      "put": {
        "tags": [
          "Payments"
        ],
        "summary": "Update payment details",
        "operationId": "updatePayment",
        "parameters": [
          {
            "type": "string",
            "format": "uuid",
            "description": "ID of payment to update",
            "name": "id",
            "in": "path",
            "required": true
          },
          {
            "description": "New payment details",
            "name": "Payment update request",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/PaymentUpdateRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Payment details",
            "schema": {
              "$ref": "#/definitions/PaymentUpdateResponse"
            }
          },
          "404": {
            "description": "Payment Not Found",
            "schema": {
              "$ref": "#/definitions/ApiError"
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "$ref": "#/definitions/ApiError"
            }
          }
        }
      },
      "delete": {
        "tags": [
          "Payments"
        ],
        "summary": "Deletes a payment resource",
        "operationId": "deletePayment",
        "parameters": [
          {
            "type": "string",
            "format": "uuid",
            "description": "ID of payment to delete",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "204": {
            "description": "Payment deleted OK. No body content will be returned"
          },
          "404": {
            "description": "Payment Not Found",
            "schema": {
              "$ref": "#/definitions/ApiError"
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "$ref": "#/definitions/ApiError"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "AccountNumber": {
      "description": "Account number",
      "type": "string",
      "example": "71268996"
    },
    "Amount": {
      "description": "Amount of money. Requires 1 to 2 decimal places.",
      "type": "string",
      "pattern": "^[0-9.]{0,20}$",
      "example": "10.00"
    },
    "ApiError": {
      "type": "object",
      "properties": {
        "error_code": {
          "type": "string",
          "format": "uuid"
        },
        "error_message": {
          "type": "string"
        }
      }
    },
    "BankId": {
      "description": "Financial institution identification",
      "type": "string",
      "example": "333333"
    },
    "BankIdCode": {
      "description": "The type of identification provided at ` + "`" + `bank_id` + "`" + ` attribute. Must be ISO code as listed in the [External Code Sets spreadsheet](https://www.iso20022.org/external_code_list.page)",
      "type": "string",
      "enum": [
        "SWBIC",
        "GBDSC",
        "BE",
        "FR",
        "DEBLZ",
        "GRBIC",
        "ITNCC",
        "PLKNR",
        "PTNCC",
        "ESNCC",
        "CHBCC"
      ],
      "example": "GBDSC"
    },
    "ChargesInformation": {
      "type": "object",
      "properties": {
        "bearer_code": {
          "description": "Specifies which party/parties will bear the charges associated with the processing of the payment transaction.",
          "type": "string",
          "enum": [
            "DEBT",
            "CRED",
            "SHAR",
            "SLEV"
          ],
          "example": "SLEV"
        },
        "receiver_charges_amount": {
          "description": "Transaction charges due to the receiver of the transaction.",
          "$ref": "#/definitions/Amount"
        },
        "receiver_charges_currency": {
          "$ref": "#/definitions/Currency"
        },
        "sender_charges": {
          "type": "array",
          "items": {
            "description": "List of transaction charges due to the sender of the transaction",
            "type": "object",
            "properties": {
              "amount": {
                "description": "Amount of each transaction charge due to the sender of the transaction.",
                "$ref": "#/definitions/Amount"
              },
              "currency": {
                "$ref": "#/definitions/Currency"
              }
            }
          }
        }
      }
    },
    "Currency": {
      "description": "Currency code as defined in [ISO 4217](http://www.iso.org/iso/home/standards/currency_codes.htm).",
      "type": "string",
      "pattern": "^[A-Z]{3}$",
      "example": "EUR"
    },
    "Links": {
      "type": "object",
      "properties": {
        "first": {
          "description": "Link to the first resource in the list",
          "type": "string",
          "example": "https://api.test.example.com/v1/api_name/resource_type"
        },
        "last": {
          "description": "Link to the last resource in the list",
          "type": "string",
          "example": "https://api.test.example.com/v1/api_name/resource_type"
        },
        "next": {
          "description": "Link to the next resource in the list",
          "type": "string",
          "example": "https://api.test.example.com/v1/api_name/resource_type"
        },
        "prev": {
          "description": "Link to the previous resource in the list",
          "type": "string",
          "example": "https://api.test.example.com/v1/api_name/resource_type"
        },
        "self": {
          "description": "Link to this resource type",
          "type": "string",
          "example": "https://api.test.example.com/v1/api_name/resource_type"
        }
      }
    },
    "Payment": {
      "type": "object",
      "required": [
        "id",
        "organisation_id",
        "attributes"
      ],
      "properties": {
        "attributes": {
          "type": "object",
          "properties": {
            "amount": {
              "description": "Amount of money moved between the instructing agent and instructed agent",
              "$ref": "#/definitions/Amount"
            },
            "beneficiary_party": {
              "$ref": "#/definitions/PaymentParty"
            },
            "charges_information": {
              "$ref": "#/definitions/ChargesInformation"
            },
            "currency": {
              "$ref": "#/definitions/Currency"
            },
            "debtor_party": {
              "$ref": "#/definitions/PaymentParty"
            },
            "end_to_end_reference": {
              "description": "Unique identification, as assigned by the initiating party, to unambiguously identify the transaction. This identification is passed on, unchanged, throughout the entire end-to-end chain.",
              "type": "string",
              "example": "PAYMENT REF: 20094"
            },
            "fx": {
              "type": "object",
              "properties": {
                "contract_reference": {
                  "description": "Reference to the foreign exchange contract associated with the transaction",
                  "type": "string",
                  "example": "FXCONTRACT/REF/123567"
                },
                "exchange_rate": {
                  "description": "Factor used to convert an amount from the instructed currency into the transaction currency. Decimal value, represented as a string, maximum length 12. Must be \u003e 0.",
                  "type": "string",
                  "example": "0.13343"
                },
                "original_amount": {
                  "description": "Amount of money to be moved between the debtor and creditor, before deduction of charges, expressed in the currency as instructed by the initiating party. Decimal value. Must be \u003e 0.",
                  "$ref": "#/definitions/Amount"
                },
                "original_currency": {
                  "description": "Currency of ` + "`" + `orginal_amount` + "`" + `.",
                  "$ref": "#/definitions/Currency"
                }
              }
            },
            "numeric_reference": {
              "description": "Numeric reference field, see scheme specific descriptions for usage",
              "type": "string",
              "example": "0001"
            },
            "payment_id": {
              "description": "Payment identification (legacy?)",
              "type": "string",
              "example": "123456789012345678"
            },
            "payment_purpose": {
              "description": "Purpose of the payment in a proprietary form",
              "type": "string",
              "example": "Paying for goods/services"
            },
            "payment_scheme": {
              "description": "Clearing infrastructure through which the payment instruction is to be processed. Default for given organisation ID is used if left empty. Currently only FPS is supported.",
              "type": "string",
              "enum": [
                "FPS"
              ],
              "example": "FPS"
            },
            "payment_type": {
              "type": "string",
              "enum": [
                "Credit"
              ]
            },
            "processing_date": {
              "description": "Date on which the payment is to be debited from the debtor account. Formatted according to ISO 8601 format YYYY-MM-DD.",
              "type": "string",
              "format": "date",
              "example": "2015-02-12"
            },
            "reference": {
              "description": "Payment reference for beneficiary use",
              "type": "string",
              "example": "rent for oct"
            },
            "scheme_payment_sub_type": {
              "description": "The scheme specific payment sub type",
              "type": "string",
              "enum": [
                "TelephoneBanking",
                "InternetBanking",
                "BranchInstruction",
                "Letter",
                "Email",
                "MobilePaymentsService"
              ],
              "example": "TelephoneBanking"
            },
            "scheme_payment_type": {
              "description": "The scheme-specific payment type",
              "type": "string",
              "enum": [
                "ImmediatePayment",
                "ForwardDatedPayment",
                "StandingOrder"
              ],
              "example": "ImmediatePayment"
            },
            "sponsor_party": {
              "description": "Sponsor party",
              "type": "object",
              "properties": {
                "account_number": {
                  "$ref": "#/definitions/AccountNumber"
                },
                "bank_id": {
                  "$ref": "#/definitions/BankId"
                },
                "bank_id_code": {
                  "$ref": "#/definitions/BankIdCode"
                }
              }
            }
          }
        },
        "id": {
          "description": "Unique resource ID",
          "type": "string",
          "format": "uuid",
          "example": "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43"
        },
        "organisation_id": {
          "description": "Unique ID of the organisation this resource is created by",
          "type": "string",
          "format": "uuid",
          "example": "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb"
        },
        "type": {
          "description": "Name of the resource type",
          "type": "string",
          "pattern": "^[A-Za-z_]*$",
          "example": "Payment"
        },
        "version": {
          "description": "Version number",
          "type": "integer",
          "minimum": 0,
          "example": 0
        }
      }
    },
    "PaymentCreationRequest": {
      "type": "object",
      "required": [
        "data"
      ],
      "properties": {
        "data": {
          "$ref": "#/definitions/Payment"
        }
      }
    },
    "PaymentCreationResponse": {
      "type": "object",
      "required": [
        "data"
      ],
      "properties": {
        "data": {
          "$ref": "#/definitions/Payment"
        },
        "links": {
          "$ref": "#/definitions/Links"
        }
      }
    },
    "PaymentDetailsListResponse": {
      "type": "object",
      "properties": {
        "data": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Payment"
          }
        },
        "links": {
          "$ref": "#/definitions/Links"
        }
      }
    },
    "PaymentDetailsResponse": {
      "type": "object",
      "properties": {
        "data": {
          "$ref": "#/definitions/Payment"
        },
        "links": {
          "$ref": "#/definitions/Links"
        }
      }
    },
    "PaymentParty": {
      "type": "object",
      "properties": {
        "account_name": {
          "description": "Name of beneficiary/debtor as given with account",
          "type": "string",
          "example": "James Bond"
        },
        "account_number": {
          "$ref": "#/definitions/AccountNumber"
        },
        "account_number_code": {
          "description": "The type of identification given at ` + "`" + `account_number` + "`" + ` attribute",
          "type": "string",
          "enum": [
            "IBAN",
            "BBAN"
          ],
          "example": "IBAN"
        },
        "account_type": {
          "description": "The type of the account given with account_number. Single digit number. Only required if requested by the beneficiary party. Defaults to 0.",
          "type": "integer",
          "example": 0
        },
        "address": {
          "description": "Beneficiary/debtor address",
          "type": "string",
          "example": "1 Clarence Mew, Horsforth, Leeds Ls18 4EP"
        },
        "bank_id": {
          "$ref": "#/definitions/BankId"
        },
        "bank_id_code": {
          "$ref": "#/definitions/BankIdCode"
        },
        "name": {
          "description": "Beneficiary/debtor name",
          "type": "string",
          "example": "Norman Smith"
        }
      }
    },
    "PaymentUpdateRequest": {
      "type": "object",
      "required": [
        "data"
      ],
      "properties": {
        "data": {
          "$ref": "#/definitions/Payment"
        }
      }
    },
    "PaymentUpdateResponse": {
      "type": "object",
      "required": [
        "data"
      ],
      "properties": {
        "data": {
          "$ref": "#/definitions/Payment"
        },
        "links": {
          "$ref": "#/definitions/Links"
        }
      }
    }
  }
}`))
}
