FROM node:12-alpine AS builder
WORKDIR /usr/src/app
COPY ./app/package*.json ./
RUN npm install
COPY ./app .
RUN npm run build

FROM nginx:1.13.12-alpine
COPY --from=builder /usr/src/app/dist/app /usr/share/nginx/html
COPY ./app/nginx.conf /etc/nginx/conf.d/default.conf
