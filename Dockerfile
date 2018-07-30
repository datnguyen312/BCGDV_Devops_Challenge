FROM golang

ARG GIT_COMMIT=unspecified
LABEL git_commit=$GIT_COMMIT

WORKDIR /go/src/app
COPY ./api/ .

RUN go get -d -v -t ./...
RUN go install -v ./...

CMD ["app"]

EXPOSE 9090
