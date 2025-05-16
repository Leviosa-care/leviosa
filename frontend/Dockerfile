FROM node:20.11.1-alpine AS build
WORKDIR /app

ARG SESSION_COOKIE_NAME="leviosa_session"
ARG CLIENT_IP_HEADER="x-client-ip"
ARG BACKEND_PORT="http://backend:5000"

ENV SESSION_COOKIE_NAME=${SESSION_COOKIE_NAME}
ENV CLIENT_IP_HEADER=${CLIENT_IP_HEADER}
ENV API_URL="http://backend:${BACKEND_PORT}"

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

RUN corepack enable
COPY ./package.json ./pnpm-lock.yaml ./
RUN pnpm install

COPY . .
RUN pnpm run build

FROM node:20.11.1-alpine

WORKDIR /app
RUN adduser -D runner

COPY --from=build /app/build build/
COPY --from=build /app/node_modules node_modules/ 
COPY package.json pnpm-lock.yaml ./

RUN chown -R runner:runner /app/build

EXPOSE 3000 3001

CMD ["node", "build"]
