FROM registry.access.redhat.com/ubi9/go-toolset:1.20 as build
USER 0
RUN mkdir /build
WORKDIR /build
COPY Makefile cmd/ internal/ pkg/ mk/ .
RUN make prep build strip GO=go

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest
COPY --from=build /build/pbackend /pbackend
USER 1001
CMD ["/pbackend", "api"]
