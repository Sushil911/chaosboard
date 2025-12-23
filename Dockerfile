FROM golang:1.25-alpine AS deps
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM deps AS lint
RUN apk add --no-cache git
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
COPY . .
RUN golangci-lint run --timeout=5m

FROM deps AS test
COPY . .
RUN go test -v ./...

FROM deps AS development
RUN go install github.com/air-verse/air@latest && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && go install github.com/go-delve/delve/cmd/dlv@latest
CMD ["air","-c",".air.toml"]

FROM deps AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" chaosboard .

FROM gcr.io/distroless/static-debian12:nonroot AS production
COPY --from=builder /app/chaosboard /chaosboard
EXPOSE 8080
ENTRYPOINT ["/chaosboard"]