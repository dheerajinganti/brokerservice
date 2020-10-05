FROM golang:1.14 as build
ARG USER="broker"
ADD . /pg-broker
WORKDIR /pg-broker
RUN mkdir binaries
RUN GO111MODULE=on go mod vendor
RUN CGO_ENABLED=0 go build -o "binaries/pg-broker" .

FROM golang:1.14
COPY --from=build /pg-broker/binaries/* /bin/
RUN mkdir /bin/templates
COPY --from=build /pg-broker/templates/* /bin/templates/
COPY --from=build /pg-broker/config.yml /bin/
WORKDIR /bin
ENTRYPOINT ["pg-broker"]