VERSION = 0.1
IMAGE = quay.io/eacmesquita96/go-postgres:v$(VERSION)

build:
	CGO_ENABLED=0 GOOS=linux go build main.go
	docker build -t $(IMAGE) .