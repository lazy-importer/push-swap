.PHONY: all build test bench clean

all: build

# The subject requires two executables named checker and push-swap.
build:
	go build -trimpath -o push-swap ./cmd/push-swap
	go build -trimpath -o checker ./cmd/checker

test:
	go test ./...

bench:
	go test -run '^$$' -bench . -benchmem ./internal/pushswap

clean:
	rm -f push-swap checker
