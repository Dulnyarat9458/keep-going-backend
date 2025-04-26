FROM golang:1.24.2

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main main.go

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
