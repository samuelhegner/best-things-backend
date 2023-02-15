FROM golang:alpine

EXPOSE 8080

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o /best-things-api

CMD [ "/best-things-api" ]