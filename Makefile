spec = pAPI-swagger.yaml
target_dir = pkg
swagger_ver = v0.19.0
swagger = docker run --rm -e GOPATH=/go -v "$(PWD)":"$(PWD)" -w "$(PWD)" quay.io/goswagger/swagger:$(swagger_ver)

swagger.clean:
	rm -r $(target_dir)/models
	rm -r $(target_dir)/restapi

swagger.generate: swagger.clean
	if [ ! -d $(target_dir) ]; then mkdir -p $(target_dir); fi
	$(swagger) generate server --spec=$(spec) --template=stratoscale --target=$(target_dir)

swagger.validate:
	$(swagger) validate $(spec)

test:
	go test ./...

.PHONY: swagger.clean swagger.generate swagger.validate test