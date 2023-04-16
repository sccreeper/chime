FROM alpine:latest

RUN apk add --no-cache go ffmpeg

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /chime ./server/.

CMD [ "/chime" ]