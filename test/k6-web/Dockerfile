FROM grafana/k6:0.49.0 AS k6

FROM node:20-alpine

WORKDIR /app

COPY --from=k6 /usr/bin/k6 /usr/bin/k6
COPY package.json ./
COPY server.js ./
COPY public ./public

ENV NODE_ENV=production
ENV PORT=8080
ENV K6_NO_USAGE_REPORT=true

EXPOSE 8080

CMD ["node", "server.js"]
