# pAPI (= payments API)

[![Build Status](https://travis-ci.org/volmedo/pAPI.svg?branch=master)][travis]
[![Go Report Card](https://goreportcard.com/badge/github.com/volmedo/pAPI)][go_report]

[travis]: https://travis-ci.org/volmedo/pAPI
[go_report]: https://goreportcard.com/report/github.com/volmedo/pAPI

pAPI is a payments API written in Go, a fictional service that offers standard CRUD functionality on `Payment` resources.

The main objective of this repo is serving as an example of what I consider modern software engineering practices applied to the development of backend architectures based in microservices. It is aimed at covering not only the development of the software itself, but also the tools and strategies needed to walk the full way from writing the first line of code to getting the service in production, following the principles around agile development and DevOps.

## Architecture overview

The payments API is implemented as a microservice that offers a REST API that allows clients to manage payment resources by offering standard CRUD functionality. API messages are written using the JSON format and are conformant with the [json:api specification](https://jsonapi.org/).

[Swagger/OpenAPI](https://swagger.io/) is used to specify the API contract with clients. Swagger allows writing API specifications using the YAML format. Apart from serving as documentation of the API's exported functionality and exptected inputs and outputs, swagger files can be used to automatically generate model and boilerplate code which takes care of common tasks such as request handling and input validation.

## Test-Driven Development/Behaviour-Driven Development

The Payments Service has been developed using TDD from start to finish. The development process starts by defining how the system as a whole should behave to satisfy business needs and requirements. These business requirements are usually expressed in the form of functionality the service is needed to offer to potential users and, thus, end to end or acceptance tests are the best kind of test to capture the essence of these requirements. Since these tests describe the system's capabilities from the viewpoint of an external user, they are usually part of the communication between technical and not technical stakeholders of the project. Because of this, [Cucumber](https://cucumber.io/) and the [Gherkin syntax](https://cucumber.io/docs/gherkin/) are great for acceptance tests, as they allow writing them in a form close to natural language, eliminating the need of having a technical background to read or write them.

This project's end to end tests are written in the form of `feature` files using Gherkin syntax. These files are then processed by [godog](https://github.com/DATA-DOG/godog), the semi-official implementation of Cucumber for Go. `feature` files, along with step implementations, can be found in the [e2e_test] folder.

The test strategy I used follows the well-known Test Pyramid approach, where e2e tests are at the top and unit tests form the base of the pyramid. Usually, e2e tests are more expensive in terms of time required to set up the environment and run them (as an example, e2e tests in this project are run against real infrastructure that gets deployed before the tests run and destroyed afterwards). Due to their higher cost, e2e tests are usually limited in number and only happy paths are checked as a way to guarantee that the system as a whole delivers the required functionality.

Deeper in the service logic, functionality is checked by unit tests. These tests are fast and are run several times during development. Writing the tests before any logic is implemented eases the process of defining what functionality is really needed, what architecture should be used and how the unit under test is expected to behave. As opposed to e2e tests, unit tests are developer tests, so the language used for development is also the most convenient to write unit tests in. Moreover, one of the great things about Go is its rich tooling, and testing is a first-class citizen in this tooling. Unit tests in this project use plain `go test` constructs. I don't even use an assertion library (with [testify](https://github.com/stretchr/testify) being one of the most prominent examples) in favor of standard mechanisms, such as `reflect.DeepEqual`. My opinion is that the little added value is not worth the time needed to learn and handle yet another API. Of course, this (as several other things) is debatable and I'm always open to being convinced of the opposite :).

## Continous Integration/Continuous Deployment

A critical requirement of Continuous Integration workflows is that every commit to the master/trunk branch must build and pass tests. CI/CD platforms and tools are key to enable high-throughput teams that aim at release software at a fast pace.

This project uses [Travis CI](https://travis-ci.org) for CI/CD. Travis configuration is done using a single [.travis.yml](.travis.yml) file that sets the environment up and then goes through each of the defined stages collecting their results. If any command returns an error, the build will be considered broken.

The configuration follows the usual lint-test-build cycle. A relevant point here is that the process is tailored to cover all tests, not only unit and integration ones but also end to end tests. These are not performed in a special, contrived environment. Instead, real infrastructure is created, the code just built is deployed in this infrastructure and end to end tests are done by using a test client to consume the API and check the results. Once the tests are done, the infrastructure is destroyed. Performing these tests in a test environment on real infrastructure, ideally as close as the production infrastructure the application will live in the future, builds a lot of confidence on the recently developed code.

In order for this approach to be practicable, infrastructure cannot be handled manually as this would be both time-consuming and error-prone. This is where Infrastructure as Code comes to the rescue.

## Infrastructure as Code

Aside from allowing managing infrastructure in a quick and efficient manner, on the main benefits one gets from using IaC is reproducible deployments that are driven by configuration files which can be commited to a repository along with the rest of the source code.

[Terraform](https://www.terraform.io/) is used in this project as IaC tool. As with other similar solutions, Terraform uses a proprietary DSL called HCL (Hashicorp Configuration Language) to declaratively describe infrastructure deployments. In this project I make use of [Terraform backends](https://www.terraform.io/docs/backends/) to store remote state, so that infrastructure can be managed from both the Travis CI environment and my local development environment.

The cloud provider I used in this project is [AWS](https://aws.amazon.com/), but all of the concepts and services can easily be translated to [GCP](https://cloud.google.com/), [Azure](https://azure.microsoft.com/) or any other IaaS provider.

## Amazon Web Services

Amazon offers a wide (wider, in fact) range of services and products as part of its cloud infrastructure offering. They provide different abstraction levels (PaaS, IaaS) and functionalities. One of the first and more popular services is [EC2](https://aws.amazon.com/ec2/), which allows creating virtual machines with shared or dedicated resources that can be used to host one or more services.

### Infrastructure planning

Since the pAPI is currently served by a single component, a single EC2 `t2.micro` instance will do for now (and it is covered by [AWS free tier](https://aws.amazon.com/free/) :)). When the service grows in complexity, I will explore the addition of new elements.

### AWS setup

Setting everything up for automatic infrastructure configuration could be as easy as using root credentials to perform every needed change, but this is an obvious call for disaster. Best practices mandate users, groups and policies to be correctly configured to grant the minimum set of permissions needed to manage the resources described in Terraform configuration.

Within IAM (AWS identity and role management service), I created a `travisci` user to be used by the CI/CD pipeline and added it to a `Terraformers` group. I then created the needed policies and attached them to the group. When defining permissions policies, bear in mind that everything in AWS is a resource and that even the simplest configuration usually involves several types of them.

As an example, the basic configuration described in this project consists of a single EC2 instance, but it comprises the following resources:

- the EC2 instance itself
- an AMI to launch the instance with
- an attached volume for storage
- security groups to define inbound and outbound traffic rules
- a key pair for SSH access

And one could also define a dedicated VPC, subnets and gateways instead of using the default ones.

The high level of granularity can sometimes make difficult to actually know what the minimum set of permissions looks like. A solution proposed by Amazon is [using CloudTrail event history](https://docs.aws.amazon.com/IAM/latest/UserGuide/best-practices.html#grant-least-privilege) to know exactly what APIs are called when operating the infrastructure and use that information to narrow privileges down. For more information, see [the discussion in this Terraform issue](https://github.com/hashicorp/terraform/issues/2834) and have a look at [this gist with example policies](https://gist.github.com/arsdehnel/70e292467ced2a39f472ddca44629c08).

I also decided to use the same AWS account to give support to Terraform's backend for remote state. Doing so implies creating an S3 bucket to store the state file itself and a DynamoDB table to enable state locking. [The instructions on the Terraform site](https://www.terraform.io/docs/backends/types/s3.html) are easy to follow. One gotcha I'd like to note here: the key in the DynamoDB table must be called `LockID` and the name is case-sensitive.
