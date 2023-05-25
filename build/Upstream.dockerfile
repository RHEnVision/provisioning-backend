# We cannot use the official golang images from docker.io because it's blocked
# on our CI. Instead, we use image maintained by Quay team which is updated every
# few hours automatically from theird GitHub. For more info:
#
#  https://github.com/quay/claircore/actions/workflows/golang-image.yml
#  https://github.com/quay/claircore/blob/main/.github/workflows/golang-image.yml

FROM quay.io/projectquay/golang:1.19 as build
USER 0
RUN mkdir /build
WORKDIR /build
COPY . .
RUN make prep build strip GO=go

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest
COPY --from=build /build/pbackend /pbackend
USER 1001
CMD ["/pbackend", "api"]
