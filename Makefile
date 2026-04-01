BINARY_NAME=ws

.PHONY: all
all: clean build

.PHONY: build
build:
	@echo "==> Building $(BINARY_NAME) binary..."
	go build -o bin/$(BINARY_NAME) ./cmd/ws/

.PHONY: clean
clean:
	@echo "==> Cleaning build artifacts..."
	rm -rf bin/

.PHONY: run
run: build
	./bin/$(BINARY_NAME)

.PHONY: install
install: clean build
	@echo "==> Installing $(BINARY_NAME) binary..."
	mv bin/$(BINARY_NAME) ~/.local/bin/$(BINARY_NAME)
