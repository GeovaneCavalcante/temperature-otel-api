run-server-a:
	go run cmd/server.go server-a

run-server-b:
	go run cmd/server.go server-b

docker-build-image:
	docker build -t geovanedeveloper/temperature-api:latest -f Dockerfile .

docker-up:
	docker-compose up -d

open-bash: 
	docker-compose exec temperature-app bash

build-mocks:
	go install go.uber.org/mock/mockgen@latest
	~/go/bin/mockgen -source=internal/entity/temperature.go -destination=internal/usecase/temperature/mock/find_temperature.go
	~/go/bin/mockgen -source=pkg/address/address.go -destination=pkg/address/mock/address.go
	~/go/bin/mockgen -source=pkg/temperature/temperature.go -destination=pkg/temperature/mock/temperature.go
	~/go/bin/mockgen -source=pkg/requester/requester.go -destination=pkg/requester/mock/requester.go
	
test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html