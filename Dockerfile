# Step #1
FROM golang:1.22 AS firststage
LABEL description="Go Monitor"
LABEL maintainer="Bagas <mbagas221@gmail.com>"
WORKDIR /build/
COPY . /build
ENV CGO_ENABLED=0
RUN go get
RUN go build -o go-monitor

# Step #2
FROM alpine:latest
WORKDIR /app/
RUN apk update
RUN apk upgrade
RUN apk add ca-certificates && update-ca-certificates
RUN apk add --no-cache tzdata gcompat
ENV TZ=Asia/Jakarta
COPY --from=firststage /build/go-monitor .
CMD ["./go-monitor"]