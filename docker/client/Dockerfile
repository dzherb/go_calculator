FROM node:alpine AS builder

WORKDIR /app

COPY ./client/package.json .

RUN yarn install

COPY ./client .

ARG BACKEND_BASE_URL

ENV VITE_BACKEND_BASE_URL=$BACKEND_BASE_URL

RUN yarn build

FROM alpine

COPY --from=builder /app/dist /tmp_static

RUN chmod +rwx /tmp_static/*

CMD cp -R /tmp_static/* /static