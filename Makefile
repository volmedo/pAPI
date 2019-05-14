.DEFAULT_GOAL := build
SHELL := /bin/bash

SPEC ?= pAPI-swagger.yaml
PKG ?= pkg
CMD ?= cmd
E2E ?= e2e_test

GO ?= go
TARGET_OS ?= linux
TARGET_ARCH ?= amd64

SWAGGER_VER ?= v0.19.0
SWAGGER = docker run --rm -e GOPATH=/go -v "$(PWD)":"$(PWD)" -w "$(PWD)" quay.io/goswagger/swagger:$(SWAGGER_VER)

GOLANGCI_LINT_VER ?= v1.16.0
GOLANGCI_LINT = docker run --rm -v "$(PWD)":"$(PWD)" -w "$(PWD)" golangci/golangci-lint:$(GOLANGCI_LINT_VER)

TF_SSH_KEY_PATH ?= "$(PWD)/tf_ssh_key"
TF_DIR ?= terraform
TF_VER ?= 0.11.13
TERRAFORM = docker run --rm -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY -v "$(PWD)":"$(PWD)" -w "$(PWD)/$(TF_DIR)" hashicorp/terraform:$(TF_VER)

SRV_BIN_NAME ?= papisrv

swagger.validate:
	$(SWAGGER) validate $(SPEC)
	
swagger.clean:
	rm -r $(PKG)/models
	rm -r $(PKG)/restapi

swagger.generate.server: swagger.clean
	$(SWAGGER) generate server --spec=$(SPEC) --template=stratoscale --target=$(PKG)

swagger.generate.client:
	$(SWAGGER) generate client --spec=$(SPEC) --template=stratoscale --target=$(PKG) --skip-models

lint:
	$(GOLANGCI_LINT) golangci-lint run --no-config --skip-dirs "$(PKG)/(client|models|restapi)" --disable unused

test.unit:
	$(GO) test -v -race ./$(PKG)/impl

test.component:
	$(GO) test -v -race ./$(CMD)/server

build: ./$(CMD)/server/main.go
	GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) $(GO) build -o $(SRV_BIN_NAME) ./$(CMD)/server/main.go

test.e2e.local:
	$(GO) build -o testsrv ./$(CMD)/server/main.go ;\
	./testsrv -port=8080 & \
	SERVER_PID=$$! ;\
	$(GO) test -v -race ./$(E2E) -host=localhost -port=8080 ;\
	kill $$SERVER_PID ;\
	rm testsrv

test.e2e:
	$(TERRAFORM) output > tf.out ;\
	HOST=$$(awk '/host-ip/{print $$NF}' tf.out) ;\
	PORT=$$(awk '/server-port/{print $$NF}' tf.out) ;\
	if [ -z "$$HOST" ] || [ -z "$$PORT" ]; then \
		echo "Couldn't retrieve current host address or port. Are you sure the infrastructure is correctly deployed?" ;\
	else \
		echo "Testing API at http://$$HOST:$$PORT" ;\
		$(GO) test -v -race ./$(E2E) -host=$$HOST -port=$$PORT ;\
	fi ;\
	rm tf.out

terraform.keygen:
	ssh-keygen -t rsa -b 4096 -f $(TF_SSH_KEY_PATH) -N ""

terraform.init:
	$(TERRAFORM) init

terraform.chkfmt:
	$(TERRAFORM) fmt -check=true

terraform.validate:
	$(TERRAFORM) validate -var "srv-bin-path=$(PWD)/$(SRV_BIN_NAME)" -var "ssh-key-path=$(TF_SSH_KEY_PATH)"

terraform.apply:
	$(TERRAFORM) apply -var "srv-bin-path=$(PWD)/$(SRV_BIN_NAME)" -var "ssh-key-path=$(TF_SSH_KEY_PATH)" -input=false -auto-approve

terraform.output:
	$(TERRAFORM) output

terraform.destroy:
	$(TERRAFORM) destroy -var "srv-bin-path=$(PWD)/$(SRV_BIN_NAME)" -var "ssh-key-path=$(TF_SSH_KEY_PATH)" -auto-approve

clean:
	rm -f "$(PWD)/$(SRV_BIN_NAME)"
	rm -f $(TF_SSH_KEY_PATH)
	rm -f $(TF_SSH_KEY_PATH).pub

.PHONY: $(patsubst %,swagger.%,validate clean generate.client generate)
.PHONY: lint test.unit test.component test.e2e
.PHONY: $(patsubst %,terraform.%,keygen init chkfmt validate apply output destroy)
.PHONY: clean

.SILENT: test-e2e $(patsubst %,terraform.%,init chkfmt validate apply output destroy)
