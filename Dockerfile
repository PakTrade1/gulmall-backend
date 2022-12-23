From golang:1.19.3
ENV AZURE_TENANT_ID=8a58d6a5-8840-4b2b-aa1c-a4b9b23db880
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download 
COPY . .
RUN go build -o ./out/dist .
CMD ./out/dist