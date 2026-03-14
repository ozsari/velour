# Stage 1: Build frontend
FROM node:22-alpine AS frontend
WORKDIR /app/web/frontend
COPY web/frontend/package*.json ./
RUN npm ci --legacy-peer-deps
COPY web/frontend/ ./
RUN npm run build

# Stage 2: Build Go backend
FROM golang:1.26-alpine AS backend
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/web/dist ./web/dist
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o /velour ./cmd/velour

# Stage 3: Final image
FROM alpine:3.20
RUN apk add --no-cache ca-certificates docker-cli sqlite
WORKDIR /app

COPY --from=backend /velour /app/velour
COPY --from=frontend /app/web/dist /app/web/dist

RUN mkdir -p /opt/velour

ENV VELOUR_HOST=0.0.0.0
ENV VELOUR_PORT=8585
ENV VELOUR_DATA_DIR=/opt/velour

EXPOSE 8585

VOLUME ["/opt/velour"]

ENTRYPOINT ["/app/velour"]
