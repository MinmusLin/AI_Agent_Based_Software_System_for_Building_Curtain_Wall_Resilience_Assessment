FROM node:22-alpine AS node-builder

WORKDIR /app

ARG VITE_API_BASE_URL
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM nginx:1.29-alpine AS nginx-runtime

RUN rm -f /etc/nginx/conf.d/default.conf \
	&& printf '%s\n' \
	'server {' \
	'  listen 80;' \
	'  server_name icw.core.web;' \
	'  root /usr/share/nginx/html;' \
	'  index index.html;' \
	'  location / {' \
	'    try_files $uri $uri/ /index.html;' \
	'  }' \
	'  location ~* \.(?:js|mjs|css|png|jpe?g|gif|webp|avif|svg|ico|bmp|ttf|otf|eot|woff2?|wasm|webmanifest|map|webm|mp4|ogg|mp3|wav)$ {' \
	'    try_files $uri =404;' \
	'    expires 30d;' \
	'    add_header Cache-Control "public, immutable";' \
	'  }' \
	'}' \
	> /etc/nginx/conf.d/default.conf

COPY --from=node-builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
