FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

COPY person-service.yaml ./person-service.yaml
COPY pkg/ ./pkg
COPY "cmd/" "./cmd"
COPY internal ./internal

RUN go mod tidy

WORKDIR /app/cmd/app

RUN CGO_ENABLED=0 go build -o /server .

FROM scratch AS run

COPY --from=build /server /server

ENTRYPOINT ["/server"]
