FROM golang:1.20 as builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -v -o /usr/local/bin/app .

FROM gcr.io/distroless/static

WORKDIR /usr/src/app
COPY --from=builder /usr/local/bin/app .
USER 65532:65532

ENTRYPOINT ["/usr/src/app/app"]