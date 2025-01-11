# Go 관련 변수 설정
APP_NAME := server
CMD_DIR := ./cmd/server
PROTO_SRC=proto/chatproto.proto
PROTO_OUT=.
BUILD_DIR := ./build

.PHONY: all clean build run test fmt vet install

all: build

# 빌드 명령어
build:
	@echo "Compiling Protobuf files..."
	protoc --go_out=$(PROTO_OUT) --go-grpc_out=$(PROTO_OUT) $(PROTO_SRC)
	@echo "==> Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)

# 실행 명령어
run: build
	@echo "==> Running $(APP_NAME)..."
	@$(BUILD_DIR)/$(APP_NAME)

# 코드 포맷팅
fmt:
	@echo "==> Formatting code..."
	go fmt ./...

# 코드 분석
vet:
	@echo "==> Running go vet..."
	go vet ./...

# 의존성 설치
install:
	@echo "==> Installing dependencies..."
	go mod tidy

# 테스트 실행
test:
	@echo "==> Running tests..."
	go test -v ./...

# 빌드 정리
clean:
	@echo "Cleaning generated files..."
	rm -rf $(PROTO_OUT)/*.go
	@echo "==> Cleaning build directory..."
	rm -rf $(BUILD_DIR)