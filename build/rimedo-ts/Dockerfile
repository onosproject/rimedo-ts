# SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
# SPDX-FileCopyrightText: 2019-present Rimedo Labs
#
# SPDX-License-Identifier: Apache-2.0

FROM onosproject/golang-build:v1.0 as build

ENV GO111MODULE=on
ENV XAPPNAME=rimedo-ts

COPY Makefile go.mod go.sum /go/src/github.com/rimedo-labs/rimedo-ts/
COPY cmd/ /go/src/github.com/rimedo-labs/rimedo-ts/cmd/
COPY vendor/ /go/src/github.com/rimedo-labs/rimedo-ts/vendor/
COPY pkg/ /go/src/github.com/rimedo-labs/rimedo-ts/pkg/
COPY build/ /go/src/github.com/rimedo-labs/rimedo-ts/build/
# COPY policytypes /go/src/github.com/rimedo-labs/rimedo-ts/policytypes/
# COPY schemas /go/src/github.com/rimedo-labs/rimedo-ts/schemas/

RUN cd /go/src/github.com/rimedo-labs/rimedo-ts && GOFLAGS=-mod=vendor make build

FROM alpine:3.11
RUN apk add libc6-compat

USER nobody

COPY --from=build /go/src/github.com/rimedo-labs/rimedo-ts/build/_output/rimedo-ts /usr/local/bin/rimedo-ts
# COPY --from=build /go/src/github.com/rimedo-labs/rimedo-ts/schemas /data/schemas
# COPY --from=build /go/src/github.com/rimedo-labs/rimedo-ts/policytypes /data/policytypes

ENTRYPOINT ["rimedo-ts"]
