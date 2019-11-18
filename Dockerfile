# builder image
FROM golang:1.13 AS builder
ENV GO111MODULE=on

RUN groupadd -g 10000 app && useradd -m -u 10001 -g app app

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix nocgo -o /dcc cmd/dcc/*.go

# go image
FROM alpine:3.10
LABEL name="dockerconfig-controller" maintainer="Kristoffer Johansson <kristoffer.johansson@gmx.com>"

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /dcc ./

USER app

ENTRYPOINT ["./dcc"]