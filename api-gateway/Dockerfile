FROM golang:1.17.5-alpine3.15

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o main

EXPOSE ${PORT}

CMD ["go", "run", "main.go"]