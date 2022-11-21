# We cannot use the official golang images from docker.io because it's blocked
# on our CI. Instead, we use image maintained by Quay team which is updated every
# few hours automatically from theird GitHub. For more info:
#
#  https://github.com/quay/claircore/actions/workflows/golang-image.yml
#  https://github.com/quay/claircore/blob/main/.github/workflows/golang-image.yml

FROM quay.io/projectquay/golang:1.18 as build
USER 0
RUN mkdir /build
WORKDIR /build
COPY . .
RUN make prep build strip

FROM registry.access.redhat.com/ubi8/ubi-minimal
COPY --from=build /build/pbapi /pbapi
COPY --from=build /build/pbworker /pbworker
COPY --from=build /build/pbstatuser /pbstatuser
COPY --from=build /build/pbmigrate /pbmigrate
USER 1001
CMD ["/pbapi"]
