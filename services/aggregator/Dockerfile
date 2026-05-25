FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /aggregator ./cmd

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /
COPY --from=builder /aggregator /aggregator

USER nonroot:nonroot
ENTRYPOINT ["/aggregator"]
