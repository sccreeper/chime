FROM alpine:latest

RUN apk update
RUN apk upgrade
RUN apk add --no-cache nodejs npm

WORKDIR /app

COPY package.json /app/
COPY package-lock.json /app/

RUN npm install

COPY . /app/

RUN npm run build

ENV PORT=3000
ENV ORIGIN=http://chime:8042

CMD [ "node", "/app/build" ]
