build:
	docker build -t imagemaker-app .
	docker run --publish 8082:8000 --rm --name imagemaker-app-run imagemaker-app
	explorer "http://localhost:8082"

test:
	go test -v ./...

clean:
	rm -rf build

.PHONY: build test lint clean
