LDFLAGS := "-X github.com/drsherluck/gobar/modules.ApiKey=${OWM_API_KEY}"

fmt:
	go fmt github.com/drsherluck/gobar/...

run: fmt
	go run -ldflags=${LDFLAGS} bar.go

build: fmt
	go build -ldflags=${LDFLAGS}

install: fmt
	sudo go build -ldflags=${LDFLAGS} -o /usr/bin/
