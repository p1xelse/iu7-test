FROM golang:1.19

WORKDIR /

COPY . /app


WORKDIR /app

RUN go mod download
RUN go mod tidy
RUN go build cmd/benchmark/main.go

EXPOSE 8081

CMD ["./main"]
# CMD ["/bin/bash"]