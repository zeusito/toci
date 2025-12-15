BINARY_NAME=myapp
DB_URL=postgres://postgres:qwerty@localhost:5432/mydb?sslmode=disable

.PHONY: lint
.PHONY: clean
.PHONY: run

lint:
	golangci-lint run --fix --config=.golangci.yaml

run:
	go run ./cmd/main.go -config=resources/config.local.toml

test:
	go test -v ./... --race -count=1

build:
	CGO_ENABLED=0 go build -o ./out/${BINARY_NAME} ./cmd/main.go

clean:
	go clean
	rm -f ./out

docker-build:
	docker build -t ${BINARY_NAME} .

migrations-up:
	dbmate -u "${DB_URL}" up
migrations-down:
	dbmate -u "${DB_URL}" down
migrations-info:
	dbmate -u "${DB_URL}" status
migrations-drop:
	dbmate -u "${DB_URL}" drop

