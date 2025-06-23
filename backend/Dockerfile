FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o rgxr-backend .

FROM texlive/texlive:latest-small

RUN tlmgr install preview

RUN apt-get update && apt-get install -y dot2tex pdf2svg

WORKDIR /app
COPY --from=builder /app/rgxr-backend .

# Create data directories
RUN mkdir -p /data/images /data/tex

EXPOSE 8080
CMD ["./rgxr-backend"]
