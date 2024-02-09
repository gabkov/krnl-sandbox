FROM golang:1.20-alpine

WORKDIR /app
COPY ./ /app

RUN go mod tidy && go build -o token_authority .

ENTRYPOINT [ "sh", "-c", "./token_authority" ]
