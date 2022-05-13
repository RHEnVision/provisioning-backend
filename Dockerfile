FROM registry.access.redhat.com/ubi8/go-toolset:latest as build
USER 0
RUN mkdir /build
WORKDIR /build
COPY . .
RUN make prep build strip

FROM registry.access.redhat.com/ubi8/ubi:latest
COPY --from=build /build/pbapi /pbapi
COPY --from=build /build/pbmigrate /pbmigrate
USER 1001
EXPOSE 8000
EXPOSE 5000
ENTRYPOINT ["/pbapi"]
