FROM golang:1.15.8-alpine as builder
# RUN apk add --update git
ADD . /webhook/
WORKDIR /webhook/
RUN go build -o ./bin/postee main.go


FROM alpine
RUN apk update && apk add wget ca-certificates
EXPOSE 8082
EXPOSE 8445
RUN mkdir /server
RUN mkdir /server/database
RUN mkdir /config
COPY --from=builder /webhook/bin /server/
COPY --from=builder /webhook/rego-templates /server/rego-templates
COPY --from=builder /webhook/cfg.yaml /config/
WORKDIR /server
RUN chmod +x postee
RUN addgroup -g 1099 webhook
RUN adduser -D -g '' -G webhook -u 1099 webhook
RUN chown -R webhook:webhook /server
RUN chown -R webhook:webhook /config
USER webhook
ENTRYPOINT ["/server/postee"]
