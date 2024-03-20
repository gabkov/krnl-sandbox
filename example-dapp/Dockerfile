FROM node:18.16-alpine

RUN apk add screen

WORKDIR /app
COPY package.json hardhat.config.ts tsconfig.json /app/

RUN npm install

CMD [ "npx", "hardhat", "node" ]