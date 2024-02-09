FROM node:18.16-alpine

WORKDIR /app
COPY start.sh start.sh

RUN npm install -g hardhat

COPY --from=golang:1.20-alpine /usr/local/go/ /usr/local/go/
 
ENV PATH="/usr/local/go/bin:${PATH}"

COPY ./ /app

WORKDIR /app/krnl
RUN go mod tidy && go build -o krnl_node .

RUN apk add screen

WORKDIR /app

RUN chmod +x start.sh
ENTRYPOINT [ "sh", "-c", "./start.sh" ]
