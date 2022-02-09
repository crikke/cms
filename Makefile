SHELL := /bin/bash

.PHONY: all
all: test vet 

.PHONY: test
test: db-seed
	go test ./... 

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
 