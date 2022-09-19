FROM golang:1.18.5-alpine3.15 AS builder

COPY . /Samplicity/cmd/samplicityserver
WORKDIR /Samplicity/cmd/samplicityserver

RUN apk  update
RUN apk add ffmpeg

RUN go build -o ./bin/main DriveApi/cmd

FROM alpine:3.15
WORKDIR /root/

COPY --from=0 /Samplicity/cmd/samplicityserver/bin/main .
COPY --from=0 /Samplicity/cmd/samplicityserver/*.env .

CMD ["./main"]