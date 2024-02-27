FROM node:18.16-alpine

RUN apk add screen

COPY --from=golang:1.20-alpine /usr/local/go/ /usr/local/go/
 
ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /app
COPY ./ /app

COPY ./docker.env /app/.env

RUN go mod tidy && go build -o krnl_node .

RUN chmod +x start.sh

WORKDIR /app/_hardhat
RUN npm install

WORKDIR /app

ENTRYPOINT [ "sh", "-c", "./start.sh" ]
