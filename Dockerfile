FROM golang:1.23.4

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o notesApp .

EXPOSE 8080

CMD ["./notesApp"]