all: smalld

smalld:
	go build smalld.go && mv ./smalld $(GOPATH)/bin/

release: smalld
	docker build Dockerfile
