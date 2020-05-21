FROM golang:1.14 AS build_img
ENV APP_DIR=/app
RUN mkdir -p $APP_DIR
WORKDIR $APP_DIR

COPY . .
RUN make

#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
#    go build -gcflags "all=-N -l" -o /kin

ENTRYPOINT /kin

FROM scratch

COPY --from=build_img /app/bin/kin /usr/bin/kin

ENTRYPOINT ["/usr/bin/kin" ]
