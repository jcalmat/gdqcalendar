API = api
GENERATOR = generator

VERBOSE			= 0
C				= $(if $(filter 1,$(VERBOSE)),,@) ## Conditional command display
M				= @echo "\033[0;35m▶\033[0m"

.PHONY: all
all: vendor api

.PHONY: api
api: ## Build api binary
	$(M) building executable api…
	$(C) cd cmd/api && go build -o ../../$(API)

.PHONY: generator
generator: ## Build generator binary
	$(M) building executable generator…
	$(C) cd cmd/generator && go build -o ../../$(GENERATOR)

.PHONY: vendor
vendor:
	$(M) running mod vendor…
	$(GO) mod vendor

.PHONY: tidy
tidy:
	$(M) running mod tidy…
	$(GO) mod tidy

.PHONY: test
test:
	$(M) running go test…
	$(GO) test -cover -race -v ./...

.PHONY: fmt
fmt:
	$(M) running mod fmt…
	$(GO) fmt ./...

.PHONY: clean
clean:
	$(M) cleaning binaries…
	rm -rf $(API) $(GENERATOR)
