FROM golang

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o ./pixelparade ./cmd/server/

CMD ["./pixelparade"]