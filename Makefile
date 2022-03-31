build:
	GOARCH=amd64 GOOS=linux go build -a -o bin/service ./application/main.go

deploy: build
	serverless deploy --stage $(STAGE) --verbose

build-local:
	docker-compose build service

run-local:
	docker-compose up

node-dependencies:
	@yarn install


