FROM golang:1.23
LABEL authors="MielPops"
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o main ./cmd/app/main.go

CMD ["./main"]