SHELL := /bin/bash

.PHONY: test
test:
	go test ./... 

.PHONY: db-up
db-up:
	docker run --rm --name mongodb -e MONGO_INITDB_ROOT_USERNAME=user -e MONGO_INITDB_ROOT_PASSWORD=strongpassword

.PHONY: db-seed
db-seed:
	docker run --rm $(docker build -q ./hack/mongo-seed/) -e USERNAME=user -e PASSWORD=strongpassword -e DATABASE_URI=0.0.0.0
 