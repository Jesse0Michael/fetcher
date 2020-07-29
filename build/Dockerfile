FROM golang:1.14-alpine AS build

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

# Copy project files
WORKDIR /go/src
COPY go.mod .
COPY go.sum .
RUN mkdir /dir

# Fetch dependencies
RUN go mod download
COPY . .

# Build project
ENV CGO_ENABLED=0
RUN go build -o fetcher ./cmd/fetcher

FROM scratch AS runtime

# Copy dependent files
COPY --from=build /go/src/fetcher ./
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV PORT=8080

ENTRYPOINT ["./fetcher"]
