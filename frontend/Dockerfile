FROM node:20.11.1-alpine AS base
WORKDIR /usr/src/app

ARG SESSION_COOKIE_NAME="leviosa_session"
ARG CLIENT_IP_HEADER="x-client-ip"
ARG BACKEND_PORT=5000

ENV SESSION_COOKIE_NAME=${SESSION_COOKIE_NAME}
ENV CLIENT_IP_HEADER=${CLIENT_IP_HEADER}
ENV API_URL="http://backend:${BACKEND_PORT}"

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

COPY package.json pnpm-lock.yaml ./
COPY .npmrc* ./

FROM base AS prod-deps
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod --frozen-lockfile

FROM base AS build
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

COPY . .
RUN pnpm run build

FROM gcr.io/distroless/nodejs20-debian12:nonroot AS runtime
WORKDIR /usr/src/app
COPY --from=prod-deps /app/node_modules/ ./node_modules/
COPY --from=build /app/build/ ./build/
COPY --from=build /app/package.json ./

EXPOSE 3000 3001

CMD ["build/index.js"]
