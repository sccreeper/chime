FROM alpine:latest

RUN apk update
RUN apk upgrade
RUN apk add python3 py3-gunicorn poetry

COPY . ./

RUN poetry install

CMD ["poetry", "run", "python3", "-m", "gunicorn", "-b", "0.0.0.0:8080", "--log-level", "debug", "app:app"]