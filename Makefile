PRODUCER_MAIN_PACKAGE_PATH := cmd/producer
PRODUCER_BINARY_NAME := producer
CONSUMER_MAIN_PACKAGE_PATH := cmd/consumer
CONSUMER_BINARY_NAME := consumer


PROTO_DIR := proto
PB_DIR := proto/pb
PROTOC := protoc
GRPC_PLUGIN := protoc-gen-go
GRPC_GATEWAY_PLUGIN := protoc-gen-grpc-gateway
PROTOC_OPTS := -I$(PROTO_DIR) --go_out=$(PB_DIR) --go_opt=paths=source_relative --go-grpc_out=$(PB_DIR) --go-grpc_opt=paths=source_relative

# ==================================================================================== #
# HELPERS
# ==================================================================================== #
## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	git diff --exit-code

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #
## format: format code
.PHONY: format
format:
	find . -name '*.go' -exec gofumpt -w {} +

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

## linter: run linters
.PHONY: linter-golangci
linter-golangci: ### check by golangci linter
	golangci-lint run

## clean: clean-up
.PHONY: clean
clean:
	go clean

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #
## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## build: build the application
.PHONY: build
build:
	# Include additional build steps, like TypeScript, SCSS, or Tailwind compilation here...
	go build -o=/tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## run: run the application
.PHONY: run
run: build
	/tmp/bin/${BINARY_NAME}

## run/live: run the application with reloading on file changes
.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" --build.bin "/tmp/bin/${BINARY_NAME}" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go,tpl,tmpl,html,css,scss,js,ts,sql,jpeg,jpg,gif,png,bmp,svg,webp,ico" \
		--misc.clean_on_exit "true"

# ==================================================================================== #
# OPERATIONS
# ==================================================================================== #
## push: push changes to the remote Git repository
.PHONY: push
push: tidy audit no-dirty
	git push

## proto: generate protobuf files
.PHONY: generate-proto
generate-proto:
	$(PROTOC) $(PROTOC_OPTS) $(PROTO_DIR)/*.proto

.PHONY: clean-proto
clean-proto:
	rm proto/pb/*.pb.go;

# ==================================================================================== #
# RUN
# ==================================================================================== #
## run: run the applications
run: run-ignored run-unparsed run-wrong

.PHONY: run-employee
run-employee:
	cd $(PRODUCER_MAIN_PACKAGE_PATH) && go mod tidy && go mod download && \
    CGO_ENABLED=0 go run rc-boxdata_storage/$(PRODUCER_MAIN_PACKAGE_PATH)

.PHONY: run-ignored
run-ignored:
	cd $(IGNORED_DATA_MAIN_PACKAGE_PATH) && go mod tidy && go mod download && \
    CGO_ENABLED=0 go run rc-boxdata_storage/$(IGNORED_DATA_MAIN_PACKAGE_PATH)

.PHONY: run-company
run-company:
	cd $(COMPANY_MAIN_PACKAGE_PATH) && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run rc-boxdata_storage/$(COMPANY_MAIN_PACKAGE_PATH)

.PHONY: run-unparsed
run-unparsed:
	cd $(CONSUMER_MAIN_PACKAGE_PATH) && go mod tidy && go mod download && \
    CGO_ENABLED=0 go run rc-boxdata_storage/$(CONSUMER_MAIN_PACKAGE_PATH)

.PHONY: run-wrong
run-wrong:
	cd $(WRONG_DATA_MAIN_PACKAGE_PATH) && go mod tidy && go mod download && \
    CGO_ENABLED=0 go run rc-boxdata_storage/$(WRONG_DATA_MAIN_PACKAGE_PATH)

.PHONY: run-vehicle
run-vehicle:
	cd $(VEHICLE_MAIN_PACKAGE_PATH) && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run rc-boxdata_storage/$(VEHICLE_MAIN_PACKAGE_PATH)
## docker-compose: run docker-compose
docker-compose: docker-compose-stop docker-compose-start
.PHONY: docker-compose

.PHONY: docker-compose-start
docker-compose-start:
	docker-compose up --build

.PHONY: docker-compose-start-dependency
docker-compose-start-dependency:
	docker-compose up --build mongodb rabbitmq

.PHONY: docker-compose-start-unparsed
docker-compose-start-unparsed:
	docker-compose up --build unparsed

.PHONY: docker-compose-start-employee
docker-compose-start-employee:
	docker-compose up --build employee

.PHONY: docker-compose-start-wrong
docker-compose-start-wrong:
	docker-compose up --build wrong

.PHONY: docker-compose-start-ignored
docker-compose-start-ignored:
	docker-compose up --build ignored

.PHONY: docker-compose-stop
docker-compose-stop:
	docker-compose down --remove-orphans -v

# ==================================================================================== #
# WIRE
# ==================================================================================== #
## wire: generate wire
.PHONY: wire
wire:
	cd internal/ignored_data/app && wire && cd - && \
	cd internal/unparsed_boxdata/app && wire && cd - && \
	cd internal/employee/app && wire && cd - && \
	cd internal/wrong_data/app && wire && cd -
	cd internal/vehicle/app && wire && cd -
	cd internal/company/app && wire && cd -

.PHONY: clean-wire
clean-wire:
	cd internal/unparsed_boxdata/app && rm wire_gen.go && cd - && \
	cd internal/employee/app && rm wire_gen.go && cd - && \
	cd internal/ignored_data/app && rm wire_gen.go && cd - && \
	cd internal/wrong_data/app && rm wire_gen.go && cd -
	cd internal/vehicle/app && rm wire_gen.go && cd -
	cd internal/company/app && rm wire_gen.go && cd -
