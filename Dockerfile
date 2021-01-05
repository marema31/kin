FROM golang:1.14 AS build_img
ENV GOOS=linux
ENV GOARCH=amd64
ENV APP_DIR=/app
RUN mkdir -p $APP_DIR
WORKDIR $APP_DIR

COPY . .
RUN make static

ENTRYPOINT /app/bin/kin

FROM alpine

COPY --from=build_img /app/bin/kin /usr/bin/kin
COPY --from=build_img /app/site /root/.kin-root


ENTRYPOINT ["/usr/bin/kin"]

CMD ["-d"]