LDFLAGS="-X github.com/drsherluck/gobar/modules.ApiKey=${OWM_API_KEY}"

run:
	go run -ldflags=${LDFLAGS} bar.go

build:
	go build -ldflags=${LDFLAGS}

install:
	sudo go build -ldflags=${LDFLAGS} -o /usr/bin/
