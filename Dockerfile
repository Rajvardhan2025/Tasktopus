# Multi-stage Dockerfile - Builds both frontend and backend in one image

# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ ./

ARG VITE_BACKEND_URL=/api
ENV VITE_BACKEND_URL=$VITE_BACKEND_URL

RUN npm run build

# Stage 2: Build Backend
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app/backend

RUN apk add --no-cache git

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 3: Final Runtime Image
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata nginx

WORKDIR /app

# Copy backend binary
COPY --from=backend-builder /app/backend/main ./backend

# Copy backend scripts
COPY backend/scripts/ ./scripts/

# Copy frontend build
COPY --from=frontend-builder /app/frontend/dist ./frontend

# Copy nginx config for serving frontend and proxying API
COPY nginx.conf /etc/nginx/http.d/default.conf

# Create startup script
RUN echo '#!/bin/sh' > /app/start.sh && \
    echo 'nginx' >> /app/start.sh && \
    echo 'exec ./backend' >> /app/start.sh && \
    chmod +x /app/start.sh

EXPOSE 8080

CMD ["/app/start.sh"]
