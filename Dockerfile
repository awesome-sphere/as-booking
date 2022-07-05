FROM golang:1.18-alpine3.16 as build-stage

COPY ./ /go/src/as-booking
WORKDIR /go/src/as-booking

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -v -o as-booking

FROM alpine:latest as production-stage

RUN apk --no-cache add ca-certificates

COPY --from=build-stage /go/src/as-booking /as-booking
WORKDIR /as-booking

CMD ["./as-booking"]