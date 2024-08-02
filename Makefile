.PHONY: test t

t:
	go test ./...

test:
	go test ./... -v -count=1