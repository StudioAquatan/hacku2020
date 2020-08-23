FROM golang:1.15 AS builder
ENV PROJECT_DIR /go/src/github.com/StudioAquatan/hacku2020
COPY ./ ${PROJECT_DIR}/
WORKDIR ${PROJECT_DIR}/
RUN go build -i -o ./bin/oinori

FROM alpine:3.12 AS prod

COPY --from=builder /go/src/github.com/StudioAquatan/hacku2020/bin/oinori /
ENTRYPOINT ["/oinori"]
CMD ["run", "-c", "./character_config.yaml", "-n", "4"]
