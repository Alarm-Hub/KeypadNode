# Build Stage
FROM lacion/alpine-golang-buildimage:1.13 AS build-stage

LABEL app="build-DoorManager"
LABEL REPO="https://github.com/Phill93/DoorManager"

ENV PROJPATH=/go/src/github.com/Phill93/DoorManager

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/Phill93/DoorManager
WORKDIR /go/src/github.com/Phill93/DoorManager

RUN make build-alpine

# Final Stage
FROM lacion/alpine-base-image:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/Phill93/DoorManager"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/DoorManager/bin

WORKDIR /opt/DoorManager/bin

COPY --from=build-stage /go/src/github.com/Phill93/DoorManager/bin/DoorManager /opt/DoorManager/bin/
RUN chmod +x /opt/DoorManager/bin/DoorManager

# Create appuser
RUN adduser -D -g '' DoorManager
USER DoorManager

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/DoorManager/bin/DoorManager"]
