.PHONY: all tidy build run clean

all: tidy test build run clean

tidy:
	@echo "[tidy]"
	@go mod tidy

test:
	@echo "Running tests..."
	@go test ./... -cover
	@mkdir -p build
	@go test -coverprofile=build/coverage.out ./... > /dev/null
	@if [ -f build/coverage.out ]; then \
		go tool cover -func=build/coverage.out | grep total | awk '{print "Total coverage: " $$3}'; \
		rm -f build/coverage.out; \
	fi

build:
	@echo "[build]"
	@mkdir -p build
	@go build -o ./build/rmt .

run:
	@echo "[run]"
	@./build/rmt

clean:
	@echo "[clean]"
	@rm -rf build