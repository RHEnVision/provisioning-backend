FROM registry.access.redhat.com/ubi8/go-toolset:1.18 as build
USER 0
RUN mkdir /build
WORKDIR /build
COPY . .
RUN make prep build strip GO=go

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
COPY --from=build /build/pbapi /pbapi
COPY --from=build /build/pbworker /pbworker
COPY --from=build /build/pbstatuser /pbstatuser
COPY --from=build /build/pbmigrate /pbmigrate
USER 1001
CMD ["/pbapi"]
