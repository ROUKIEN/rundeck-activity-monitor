# let's start by building the javascript app
FROM node:16-alpine AS frontend-builder

WORKDIR /app

COPY frontend/package.json frontend/package-lock.json /app/

RUN npm ci

COPY frontend/ /app

RUN npm run build

# golang builder, retrieving the bundled app from the previous step
FROM golang:1.19 AS backend-builder

WORKDIR /app

COPY backend/go.mod backend/go.sum /app/

RUN go mod download

COPY backend /app

COPY --from=frontend-builder /app/dist/ /app/cmd/static/

RUN ls -lah /app/cmd/static
RUN cat /app/cmd/static/index.html

RUN go build ram.go

# final stage
FROM gcr.io/distroless/base-debian11

COPY --from=backend-builder /app/ram /app/ram

EXPOSE 4000

USER nonroot:nonroot

ENTRYPOINT ["/app/ram"]

CMD ["help"]
