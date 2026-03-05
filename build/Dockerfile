FROM golang:1.26-alpine AS build

# Copy project files
WORKDIR /go/src
COPY go.mod go.sum ./

# Fetch dependencies
RUN go mod download
COPY . .

# Build project
ENV CGO_ENABLED=0
RUN go build -o fetcher ./cmd/fetcher

FROM gcr.io/distroless/static

COPY --from=build /go/src/fetcher ./

ENTRYPOINT ["./fetcher"]
