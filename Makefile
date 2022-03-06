SHELL := /bin/bash
RABBITMQ_HOSTNAME := cms_rabbitmq
IMAGE_NAME := contentdelivery
IMAGE_TAG := local

.PHONY: test
test: 
	go test ./... -tags=unit

.PHONY: integration
integration: db-up
	go test ./... -tags=integration -p 1

.PHONY: vet
vet:
	go vet ./...

.PHONY: db-up
db-up:
	docker start mongodb 2>/dev/null|| docker run --rm  -p 27017:27017 --name mongodb -d mongo

.PHONY: db-seed
db-seed: db-up
	docker run --rm \
	--name mongo-seed \
	-v $(shell pwd)/hack/mongo-seed:/seed \
	--network host \
	--entrypoint /seed/seed-database.sh mongo
 


 .PHONY: rabbitmq
 rabbitmq: db-seed
	# docker build . -t $(IMAGE_NAME):$(IMAGE_TAG) 

	docker run --rm -d \
	-p 5672:5672 \
	--hostname $(RABBITMQ_HOSTNAME) \
	rabbitmq  

	# docker run --rm $(IMAGE_NAME):$(IMAGE_TAG)

# todo
# docker run -p 8081:8080 --user $(id -u):$(id -g) -e SWAGGER_JSON=/swagger.json -v $(pwd)/swagger.json:/swagger.json swaggerapi/swagger-ui


# .PHONY: swagger
# swagger:
# 	docker run --rm -it -e GOPATH=$(go env GOPATH):/go -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger'
