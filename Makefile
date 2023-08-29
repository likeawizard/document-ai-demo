clean:
	go clean

build-native:
	GOAMD64=v3 go build -o expense-bot cmd/expense-bot/main.go

build: clean build-native