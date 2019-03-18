FROM golang:alpine3.8 as builder

RUN apk --update upgrade
RUN apk --no-cache --no-progress add make git gcc musl-dev

WORKDIR /build
COPY . .
RUN go build .

FROM node:10-alpine
RUN apk update && apk add --no-cache --virtual ca-certificates
COPY --from=builder /build/jest-action /usr/bin/jest-action

LABEL version="1.0.0"
LABEL repository="https://github.com/rkusa/jest-action"
LABEL homepage="https://github.com/rkusa/jest-action"
LABEL maintainer="Markus Ast <m@rkusa.st>"

LABEL com.github.actions.name="Annotated Jest"
LABEL com.github.actions.description="Execute jest tests and test failure annotations"
LABEL com.github.actions.icon="check"
LABEL com.github.actions.color="green"

ENV JEST_CMD ./node_modules/.bin/jest
COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD [""]