TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=trendmicro.com
NAMESPACE=alicloudsecurity
NAME=alicloudsecurity
BINARY=terraform-provider-${NAME}
VERSION=0.4
OS_ARCH=darwin_arm64

default: install

build:
	go build -o ${BINARY}

release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test:  
	@go test $(TEST) -race -coverprofile=coverage.txt -covermode=atomic -timeout=30s -parallel=4
	go tool cover -html=coverage.txt 
	go tool cover -func coverage.txt 
	rm coverage.txt                   

testacc-local: 
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m -coverprofile=coverage.txt -covermode=atomic
	go tool cover -html=coverage.txt 
	go tool cover -func coverage.txt
	rm coverage.txt

testacc:
	TF_ACC=1 go test $(TEST) $(TESTARGS) -timeout 120m -coverprofile=count.out
	go tool cover -func=count.out
