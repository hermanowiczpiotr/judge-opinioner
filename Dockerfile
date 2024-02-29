FROM golang:alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o judge-opinioner
COPY .env ./

CMD ["/app/judge-opinioner"]