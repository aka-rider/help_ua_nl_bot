### Build
FROM golang AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o main ./...


## Deploy
FROM gcr.io/distroless/static

WORKDIR /app

USER nonroot:nonroot
COPY --from=build --chown=nonroot:nonroot /app /app

ENTRYPOINT ["/app/main"]