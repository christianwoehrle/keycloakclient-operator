FROM golang:1.18.2-buster  AS build-env

COPY . /src/

RUN cd /src && \
    make code/compile && \
    echo "Build SHA1: $(git rev-parse HEAD)" && \
    echo "$(git rev-parse HEAD)" > /src/BUILD_INFO

#yum install make

FROM ubuntu:latest AS user-distro
RUN useradd -u 10001 scratchuser



# final stage
FROM scratch

COPY --from=user-distro /etc/passwd /etc/passwd
USER scratchuser
##LABELS

COPY --from=build-env /src/BUILD_INFO /src/BUILD_INFO
COPY --from=build-env /src/tmp/_output/bin/keycloakclient-operator /

ENTRYPOINT ["/keycloakclient-operator"]
