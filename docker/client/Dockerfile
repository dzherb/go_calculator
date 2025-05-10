FROM node:alpine AS builder

WORKDIR /app

COPY ./client/package.json .

RUN yarn install

COPY ./client .

ARG BACKEND_BASE_URL

ENV VITE_BACKEND_BASE_URL=$BACKEND_BASE_URL

RUN yarn build

FROM nginx:alpine AS prod

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80