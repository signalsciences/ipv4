all: lint test clean

lint:
	gometalinter \
		--vendor \
		--vendored-linters \
		--deadline=60s \
		--disable-all \
		--enable=goimports \
		--enable=aligncheck \
		--enable=vetshadow \
		--enable=varcheck \
		--enable=structcheck \
		--enable=deadcode \
		--enable=ineffassign \
		--enable=unconvert \
		--enable=goconst \
		--enable=golint \
		--enable=gofmt \
		--enable=errcheck \
		--enable=misspell \
		./...

test:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

clean:
	go clean
	rm -rf coverage.out
