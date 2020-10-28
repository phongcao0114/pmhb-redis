FROM paymenthub-789-docker.artifactory.kasikornbank.com:8443/common/golang:1.13.5-alpine3.11 as builder
LABEL vendor="kbtg" project="pmhb-redis"

COPY . /go/src/app

WORKDIR /go/src/app  
ENV PATH=/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin  
ENV GOPATH=/go  
ENV GOLANG_VERSION=1.13.5

# Add zoneinfo.zip to GOROOT/lib so time.LoadLocation can work inside alpine image.
RUN mkdir -p /usr/lib/go-1.13/lib/time
COPY deployment/docker/assets/zoneinfo.zip /usr/lib/go-1.13/lib/time/zoneinfo.zip

# Add troubleshooting tools
COPY deployment/docker/assets/httpstat /usr/bin/httpstat
RUN chmod +x /usr/bin/httpstat

# Build
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o /pmhb-redis.bin .

########## Open new clean box image ##########
FROM paymenthub-789-docker.artifactory.kasikornbank.com:8443/common/phub-alpine:1.0.0
COPY --from=builder /pmhb-redis.bin ./
RUN mkdir configs
COPY configs ./configs
RUN ls -l
RUN chmod +x pmhb-redis.bin


EXPOSE 8081
CMD ["./pmhb-redis.bin", "-port", "8081", "-config", "./configs/"]