build:
	GOARCH=amd64 GOOS=linux go build -a -o bin/service ./application/main.go

deploy:
	serverless deploy --verbose

build-local:
	docker-compose build service

run-local:
	docker-compose up

node-dependencies:
	npm ci


