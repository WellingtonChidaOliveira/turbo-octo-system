FROM golang:1.26 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /bin/processor ./cmd

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /bin/processor /processor

USER nonroot:nonroot
ENTRYPOINT ["/processor"]
