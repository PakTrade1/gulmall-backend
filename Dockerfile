From golang:1.16.6
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download 
COPY . . 
RUN go build -o ./out/dist .
CMD ./out/dist