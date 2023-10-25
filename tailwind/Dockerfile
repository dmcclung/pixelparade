FROM alpine

RUN apk add --update nodejs npm

RUN mkdir /tailwind

WORKDIR /tailwind

RUN npm install -D tailwindcss
RUN npx tailwindcss init

CMD npx tailwindcss -c /src/tailwind.config.js -i /src/styles.css -o /dst/styles.css --watch --poll
