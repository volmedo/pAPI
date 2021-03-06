# pAPI (payments API)

[![Build Status](https://travis-ci.org/volmedo/pAPI.svg?branch=master)][travis]
[![Go Report Card](https://goreportcard.com/badge/github.com/volmedo/pAPI)][go_report]

[travis]: https://travis-ci.org/volmedo/pAPI/
[go_report]: https://goreportcard.com/report/github.com/volmedo/pAPI/

<p align="center"><img src="logo.png" alt="pAPI logo" width="250"/></p>

pAPI is a payments API written in Go, a fictional cloud service that offers standard CRUD functionality on `payment` resources.

## Table of Contents

- [API design](#api-design)
  - [Operations](#operations)
    - [Common status codes](#common-status-codes)
    - [Create payment](#create-payment)
    - [Fetch payment](#fetch-payment)
    - [Update payment](#update-payment)
    - [Delete payment](#delete-payment)
    - [List payments](#list-payments)
  - [Rate limits](#rate-limits)
  - [Additional endpoints](#additional-endpoints)
- [Implementation details](#implementation-details)
  - [Automatic code generation](#automatic-code-generation)
  - [Test-Driven Development / Behaviour-Driven Development](#test-driven-development--behaviour-driven-development)
  - [Resource state persistence](#resource-state-persistence)
  - [Continuous Integration](#continuous-integration)
  - [Infrastructure as Code
    ](#infrastructure-as-code)
  - [Instrumentation and logging](#instrumentation-and-logging)
  - [Rate limiting](#rate-limiting)
  - [Configuration from the environment](#configuration-from-the-environment)
  - [Containerization](#containerization)
  - [Cluster deployment](#cluster-deployment)
- [Further work](#further-work)

## API design

Payments API handles `payment` resources. A payment represents a money transaction between a beneficiary party and a debtor party and contains all information needed to correctly process the transaction as well as maintaining adequate book-keeping. Payment objects are written using JSON and are conformant with the [json:api specification](https://jsonapi.org/).

An example payment could look like the following:

```json
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
}
```

For simplicity, it is assumed that full representations of `payment` resources will be used, i.e. `payment` objects in requests and responses will always contain every attribute except for `type` and `version`, which will be handled by the server.

### Operations

The API offers basic CRUD operations on a collection of `payment` resources. The API is designed following RESTful conventions around endpoint names and HTTP methods. The following table summarizes available actions:

| Action         |  Method  | Endpoint         | Description                                                    | Status codes            |
| -------------- | :------: | ---------------- | -------------------------------------------------------------- | ----------------------- |
| Create payment |  `POST`  | `/payments`      | Creates a new payment resource with the given details          | 201, 409, 422, 429, 500 |
| Fetch payment  |  `GET`   | `/payments/{id}` | Requests details about the payment resource identified by `id` | 200, 404, 422, 429, 500 |
| Update payment |  `PUT`   | `/payments/{id}` | Uses the provided data to update the payment with `id`         | 200, 404, 422, 429, 500 |
| Delete payment | `DELETE` | `/payments/{id}` | Deletes the payment resource identified by `id`                | 204, 404, 422, 429, 500 |
| List payments  |  `GET`   | `/payments`      | Fetches details about more than one payment as a collection    | 200, 400, 422, 429, 500 |

#### Common status codes

Status codes `422`, `429` and `500` are common to all endpoints:

- `422 Unprocessable Entity`: the client sent syntactically correct but semantically wrong data. Parameters with invalid values and missing fields in payment objects are the most common causes of this error.
- `429 Too Many Requests`: request rate limit reached.
- `500 Internal Server Error`: the server encountered an error while processing the request.

#### Create payment

Creates a new payment with the information given by the client in the request body. The new payment's `id` will be generated by the client and included in the payment object.

##### Request

|     Request      | Params |   Body    |
| :--------------: | :----: | :-------: |
| `POST /payments` |   -    | `payment` |

##### Response

| Status code    |   Body    | Description                                    |
| -------------- | :-------: | ---------------------------------------------- |
| `201 Created`  | `payment` | Resource created successfully                  |
| `409 Conflict` |     -     | There is already a payment with the given `id` |

#### Fetch payment

Asks the server for details about the payment with `id`.

##### Request

|       Request        | Params | Body |
| :------------------: | :----: | :--: |
| `GET /payments/{id}` |  `id`  |  -   |

##### Response

| Status code     |   Body    | Description                              |
| --------------- | :-------: | ---------------------------------------- |
| `200 OK`        | `payment` | Requested details retrieved successfully |
| `404 Not Found` |     -     | A payment with `id` could not be found   |

#### Update payment

Updates the information about the payment identified by `id` with the data contained in the request body. The `id` in the URI will be used to identify the payment. If the payment object sent in the request body contains an `id` field, it will be ignored.

`PUT` is used instead of `PATCH` to indicate that partial updates (i.e. updating only some attributes) are not allowed. Payment details will be updated by replacing payment representations as a whole.

##### Request

|       Request        | Params |   Body    |
| :------------------: | :----: | :-------: |
| `PUT /payments/{id}` |  `id`  | `payment` |

##### Response

| Status code     |   Body    | Description                            |
| --------------- | :-------: | -------------------------------------- |
| `200 OK`        | `payment` | Payment resource updated successfully  |
| `404 Not Found` |     -     | A payment with `id` could not be found |

#### Delete payment

Deletes the payment with the given `id`.

##### Request

|         Request         | Params | Body |
| :---------------------: | :----: | :--: |
| `DELETE /payments/{id}` |  `id`  |  -   |

##### Response

| Status code      | Body | Description                            |
| ---------------- | :--: | -------------------------------------- |
| `204 No Content` |  -   | Payment resource deleted successfully  |
| `404 Not Found`  |  -   | A payment with `id` could not be found |

#### List payments

Gets a view of payment resources as a collection.

##### Request

|     Request     |            Params            | Body |
| :-------------: | :--------------------------: | :--: |
| `GET /payments` | `page[number]`, `page[size]` |  -   |

This action supports pagination parameters:

- `page[number]`: The page number that is being requested. A page is a sub-collection of `page[size]` elements. The first page is number 0. This parameter defaults to 0 and, obviously, cannot be negative.
- `page[size]`: Number of elements per page. This parameter defaults to 10 and must be in the range (0, 100]. Requests with values outside this range will result in a `422 Unprocessable Entity` response.

##### Response

| Status code     |        Body        | Description                                                                                                          |
| --------------- | :----------------: | -------------------------------------------------------------------------------------------------------------------- |
| `200 OK`        | Array of `payment` | Requested details retrieved successfully                                                                             |
| `404 Not Found` |         -          | No payment matches the query. Either there are no payments or pagination parameters make the query return no results |

### Rate limits

The API implements request rate limit to avoid intentional or unintentional misuse of server resources. By default, a limit of 100 requests per second per client is imposed. If the client sends requests at higher rates, the server will return `429 Too Many Requests` to any request beyond the limit.

### Additional endpoints

Aside from the application endpoints, the service also offers additional endpoints that are useful from an operational point of view:

- `/health`: Returns `200` if the service is available and able to handle requests and `500` otherwise, following the guidelines for [Kubernetes liveness and readiness probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/).
- `/metrics`: Returns service metrics in Prometheus format.

## Implementation details

The payments API is implemented as a microservice that offers a REST API that allows clients to manage payment resources by offering standard CRUD functionality.

[Swagger/OpenAPI](https://swagger.io/) is used to specify the API contract with clients. Swagger allows writing API specifications using the YAML format. Apart from serving as documentation of the API's exported functionality and expected inputs and outputs, swagger files can be used to automatically generate model and boilerplate code which takes care of common tasks such as request handling and input validation.

Part of the code in the repo is automatically generated from the spec indeed. Specifically, code in [pkg/client](pkg/client), [pkg/models](pkg/models) and [pkg/restapi](pkg/restapi) is autogenerated.

The rest of the section highlights the key features implemented in the project. This is a quick rundown, please ask if you want more details ;).

### Automatic code generation

[go-swagger](https://goswagger.io/) is used to generate routing and validation code from the swagger spec. The tool was used with the great [Stratoscale templates](https://github.com/Stratoscale/swagger/) that carefully isolate implementation from generated code, producing interfaces that improve testability and play nicely with the standard `net/http` library.

### Test-Driven Development / Behaviour-Driven Development

Tests are used to describe desired behaviour and to drive design decisions. Required functionality is captured in [Cucumber](https://cucumber.io/) features, written using [Gherkin](https://cucumber.io/docs/gherkin/) syntax. Acceptance end to end tests are based on [DATA-DOG/godog](https://github.com/DATA-DOG/godog/), the semi-official implementation of Cucumber for Go.

End to end tests are then followed by integration and unit tests. Standard Go testing facilities are used for these tests, since they are developer tests, with the help of auxiliary libraries such as [mitchellh/copystructure](https://github.com/mitchellh/copystructure/) or [google/go-cmp](https://github.com/google/go-cmp/). [DATA-DOG/go-sqlmock](https://github.com/DATA-DOG/go-sqlmock/) serves as a database test double when unit testing the code that interfaces with the data backend.

### Resource state persistence

Data persistence is achieved by using [PostgreSQL](https://www.postgresql.org/) as a data backend. It is difficult to choose the right database for an application when one has little information about what the data looks like (e.g. what other resources could exist and/or how they relate with payment resources) or how it will be used. However, an RDBMS was chosen over NoSQL or document-oriented alternatives on the assumption that strong consistency (and ACID-compliancy in general) is a must.

The interface with the database is abstracted by the notion of Repository for payment resources and the implementation of a DB-based one. No ORMs were used to avoid the need to edit autogenerated model code to add annotations.

[lib/pq](https://github.com/lib/pq/) is used as driver and schema migrations are handled by means of [golang-migrate/migrate](https://github.com/golang-migrate/migrate/).

### Continuous Integration

An automated build pipeline is configured on [Travis CI](https://travis-ci.org/). The CI pipeline, which is triggered on every commit, lints the code (using [golangci/golangci-lint](https://github.com/golangci/golangci-lint/)), runs tests (configuring infrastructure when necessary) and publishes Docker images.

### Infrastructure as Code

End to end tests are run by deploying the service on real infrastructure in [Amazon Web Services](https://aws.amazon.com/). Thus, automating tests imply automating the generation and configuration of such infrastructure. [Terraform](https://www.terraform.io/) is the IaC tool of choice.

### Instrumentation and logging

Service metrics are exposed in [Prometheus](https://prometheus.io/) format thanks to an additional endpoint implemented using [slok/go-http-metrics](https://github.com/slok/go-http-metrics/). Collected metrics follow [the RED method](https://www.weave.works/blog/the-red-method-key-metrics-for-microservices-architecture/).

Logging in the application is done in a lightweight manner, avoiding huge amounts of logs that make finding the information needed to solve an issue quite hard. Instead, metrics are favoured as the main source of information about the service's status. Request logging is avoided, and only server errors are logged. Go's standard log package is enough for this use.

### Rate limiting

Imposing rate limits is essential to avoid server resource misuse. Rate limiting is implemented by adding [ulule/limiter](https://github.com/ulule/limiter/) middleware to the handler chain.

### Configuration from the environment

Service configuration can be stored in the environment, following guidelines and conventions such as those proposed by [The Twelve-Factor App](https://12factor.net/). [namsral/flag](https://github.com/namsral/flag/) is a drop-in replacement for Go stdlib's `flag` package that is able to read configuration parameters from environment variables as well as regular command-line arguments.

### Containerization

To ease deployment, [Docker](https://www.docker.com/) container images are generated for the service and uploaded to a repository on [Docker Hub](https://cloud.docker.com/repository/docker/volmedo/papi/).

### Cluster deployment

Once the service is containerized, it can be easily deployed in a cluster using [Kubernetes](https://kubernetes.io/). To check K8s configuration, the service is deployed in a local cluster created with [kind](https://kind.sigs.k8s.io/) and the end to end tests are run against this local deployment.

## Further work

Some features have been left out for simplicity, but would be mandatory in a real-world scenario.

- **HTTPS:** [Core Go crypto libraries](https://godoc.org/golang.org/x/crypto/acme/autocert/) allow automatic generation of SSL certificates from [Let's Encrypt](https://letsencrypt.org/). HTTPS features could, however, be handled at the infrastructure level, as SSL termination is a common option in most load balancers and API gateways.
- **Authorization:** A library such as [go-oauth2/oauth2](https://github.com/go-oauth2/oauth2/) could be used to implement server-side OAuth2 to handle client authorization to use of API resources. As was the case with SSL termination, authorization is likely to be handled centrally in most microservices architectures.
- **Continuous Deployment:** Currently, Travis is only used for CI, but build artifacts are not deployed automatically on success. Terraform configurations could be rewritten into modules that could then be parametrized to allow for different deployments depending on the environment (staging, production, etc.)
- **A lot of other things that doesn't make sense in a restricted-scope scenario like this project's one:** Advanced tools and services useful in complex architectures, such as service meshes, API gateways, log ingestion platforms, credential managers...

<sub>\*Logo font: Aristotelica Family by Zetafonts - http://www.zetafonts.com/collection/1077</sub>
