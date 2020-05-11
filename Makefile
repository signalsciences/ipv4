
build:  ## compile everything
	go build ./...

lint:  ## lint time
	./scripts/lint.sh

test:  ## run all unit tess
	go test

bench: ## run all benchmarks
	go test -bench=. -benchmem

coverage:  ## test with coverage
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

clean:  ## cleanup
	./scripts/clean.sh

# https://www.client9.com/self-documenting-makefiles/
help:
	@awk -F ':|##' '/^[^\t].+?:.*?##/ {\
		printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF \
		}' $(MAKEFILE_LIST)
.DEFAULT_GOAL=help
.PHONY=help

