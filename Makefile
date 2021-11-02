.PHONY: build
COVERAGEDIR = .coverage

gen:
	go generate ./...

	docker run -v ${PWD}:/fetcher  openapitools/openapi-generator-cli generate -i /fetcher/api/openapi.yaml -g typescript-node -o /fetcher/client/ts/ --git-user-id jesse0michael --git-repo-id fetcher --additional-properties=npmName=@jesse0michael/fetcher,npmVersion=1.0.0

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
