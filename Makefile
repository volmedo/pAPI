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

POSTGRES_VER ?= 11.3-alpine
CONTAINER_NAME = db
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= papi_user
DB_PASS ?= papi_test_pass
DB_NAME ?= papi_db
POSTGRES_START = docker run --name $(CONTAINER_NAME) \
					-e POSTGRES_USER=$(DB_USER) \
					-e POSTGRES_PASSWORD=$(DB_PASS) \
					-e POSTGRES_DB=$(DB_NAME) \
					-p $(DB_PORT):5432 \
					-d postgres:$(POSTGRES_VER)
POSTGRES_WAIT = until docker run --rm --link $(CONTAINER_NAME):pg postgres:$(POSTGRES_VER) pg_isready -U postgres -h pg; do sleep 1; done
POSTGRES_STOP = docker stop $(CONTAINER_NAME) && docker rm $(CONTAINER_NAME)

TF_SSH_KEY_PATH ?= "$(PWD)/tf_ssh_key"
TF_DIR ?= terraform
TF_VER ?= 0.11.13
TERRAFORM = docker run --rm -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY -v "$(PWD)":"$(PWD)" -w "$(PWD)/$(TF_DIR)" hashicorp/terraform:$(TF_VER)

SRV_BIN_NAME ?= papisrv

PAPI_IMG_TAG ?= test

KIND_CLUSTER = papi

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
	golangci-lint run --no-config --skip-dirs "$(PKG)/(client|models|restapi)" --deadline 2m

test.unit:
	$(GO) test -v -race ./$(PKG)/service
	$(GO) test -v -race ./$(CMD)/server

# The exit code of the test command is saved in a variable to call POSTGRES_STOP
# no matter if tests fails or not but make the final exit code of the command depend
# on the result of tests, so that the step fails in CI if tests fail
test.integration:
	if [ -z $$CI ]; then \
		$(POSTGRES_START) ;\
		$(POSTGRES_WAIT) ;\
	fi ;\
	$(GO) test -v -race -tags=integration ./$(PKG)/service \
		-dbhost=$(DB_HOST) \
		-dbport=$(DB_PORT) \
		-dbuser=$(DB_USER) \
		-dbpass=$(DB_PASS) \
		-dbname=$(DB_NAME) \
		-migrations=./migrations ;\
	TEST_RESULT=$$? ;\
	if [ -z $$CI ]; then \
		$(POSTGRES_STOP) ;\
	fi ;\
	exit $$TEST_RESULT

build:
	GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) $(GO) build -o $(SRV_BIN_NAME) ./$(CMD)/server

test.e2e.local:
	$(GO) build -o testsrv ./$(CMD)/server
	if [ -z $$CI ]; then \
		$(POSTGRES_START) ;\
		$(POSTGRES_WAIT) ;\
	fi ;\
	./testsrv \
		-port=8080 \
		-dbhost=$(DB_HOST) \
		-dbport=$(DB_PORT) \
		-dbuser=$(DB_USER) \
		-dbpass=$(DB_PASS) \
		-dbname=$(DB_NAME) \
		-migrations=./$(PKG)/service/migrations & \
	SERVER_PID=$$! ;\
	$(GO) test -v -race ./$(E2E) -host=localhost -port=8080 ;\
	TEST_RESULT=$$? ;\
	kill $$SERVER_PID ;\
	rm testsrv ;\
	if [ -z $$CI ]; then \
		$(POSTGRES_STOP) ;\
	fi ;\
	exit $$TEST_RESULT

docker.build:
	docker build -t volmedo/papi:$(PAPI_IMG_TAG) .

docker.push:
	docker push volmedo/papi:$(PAPI_IMG_TAG)

test.e2e.k8s:
	kind create cluster --name $(KIND_CLUSTER) --wait 5m
	export KUBECONFIG="$$(kind get kubeconfig-path --name $(KIND_CLUSTER))" ;\
	kubectl apply -f k8s/ ;\
	kubectl wait --for condition=Ready pod -l tier=backend ;\
	PROXY_PORT=8000 ;\
	kubectl proxy --port=$$PROXY_PORT & \
	PROXY_PID=$$! ;\
	$(GO) test -v -race ./$(E2E) \
		-host=localhost \
		-port=$$PROXY_PORT \
		-api-path=/api/v1/namespaces/default/services/api/proxy/v1 \
		-health-path=/api/v1/namespaces/default/services/api/proxy/health ;\
	TEST_RESULT=$$? ;\
	kill $$PROXY_PID ;\
	kubectl delete -f k8s/ ;\
	unset KUBECONFIG ;\
	kind delete cluster --name $(KIND_CLUSTER) ;\
	exit $$TEST_RESULT

test.e2e:
	$(TERRAFORM) output > tf.out ;\
	HOST=$$(awk '/srv-ip/{print $$NF}' tf.out) ;\
	PORT=$$(awk '/srv-port/{print $$NF}' tf.out) ;\
	if [ -z "$$HOST" ] || [ -z "$$PORT" ]; then \
		echo "Couldn't retrieve current host address or port. Are you sure the infrastructure is correctly deployed?" ;\
	else \
		echo "Testing API at http://$$HOST:$$PORT" ;\
		$(GO) test -v -race ./$(E2E) -host=$$HOST -port=$$PORT ;\
		TEST_RESULT=$$? ;\
	fi ;\
	rm tf.out ;\
	exit $$TEST_RESULT

terraform.keygen:
	ssh-keygen -t rsa -b 4096 -f $(TF_SSH_KEY_PATH) -N ""

terraform.init:
	$(TERRAFORM) init

terraform.chkfmt:
	$(TERRAFORM) fmt -check=true

terraform.validate:
	$(TERRAFORM) validate \
		-var "srv-bin-path=$(PWD)/$(SRV_BIN_NAME)" \
		-var "ssh-key-path=$(TF_SSH_KEY_PATH)" \
		-var "db-name=$(DB_NAME)" \
		-var "db-port=$(DB_PORT)" \
		-var "db-user=$(DB_USER)" \
		-var "db-pass=$(DB_PASS)" \
		-var "db-migrations-path=$(PWD)/$(PKG)/service/migrations"

terraform.apply:
	$(TERRAFORM) apply \
		-var "srv-bin-path=$(PWD)/$(SRV_BIN_NAME)" \
		-var "ssh-key-path=$(TF_SSH_KEY_PATH)" \
		-var "db-name=$(DB_NAME)" \
		-var "db-port=$(DB_PORT)" \
		-var "db-user=$(DB_USER)" \
		-var "db-pass=$(DB_PASS)" \
		-var "db-migrations-path=$(PWD)/$(PKG)/service/migrations" \
		-input=false \
		-auto-approve

terraform.output:
	$(TERRAFORM) output

terraform.destroy:
	$(TERRAFORM) destroy \
		-var "srv-bin-path=$(PWD)/$(SRV_BIN_NAME)" \
		-var "ssh-key-path=$(TF_SSH_KEY_PATH)" \
		-var "db-name=$(DB_NAME)" \
		-var "db-port=$(DB_PORT)" \
		-var "db-user=$(DB_USER)" \
		-var "db-pass=$(DB_PASS)" \
		-var "db-migrations-path=$(PWD)/$(PKG)/service/migrations" \
		-auto-approve

clean:
	rm -f "$(PWD)/$(SRV_BIN_NAME)"
	rm -f $(TF_SSH_KEY_PATH)
	rm -f $(TF_SSH_KEY_PATH).pub

.PHONY: $(patsubst %,swagger.%,validate clean generate.client generate.server)
.PHONY: lint
.PHONY: $(patsubst %,test.%,unit integration e2e.local e2e.k8s e2e)
.PHONY: $(patsubst %,docker.%,build push)
.PHONY: $(patsubst %,terraform.%,keygen init chkfmt validate apply output destroy)
.PHONY: clean

.SILENT: $(patsubst %,terraform.%,init chkfmt validate apply output destroy)
