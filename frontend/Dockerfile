FROM node:16.16.0 as builder

WORKDIR /frontend
COPY . .
RUN yarn install
RUN yarn build

FROM nginx:1.23.3

COPY --from=builder /frontend/dist/animal-facts /usr/share/nginx/html
