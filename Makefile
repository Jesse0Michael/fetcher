.PHONY: build

gen:
	docker run -v ${PWD}:/fetcher  openapitools/openapi-generator-cli generate -i /fetcher/api/openapi.yaml -g go-server -o /fetcher -c /fetcher/.openapi-generator/config.json --git-user-id jesse0michael --git-repo-id fetcher --enable-post-process-file -t /fetcher/.openapi-generator/templates/go-server
	# GO_POST_PROCESS_FILE="gofmt -w"   

	docker run -v ${PWD}:/fetcher  openapitools/openapi-generator-cli generate -i /fetcher/api/openapi.yaml -g typescript-node -o /fetcher/client/ts/ --git-user-id jesse0michael --git-repo-id fetcher --additional-properties=npmName=@jesse0michael/fetcher,npmVersion=1.0.0

gen-docs:
	redoc-cli bundle api/openapi.yaml -o docs/index.html --options.disableSearch --options.hideDownloadButton

build:
	go build -o bin/fetcher ./cmd/fetcher/
	GOOS=linux go build -o bin/lambda ./cmd/lambda/
	cd bin && zip function.zip lambda

build-docker:
	docker build -t fetcher  --file ./build/Dockerfile .

run:
	go run ./cmd/fetcher
	
run-docker:	
	docker run -p 80:80 fetcher
