FROM golang:1.18 as builder
ENV GO111MODULE=on
WORKDIR /internal
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ./cmd/go-chat-client
RUN wget https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh && chmod +x wait-for-it.sh
CMD ["./wait-for-it.sh", "redis:6379", "--timeout=15", "--", "./go-chat-client"]