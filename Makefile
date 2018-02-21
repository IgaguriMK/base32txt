.PHONY: build
build: enc32 dec32

.PHONY: enc32
enc32:
	go build enc32.go

.PHONY: dec32
dec32:
	go build dec32.go

.PHONY: deps
deps:
	true
