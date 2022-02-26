SHELL := /bin/bash
RABBITMQ_HOSTNAME := cms_rabbitmq
IMAGE_NAME := contentdelivery
IMAGE_TAG := local

.PHONY: all
all: test vet 

.PHONY: test
test: db-seed
	go test ./... -cover 

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

