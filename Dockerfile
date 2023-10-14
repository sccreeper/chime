FROM alpine:latest

RUN apk update
RUN apk upgrade
RUN apk add --no-cache go ffmpeg nodejs npm curl

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./server ./server
COPY ./web ./web

# Words file used for generating user friendly server IDs
RUN curl https://www.mit.edu/~ecprice/wordlist.10000 > ./server/words.txt
RUN go build -o /chime ./server/.

WORKDIR /app/web

# Build web frontend.
RUN npm install
RUN npx vite build

RUN mkdir /dist/
RUN mkdir /dist/assets/
RUN cp ./dist/index.html /dist/index.html
RUN cp -r ./dist/assets/* /dist/assets/

WORKDIR /

CMD [ "/chime" ]