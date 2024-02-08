FROM node:18

WORKDIR /app
COPY start.sh start.sh

RUN npm install -g hardhat

RUN apt-get update && apt-get install -y golang-go

COPY ./ /app

WORKDIR /app/krnl
RUN go mod tidy && go build -o krnl_node .


WORKDIR /app

RUN chmod +x start.sh
ENTRYPOINT [ "bash", "-c", "./start.sh" ]
