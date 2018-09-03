FROM golang:1.10 AS builder

RUN mkdir -p /go/bin/match-data-ingest
WORKDIR /go/bin/match-data-ingest
COPY . .
RUN go get -u github.com/go-sql-driver/mysql
RUN go get -u github.com/tidwall/gjson
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

FROM harborfront/base
COPY --from=builder /app ./
EXPOSE 3306:3306
# NEEDS CHANGING DEPENDING ON HOST OS
# COPY ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["./app"]
