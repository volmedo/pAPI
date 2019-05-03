spec = pAPI-swagger.yaml
pkg_dir = pkg
cmd_dir = cmd
swagger_ver = v0.19.0
swagger = docker run --rm -e GOPATH=/go -v "$(PWD)":"$(PWD)" -w "$(PWD)" quay.io/goswagger/swagger:$(swagger_ver)

swagger.validate:
	$(swagger) validate $(spec)
	
swagger.clean:
	rm -r $(pkg_dir)/models
	rm -r $(pkg_dir)/restapi

swagger.generate.server: swagger.clean
	$(swagger) generate server --spec=$(spec) --template=stratoscale --target=$(pkg_dir)

swagger.generate.client:
	$(swagger) generate client --spec=$(spec) --template=stratoscale --target=$(pkg_dir) --skip-models

test:
	go test ./...

.PHONY: swagger.validate swagger.clean swagger.generate.client swagger.generate test
