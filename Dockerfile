FROM alpine AS tailwind-builder
RUN apk add --update nodejs npm
RUN mkdir /tailwind
WORKDIR /tailwind
RUN npm install -D tailwindcss
RUN npx tailwindcss init

COPY ./templates /templates
COPY ./tailwind/tailwind.config.js /src/tailwind.config.js
COPY ./tailwind/styles.css /src/styles.css

RUN npx tailwindcss -c /src/tailwind.config.js -i /src/styles.css -o /styles.css --minify

FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o ./server ./cmd/server/

FROM alpine
WORKDIR /app
COPY ./assets ./assets
COPY .env .env
COPY --from=tailwind-builder /styles.css ./assets/styles.css
COPY --from=builder /app/server ./server
CMD ["./server"]