FROM golang:1.23-bookworm

WORKDIR /app/backend/auth

COPY backend/auth/go.mod backend/auth/go.sum ./
RUN go mod download

COPY backend/auth ./

EXPOSE 8080 50051

CMD ["go", "run", "./cmd/auth/main.go"]
