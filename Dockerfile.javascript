FROM node:18-alpine

COPY tests/ /tests
COPY javascript/ /javascript

WORKDIR /javascript

RUN npm ci

CMD ["npm", "run", "test"]
