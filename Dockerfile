FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /app
COPY main.go . 
COPY web web
COPY filedata.json .
 
RUN go mod init bdasrv
RUN go mod tidy
RUN GOOS=linux go build -ldflags="-s -w" -o ./main ./main.go

FROM alpine:3.13
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=build /app /app
RUN ls -laR .

EXPOSE 80
ENTRYPOINT /app/main --port 80