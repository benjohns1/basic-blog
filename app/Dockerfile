FROM node:12-alpine AS builder
WORKDIR /usr/src/app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:1.13.12-alpine
COPY --from=builder /usr/src/app/dist/app /usr/share/nginx/html
COPY ./nginx.conf /etc/nginx/conf.d/default.conf
