.PHONY: all test coverage

all: test bench profile coverage

test: ## Run tests
	@go test -v ./...

bench: ## Run benchmarks
	@go test -run=XXX -bench=. ./...

cpu.prof:
	@go test -run=XXX -bench=. ./... -cpuprofile=$@
cpu.pdf: cpu.prof
	@go tool pprof -pdf -output $@ $<
profile: cpu.pdf ## Run benchmarks, gen profile (requires graphviz)

coverage: ## Generate global code coverage report
	@go test -covermode=count ./...