# Get the latest go image
FROM docker.io/golang:1.22.1-alpine AS builder
WORKDIR /work 
# COPY go.mod go.sum ./
COPY go.mod ./
# RUN go mod download -x
COPY . ./ 
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
# RUN CGO_ENABLED=0 GOOS=linux go build -a -o main github.com/myhops/cmd/authorizer
RUN CGO_ENABLED=0 GOOS=linux go build -a -o authorizer ./cmd/authorizer

FROM scratch
COPY --from=builder /work/authorizer /authorizer
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
ENV port=8080
ENTRYPOINT ["/authorizer"]
