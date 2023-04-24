FROM alpine:latest

RUN apk update
RUN apk upgrade
RUN apk add --no-cache go ffmpeg nodejs npm

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /chime ./server/.

WORKDIR /app/web

RUN npm install
RUN npx vite build

RUN mkdir /dist/
RUN mkdir /dist/assets/
RUN cp ./dist/index.html /dist/index.html
RUN cp -r ./dist/assets/* /dist/assets/

WORKDIR /

CMD [ "/chime" ]