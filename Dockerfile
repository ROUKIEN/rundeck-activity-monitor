FROM node:16-alpine AS frontend-builder

WORKDIR /app

COPY frontend/package.json frontend/package-lock.json /app/

RUN npm ci

COPY frontend/ /app

RUN npm run build

FROM golang:1.18 AS backend-builder

WORKDIR /app

COPY backend/go.mod backend/go.sum /app/

RUN go mod download

COPY backend /app

RUN go build ram.go

FROM gcr.io/distroless/base-debian11
# FROM alpine

COPY --from=backend-builder /app/ram /app/ram

EXPOSE 4000

USER nonroot:nonroot

ENTRYPOINT ["/app/ram"]

CMD ["help"]
