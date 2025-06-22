FROM oven/bun:1.2-alpine AS builder

WORKDIR /app

# Copy package files
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

# Copy source code
COPY . .

# Build the app
RUN bun run build

# Production stage
FROM nginx:1.28-alpine

# Copy built app to nginx
COPY --from=builder /app/build /usr/share/nginx/html

# Copy nginx config
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
