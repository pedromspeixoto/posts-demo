# Build
FROM --platform=linux/amd64 node:14.5.0 AS build

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . .

RUN npm run build

# Distribution
FROM --platform=linux/amd64 nginx:stable-alpine

RUN apk update && apk add bash

COPY --from=build /app/build /usr/share/nginx/html
COPY --from=build /app/.env /usr/share/nginx/html/default.env
COPY --from=build /app/nginx/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/scripts/read_env.sh /usr/share/nginx/html/read_env.sh

RUN chmod +x /usr/share/nginx/html/read_env.sh

WORKDIR /usr/share/nginx/html

EXPOSE 80
CMD ["/bin/sh", "-c", "/usr/share/nginx/html/read_env.sh && nginx -g \"daemon off;\""]