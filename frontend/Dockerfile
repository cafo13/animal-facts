FROM node:16.16.0 as builder

WORKDIR /frontend
COPY . .
RUN npm install
RUN npm run build

FROM nginx:1.23.3

COPY --from=builder /frontend/build /usr/share/nginx/html
