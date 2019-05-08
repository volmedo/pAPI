.DEFAULT_GOAL := build

SPEC = pAPI-swagger.yaml
PKG = pkg
CMD = cmd

SWAGGER_VER = v0.19.0
SWAGGER = docker run --rm -e GOPATH=/go -v "$(PWD)":"$(PWD)" -w "$(PWD)" quay.io/goswagger/swagger:$(SWAGGER_VER)

TF_SSH_KEY_PATH = "$(PWD)/tf_ssh_key"
TF_DIR = terraform
TF_VER = 0.11.13
TERRAFORM = docker run --rm -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY -v "$(PWD)":"$(PWD)" -w "$(PWD)/$(TF_DIR)" hashicorp/terraform:$(TF_VER)

SRV_BIN_NAME = papisrv

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
	golangci-lint run --no-config --skip-dirs "$(PKG)/(client|models|restapi)" --disable unused

test:
	go test -v -race ./$(PKG)/impl ./$(CMD)/server

build: ./$(CMD)/server/main.go
	go build -o $(SRV_BIN_NAME) ./$(CMD)/server/main.go

terraform.keygen:
	ssh-keygen -t rsa -b 4096 -f $(TF_SSH_KEY_PATH) -N ""

terraform.init:
	@$(TERRAFORM) init

terraform.chkfmt:
	@$(TERRAFORM) fmt -check=true

terraform.validate:
	@$(TERRAFORM) validate -var "srv-bin-path=$(PWD)/$(SRV_BIN_NAME)" -var "ssh-key-path=$(TF_SSH_KEY_PATH)"

terraform.apply:
	@$(TERRAFORM) apply -var "srv-bin-path=$(PWD)/$(SRV_BIN_NAME)" -var "ssh-key-path=$(TF_SSH_KEY_PATH)" -input=false -auto-approve

terraform.destroy:
	@$(TERRAFORM) destroy -var "srv-bin-path=$(PWD)/$(SRV_BIN_NAME)" -var "ssh-key-path=$(TF_SSH_KEY_PATH)" -auto-approve

.PHONY: swagger.validate swagger.clean swagger.generate.client swagger.generate lint test terraform.keygen terraform.init terraform.apply terraform.destroy
