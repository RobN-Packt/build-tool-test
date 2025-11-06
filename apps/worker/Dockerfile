# syntax=docker/dockerfile:1.6

FROM node:20-alpine AS deps
WORKDIR /app
COPY package.json pnpm-lock.yaml* ./
RUN corepack enable && pnpm install --frozen-lockfile

FROM node:20-alpine AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN corepack enable && pnpm build

FROM public.ecr.aws/lambda/nodejs:20
WORKDIR /var/task
COPY --from=builder /app/dist ./dist
COPY package.json ./package.json
CMD ["dist/src/handler.handler"]
