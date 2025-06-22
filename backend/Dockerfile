FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o rgxr-backend .

FROM alpine:latest

RUN apk add --no-cache \
    graphviz \
    pdf2svg \
    python3 \
    py3-pip

# Install dot2tex
RUN pip3 install dot2tex

WORKDIR /app
COPY --from=builder /app/rgxr-backend .

# Create data directories
RUN mkdir -p /data/images /data/tex

EXPOSE 8080
CMD ["./rgxr-backend"]
