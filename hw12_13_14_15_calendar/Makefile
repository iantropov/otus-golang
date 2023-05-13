BIN_CALENDAR := "./bin/calendar"
BIN_CALENDAR_SCHEDULER := "./bin/calendar_scheduler"
BIN_CALENDAR_SENDER := "./bin/calendar_sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build: build-calendar build-calendar-scheduler build-calendar-sender

build-calendar:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar

run-calendar: build-calendar
	$(BIN_CALENDAR) -config ./configs/calendar.sql.toml

build-calendar-scheduler:
	go build -v -o $(BIN_CALENDAR_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler

run-calendar-scheduler: build-calendar-scheduler
	$(BIN_CALENDAR_SCHEDULER) -config ./configs/calendar_scheduler.sql.toml

build-calendar-sender:
	go build -v -o $(BIN_CALENDAR_SENDER) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

run-calendar-sender: build-calendar-sender
	$(BIN_CALENDAR_SENDER) -config ./configs/calendar_sender.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version-calendar: build-calendar
	$(BIN_CALENDAR) version

version-calendar-scheduler: build-calendar-scheduler
	$(BIN_CALENDAR_SCHEDULER) version

version-calendar-sender: build-calendar-sender
	$(BIN_CALENDAR_SENDER) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-lint-deps
	golangci-lint run ./...

install-grpc-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate: install-grpc-deps
	mkdir -p pkg/event_service_v1
	protoc 	--proto_path=api --go_out=pkg/event_service_v1 --go_opt=paths=source_relative \
			--go-grpc_out=pkg/event_service_v1  --go-grpc_opt=paths=source_relative \
			EventService.proto


.PHONY: up
up:
	docker-compose up migrate_db
	docker-compose up -d --build

.PHONY: down
down:
	docker-compose down


.PHONY: functest
functest:
	set -e ;\
	docker-compose -f docker-compose.functest.yml up migrate_db ;\
	docker-compose -f docker-compose.functest.yml up --build -d ;\
	test_status_code=0 ;\
	docker-compose -f docker-compose.functest.yml run functest go test || test_status_code=$$? ;\
	docker-compose -f docker-compose.functest.yml down ;\
	exit $$test_status_code ;

.PHONY: functest-cleanup
functest-cleanup:
	docker-compose -f docker-compose.functest.yml down \
        --rmi local \
		--volumes \
		--remove-orphans \
		--timeout 60; \
  	docker-compose -f docker-compose.functest.yml rm -f

.PHONY: build build-calendar build-calendar-scheduler build-calendar-sender
.PHONY: run-calendar run-calendar-scheduler run-calendar-sender
.PHONY: version-calendar version-calendar-scheduler version-calendar-sender
.PHONY: build-img run-img test lint generate
