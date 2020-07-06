.PHONY: build

gen:
	docker run -v ${PWD}:/fetcher  openapitools/openapi-generator-cli generate -i /fetcher/api/openapi.yaml -g go-server -o /fetcher -c /fetcher/.openapi-generator/config.json --git-user-id jesse0michael --git-repo-id jessemichael.me --enable-post-process-file -t /fetcher/.openapi-generator/templates/go-server
	# GO_POST_PROCESS_FILE="gofmt -w"   

	docker run -v ${PWD}:/fetcher  openapitools/openapi-generator-cli generate -i /fetcher/api/openapi.yaml -g typescript-node -o /fetcher/client/ts/ --git-user-id jesse0michael --git-repo-id jessemichael.me --additional-properties=npmName=@jesse0michael/jessemichael.me,npmVersion=1.0.0

gen-docs:
	redoc-cli bundle api/openapi.yaml -o docs/index.html --options.disableSearch --options.hideDownloadButton

build:
	go build -o bin/fetcher ./cmd/fetcher/
	GOOS=linux go build -o bin/lambda ./cmd/lambda/
	cd bin && zip function.zip lambda

build-docker:
	docker build -t fetcher  --file ./build/Dockerfile .

run:
	docker run -p 8080:8080 fetcher
