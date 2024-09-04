FROM registry.access.redhat.com/ubi9/go-toolset:1.21.11-9 as build
USER 0
RUN mkdir /build
WORKDIR /build
COPY . .
RUN make prep build strip GO=go

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest
COPY --from=build /build/pbackend /pbackend
USER 1001
CMD ["/pbackend", "api"]
