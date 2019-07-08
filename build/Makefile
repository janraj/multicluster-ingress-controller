OWNER=janrajc
IMAGE_NAME=citrix-k8s-multicluster-ingress-controller
version=1.1.0
.PHONY: build
build:
	docker build -f Dockerfile -t $(IMAGE_NAME) .. -m 4g
	docker save -o $(IMAGE_NAME).tar $(IMAGE_NAME)
test:
	 go test -v ../pkg/... -coverprofile cp.txt 
	 go tool cover -html=cp.txt -o cover.html
run:
	cd ../pkg;go run main.go

coverage:
	 go test -v ../... -coverprofile=coverage.txt -covermode=atomic

publish: repo-login publish

repo-login:
	@docker login -u "$(QUAY_USERNAME)" -p "$(QUAY_PASSWORD)" quay.io

publish:
	@echo 'publish latest to $(DOCKER_REGISTRY)'
	docker tag  $(IMAGE_NAME):latest $(DOCKER_REGISTRY)/$(IMAGE_NAME):latest
	docker push $(DOCKER_REGISTRY)/$(IMAGE_NAME):latest
	docker tag  $(IMAGE_NAME):latest $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(version)
	docker push $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(version)
proto:
	protoc -I/usr/local/include -I. \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=./../model/v1/ --proto_path=./../model/v1/  controlplane.proto
	swagger generate spec -i ./../cmd/k8s-control-plane/main.go -i ./../handler/api/v1/server.go  -o ./swagger.json
clean:
	docker rmi -f $$(docker images -q -f dangling=true) || true
	docker rmi -f $$(docker images | awk '$$1 ~ /$(IMAGE_NAME)/ { print $$3}') || true
         
