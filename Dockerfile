FROM golang:1.19 as build

WORKDIR /go/src/app
COPY . .

RUN make BIN_DIR=/go/bin/app build

FROM gcr.io/distroless/static-debian11

COPY --from=build /go/bin/app /bin/credstore-csi-driver
ENTRYPOINT [ "/bin/credstore-csi-driver" ]
