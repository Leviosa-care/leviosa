FROM node:20.11.1-alpine AS base
WORKDIR /app

ARG SESSION_COOKIE_NAME="leviosa_session"
ARG CLIENT_IP_HEADER="x-client-ip"
ARG BACKEND_PORT=5000
ARG API_URL="http://backend:${BACKEND_PORT}"

ENV SESSION_COOKIE_NAME=${SESSION_COOKIE_NAME}
ENV CLIENT_IP_HEADER=${CLIENT_IP_HEADER}
ENV API_URL=${API_URL}

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

COPY . .

FROM base AS prod-deps
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod --frozen-lockfile

FROM base AS build
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
RUN pnpm run build

FROM gcr.io/distroless/nodejs20-debian12:nonroot
WORKDIR /app
COPY --from=prod-deps /app/node_modules/ ./node_modules/
COPY --from=build /app/build/ ./build/
EXPOSE 3000 3001

CMD ["build/index.js"]
