FROM golang:1-13-alpine

WORKDIR /app
COPY . .
RUN go install

EXPOSE 80
ENTRYPOINT ["flights"]