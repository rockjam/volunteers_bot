build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/service ./application/main.go

deploy:
	serverless deploy --stage $(STAGE) --verbose

build-local:
	docker-compose build service

run-local:
	docker-compose up

node-dependencies:
	@yarn install


