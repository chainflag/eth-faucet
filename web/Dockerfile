FROM node:lts-alpine as build-stage

WORKDIR /app

COPY rollup.config.js ./
COPY package*.json ./
RUN npm install

COPY ./src ./src
COPY ./public ./public
RUN npm run build

FROM nginx:stable-alpine as production-stage
COPY --from=build-stage /app/public /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
