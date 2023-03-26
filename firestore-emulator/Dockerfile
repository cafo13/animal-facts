FROM node:16.16.0-alpine3.15

RUN apk add --no-cache openjdk11-jre
RUN apk add --no-cache bash

RUN npm install -g firebase-tools@11.25.1

COPY start.sh .
COPY firebase.json .
COPY .firebaserc .

ENTRYPOINT ["./start.sh"]
