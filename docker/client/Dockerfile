FROM node:alpine AS builder

WORKDIR /app

COPY ./client/package.json .

RUN yarn install

COPY ./client .

ARG BACKEND_SERVER_HOST
ARG BACKEND_SERVER_PORT

ENV VITE_BACKEND_SERVER_HOST=$BACKEND_SERVER_HOST
ENV VITE_BACKEND_SERVER_PORT=$BACKEND_SERVER_PORT

RUN yarn build

FROM nginx:alpine AS prod

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80