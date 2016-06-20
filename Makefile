

test:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

clean:
	go clean
	rm -rf coverage.out
