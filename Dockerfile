# Start by building the application.
FROM golang:1.21 as build

WORKDIR /go/src/app
RUN mkdir bin
COPY cmd/ cmd
COPY pkg/ pkg
COPY Makefile go.mod go.sum ./

RUN go mod download
RUN CGO_ENABLED=0 make bin/authentication-server

# Now copy it into our base image.
FROM gcr.io/distroless/static-debian12
COPY --from=build /go/src/app/bin/authentication-server /
CMD ["/authentication-server", "--otp-secret-file" , "/mnt/assets/secret", "--assets-dir", "/mnt/assets/static/"]
