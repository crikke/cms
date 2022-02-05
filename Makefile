SHELL := /bin/bash

.PHONY: test
test:
	go test ./... 

.PHONY: db-up
db-up:
	docker run --rm  -p 27017:27017 --name mongodb -d mongo

.PHONY: db-seed
db-seed:
	docker run --rm \
	--name mongo-seed \
	-v $(shell pwd)/hack/mongo-seed:/seed \
	--network host \
	--entrypoint /seed/seed-database.sh mongo
 