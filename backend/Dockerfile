FROM golang:1.24-alpine
RUN apk add --no-cache graphviz python3 py3-pip
RUN pip install dot2tex --break-system-packages
WORKDIR /app
COPY . .
RUN go build -o server .
EXPOSE 8080
CMD ["./server"]
