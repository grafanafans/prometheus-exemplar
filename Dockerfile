FROM golang:alpine AS builder
WORKDIR /build
ADD go.mod .
COPY . .
RUN GOPROXY=off go build -o exemplar -mod=vendor cmd/main.go

FROM alpine
WORKDIR /app 
COPY --from=builder /build/exemplar /app/exemplar
CMD ["./exemplar"]