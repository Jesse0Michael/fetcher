.PHONY: build
COVERAGEDIR = .coverage

gen:
	go generate ./...

gen-docs:
	redoc-cli bundle api/openapi.yaml -o docs/index.html --options.disableSearch --options.hideDownloadButton

build:
	go build -o bin/fetcher ./cmd/fetcher/
	GOOS=linux go build -o bin/lambda ./cmd/lambda/
	cd bin && zip function.zip lambda

build-docker:
	docker build -t fetcher  --file ./build/Dockerfile .


test:
	go test -cover ./internal/... 
	golangci-lint run ./internal/...

test-coverage:
	if [ ! -d $(COVERAGEDIR) ]; then mkdir $(COVERAGEDIR); fi
	go test -coverpkg ./internal/... -coverprofile $(COVERAGEDIR)/request.coverprofile ./internal/...
	go tool cover -html $(COVERAGEDIR)/request.coverprofile


run:
	go run ./cmd/fetcher
	
run-docker:	
	docker run -p 8080:8080 fetcher
