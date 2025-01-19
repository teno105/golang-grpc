아래는 실습 순서에 맞춰 다시 작성한 `README.md`입니다.

---

# gRPC

`gRPC`는 Golang으로 작성된 gRPC 채팅프로그램의 구조와 테스트 방법을 익히기 위한 실습입니다.

## 실습 준비

### 1. 패키지 구조를 위한 디렉토리 생성

먼저 프로젝트 디렉터리를 설정하고 필요한 디렉터리들을 생성합니다.

```bash
mkdir golang-grpc
cd golang-grpc
go mod init golang-grpc

mkdir -p cmd/server
mkdir -p cmd/client
mkdir -p proto
```

### 2. `Makefile` 작성

이제 프로젝트의 빌드 및 실행을 자동화하기 위한 `Makefile`을 프로젝트 루트에 작성합니다.

```makefile
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
```

`Makefile`을 이용하여 코드를 빌드하고 실행할 수 있습니다.

```bash
make run
```

이 명령어를 통해 `server.go`에서 작성한 grpc 서버를 실행할 수 있습니다.

### 30.4 gRPC란?
gRPC란 구글에서 만든 오픈 소스 원격 프로시져 콜(Remote Procedure Call) 프레임워크입니다.<br/>
RPC는 네트워크를 통해서 다른 컴퓨터에서 원하는 함수(또는 기능)을 실행하는 것입니다.<br/>
gRPC는 이것을 편하게 하는 프레임워크로 사용이 쉽고 성능이 빠르기 때문에, 가장 많이 사용되는 RPC 프레임워크입니다.<br/>
gRPC개요: https://cloud.google.com/api-gateway/docs/grpc-overview?hl=ko

gRPC를 사용하기 때문에 gRPC 라이브러리를 설치합니다.
```bash
go get -u google.golang.org/grpc
```
### 30.4.1 프로토버퍼(Protobuf)
gRPC는 내부에서 메시지를 직렬화(Serialize)하려면 프로토버퍼 컴파일러를 사용하고 있습니다.<br/>
직렬화라는 것은 구조체 형태의 데이터를 하나의 바이너리 배열로 바꾸는 과정이라고 보면 됩니다.<br/>
역직렬화(Deserialize)는 다시 바이너리 배열을 구조체 형태의 데이터로 역변환하는 과정입니다.<br/>
<br/>
프로토버퍼는 구조체 데이터를 정의하고 그 데이터를 직렬화/역직렬화하는 기능을 제공합니다.

### 30.4.2 Linux에서 설치
```bash
sudo apt install -y protobuf-compiler
```

### 30.4.3 Mac에서 설치
```bash
brew install protobuf
```

### 30.4.5 protoc-gen-go 설치
protoc 설치 이후 go 파일 생성 시 사용되는 protoc-gen-go 패키지를 설치해야 합니다.
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
* 주의
protoc-gen-go 의 PATH를 못 찾는 경우가 있으므로 환경변수에 PATH를 추가해야합니다.
```bash
export PATH="$HOME/go/bin:$PATH"
```
### 30.5 gRPC를 이용한 채팅 프로그램
### 30.5.1 제작 과정
gRPC를 이용해서 채팅 프로그램을 만들기 위해서는 먼저 클라이언트와 서버 간 주고받을 데이터를 정의하는 서비스 정의 파일을 만들어야 합니다.<br/>
이 서비스 정의 파일은 프로토버퍼에 정한 형식을 따라야 합니다.<br/>
서비스 정의파일에 정의한 서버스에 해당하는 go 코드를 생성하면, 클라이언트와 서버에 사용할 수 있는 인터페이스 형식이 포함되어 있습니다.<br/>
그래서 우리는 이 인터페이스 형식에 맞는 클라이언트와 서버코드만 만들어주면 프로그램이 완성됩니다.

### 30.5.2 폴더 구조
```plaintext
golang-grpc/
│
├── cmd/
│   ├── client/
│   │    └── client.go
│   │
│   └── server/
│        └── server.go
│
├── proto/
│   └── chatproto.proto
│
├── go.mod
├── Makefile
└── README.md
```

### 30.5.3 서비스 정의 파일
proto 폴더 안에 서비스 정의 파일을 만들어 줍니다.
```proto
// proto/chatproto.proto
syntax = "proto3";

// 1. 패키지 이름이 들어갑니다.
option go_package = "pkg/server/generated;generated";

package chatproto;

// 2. 서비스 정의입니다. Chat() 함수를 포함하고 있습니다.
service ChatService {
  rpc Chat (stream ChatMsg) returns (stream ChatMsg) {}
}

// 3. Chat 기능에 사용되는 구조체 정의입니다.
message ChatMsg {
  string sender = 1;
  string message = 2;
}
```
2. 서비스 정의입니다.
```go
rpc Chat (stream ChatMsg) returns (stream ChatMsg) {}
```
첫번째 rpc는 서버에서 실행되는 함수임을 나타내는 키워드입니다.<br/>
그다음 함수 이름이 나옵니다. 이름은 Chat()입니다.<br/>
그다음은 입력 인수를 정의합니다. 입력인수는 ChatMsg라는 구조체가 입력됩니다.<br/>
ChatMsg는 하단에 정의 하고 있습니다. ChatMsg는 stream이라는 키워드가 붙어 있습니다. 이것은 ChatMsg 입력이 스트림 현태로 연속적으로 들어올 수 있음을 나타냅니다.<br/>
return 키워드를 쓰고 출력 형태를 정의합니다. 출력 역시 ChatMsg라는 구조체 형태로 출력됩니다.<br/>
<br/>
3. Chat 기능에 사용되는 ChatMsg 구조체를 정의 합니다.<br/>
message 라는 키워드를 통해서 메시지 정의임을 나타냅니다.<br/>
ChatMsg 구조체에는 두 개의 필드가 있는 문자열 타입이고 각각 sender와 message라는 이름을 가지고 있습니다.<br/>
```go
message ChatMsg {
  string sender = 1;
  string message = 2;
}
```
한 가지 특이한 점은 뒤에 =1, =2 와 같이 값이 붙어 있는데 이것은 해당 필드의 초깃값이 아니라 <br/>
해당 필드가 메시지 어디에 위치하는지를 나타내는 필드의 인덱스 입니다.
<br/>
```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative chatproto.proto
```
위의 명령어를 사용하면 generate 된 chatproto.pb.go 파일과 chatproto_grpc.pb.go 파일이 생성이 되는데,<br/>
이번 실습에서는 Makefile 에서 build시 생성되는 것을 사용하여 실습하겠습니다.

### 30.5.5 서버 구현

`cmd/server/` 디렉터리 아래에 `server.go` 파일을 생성해 구현합니다.

```go
// cmd/server/server.go
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	pb "golang-grpc/pkg/server/generated"

	"google.golang.org/grpc"
)

var port = flag.Int("port", 50051, "The server port")

func main() {
	flag.Parse()

	// 1. 연결을 기다립니다.
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	// 2. Chat 서비스를 등록합니다.
	pb.RegisterChatServiceServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}

func newServer() *chatServer {
	return &chatServer{}
}

// 3. Chat 서비스 인터페이스를 구현한 객체입니다.
type chatServer struct {
	pb.UnimplementedChatServiceServer
	mu      sync.Mutex
	streams []pb.ChatService_ChatServer
}

func (s *chatServer) Chat(stream pb.ChatService_ChatServer) error {
	// 4. stream 리스트에 추가합니다.
	s.mu.Lock()
	s.streams = append(s.streams, stream)
	s.mu.Unlock()

	var err error
	for {
		// 5. 클라이언트로 전송된 입력값을 읽습니다.
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		// 6. 전체 클라이언트로 방송합니다.
		s.mu.Lock()
		for _, strm := range s.streams {
			strm.Send(&pb.ChatMsg{
				Sender:  in.Sender,
				Message: in.Message,
			})
		}
		s.mu.Unlock()
	}

	// 7. 연결이 끊어졌기 때문에 리스트에서 삭제합니다.
	s.mu.Lock()
	for i, strm := range s.streams {
		if strm == stream {
			s.streams = append(s.streams[:i], s.streams[i+1:]...)
			break
		}
	}
	s.mu.Unlock()
	return err
}
```
2. gRPC 서버를 생성한다. grpc 서버의 구동은 ①서버 생성 - ②실행 단계로 구분된다. 서버를 생성할 때에는 NewServer()를, 실행할 때에는 Serve() 함수를 사용한다.
```go
pb.RegisterChatServiceServer(grpcServer, newServer())
```
서버를 생성한 후에는 grpc 서비스를 등록해야 하는데, 이때 어떤 API와 어떤 핸들러 함수를 연결할 것인지 매핑 정보를 등록해주는 단계가 필요하다. <br/>
pb.RegisterChatServiceServer() 함수를 호출하는 것이 바로 그 단계이다. 만약 grpc 서비스가 여러 개라면 각각의 서비스마다 등록해주어야 한다. 이 과정은 서버 실행에 앞서 처리되어야 한다. <br/>

이 코드를 작성 후, 아래의 bash 구문을 실행하면 chatproto.proto 가 generate가 되며, 서버가 실행됩니다.
```bash
make run
```

### 30.5.4 클라이언트 구현
```go
// cmd/client/client.go
package main

import (
	"bufio"
	"context"
	"flag"
	pb "golang-grpc/pkg/server/generated"
	"io"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)
// 1. 실행 인수를 정의합니다.
var id = flag.String("id", "unknown", "The id name")
var serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")

func main() {
	flag.Parse()
	// 2. grpc 서버에 연결합니다.
	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	// Chat 서비스 클라이언트를 실행합니다.
	client := pb.NewChatServiceClient(conn)

	runChat(client)
}

func runChat(client pb.ChatServiceClient) {
	// 4. Chat 기능을 호출합니다.
	stream, err := client.Chat(context.Background())
	if err != nil {
		log.Fatalf("client. Char failed: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			// 5. 출력 스트림으로 출력값이 나오면 화면에 출력합니다.
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("client. Chat failed: %v", err)
			}
			log.Printf("Sender:%s Message:%s", in.Sender, in.Message)
		}
	}()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		if strings.ToLower(msg) == "exit" {
			break
		}
		// 6. 키보드로 한 줄을 입력받아 입력으로 넣어줍니다.
		stream.Send(&pb.ChatMsg{
			Sender:  *id,
			Message: msg,
		})
	}
	stream.CloseSend()
	<-waitc
}
```

위 코드를 작성 후, build 후에 서버에 접속을 합니다.
```bash
cd cmd/client
go build client.go
./client -id=teno -addr=127.0.0.1:50051
```
### 실행화면
![스크린샷 2025-01-12 오후 6 53 24](https://github.com/user-attachments/assets/936a2ec2-1c6e-4635-ae4b-b518fa9b4da7)
