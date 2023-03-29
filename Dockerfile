FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY src/*.go ./

RUN go build -o main ./src/main.go

EXPOSE 8080

CMD [ "/main" ]