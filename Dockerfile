From golang:1.19.3
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download 
COPY . .
RUN go build -o ./out/dist .
CMD ./out/dist