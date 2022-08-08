FROM node:lts-slim

WORKDIR /usr/src/app
COPY package.json yarn.lock ./

RUN yarn install --pure-lockfile
COPY . .

RUN yarn build
RUN yarn autoclean --init
RUN yarn autoclean --force

CMD yarn start:prod