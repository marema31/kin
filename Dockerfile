FROM golang:1.14 AS build_img
ENV APP_DIR=/app
RUN mkdir -p $APP_DIR
WORKDIR $APP_DIR

COPY . .
RUN make static

ENTRYPOINT /usr/bin/kin

FROM scratch

COPY --from=build_img /app/bin/kin /usr/bin/kin

ENTRYPOINT ["/usr/bin/kin"]

CMD ["-d"]