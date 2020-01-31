FROM golang:1.13-alpine AS build

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

# Copy project files
WORKDIR /go/src
COPY . .

# Fetch dependencies (go mod)
ENV GO111MODULE=on
RUN go mod download
RUN go mod verify

# Build project
ENV CGO_ENABLED=0
RUN go build -a -installsuffix cgo -o fetcher .

FROM scratch AS runtime

# Copy dependent files
COPY --from=build /go/src/fetcher ./
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080/tcp

ENTRYPOINT ["./fetcher"]
