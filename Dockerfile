FROM golang:1.25-alpine AS deps
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM deps AS lint
RUN apk add --no-cache git
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
COPY . .
RUN golanci-lint run timeout=5m

FROM deps AS test
COPY . .
RUN go test -v ./...

FROM deps AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" chaosboard .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /app/chaosboard /chaosboard
EXPOSE 8080
ENTRYPOINT ["/chaosboard"]