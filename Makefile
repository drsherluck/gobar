fmt:
	go fmt github.com/drsherluck/gobar/...

run: fmt
	go run bar.go
build: fmt
	go build 
install: fmt
	sudo go build -o /usr/bin/
