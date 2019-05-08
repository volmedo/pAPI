spec = pAPI-swagger.yaml
pkg_dir = pkg
cmd_dir = cmd

swagger_ver = v0.19.0
swagger = docker run --rm -e GOPATH=/go -v "$(PWD)":"$(PWD)" -w "$(PWD)" quay.io/goswagger/swagger:$(swagger_ver)

terraform_ssh_key_path = "$(PWD)/tf_ssh_key"
terraform_dir = terraform
terraform_ver = 0.11.13
terraform = docker run --rm -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY -v "$(PWD)":"$(PWD)" -w "$(PWD)/$(terraform_dir)" hashicorp/terraform:$(terraform_ver)

swagger.validate:
	$(swagger) validate $(spec)
	
swagger.clean:
	rm -r $(pkg_dir)/models
	rm -r $(pkg_dir)/restapi

swagger.generate.server: swagger.clean
	$(swagger) generate server --spec=$(spec) --template=stratoscale --target=$(pkg_dir)

swagger.generate.client:
	$(swagger) generate client --spec=$(spec) --template=stratoscale --target=$(pkg_dir) --skip-models

lint:
	golangci-lint run --no-config --skip-dirs "$(pkg_dir)/(client|models|restapi)" --disable unused

test:
	go test -v -race ./$(pkg_dir)/impl ./$(cmd_dir)/server

terraform.keygen:
	ssh-keygen -t rsa -b 4096 -f $(terraform_ssh_key_path) -N ""

terraform.init:
	@$(terraform) init

terraform.apply:
	@$(terraform) apply -var "ssh-key-path=$(terraform_ssh_key_path)" -input=false -auto-approve

terraform.destroy:
	@$(terraform) destroy -var "ssh-key-path=$(terraform_ssh_key_path)" -auto-approve

.PHONY: swagger.validate swagger.clean swagger.generate.client swagger.generate lint test terraform.keygen terraform.init terraform.apply terraform.destroy
