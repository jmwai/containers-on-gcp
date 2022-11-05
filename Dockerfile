FROM golang:1.16-alpine AS build_base

WORKDIR /tmp/app

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build -o  ./out/places /tmp/app/server/

FROM alpine:3.16.2

WORKDIR /root/

COPY --from=build_base /tmp/app/out/places .
# COPY --from=build_base /tmp/app/.env .
COPY --from=build_base /tmp/app/ .

EXPOSE 8080

CMD [ "./places" ]