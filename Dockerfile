FROM golang AS builder

RUN go get -d github.com/golang/dep/cmd/dep && \
    go install github.com/golang/dep/cmd/dep

WORKDIR $GOPATH/src/gitlab/nefco/platform

COPY Gopkg.toml Gopkg.lock ./

RUN dep ensure --vendor-only -v

COPY . ./

RUN go run ./platform/main.go codegen --clean

RUN CGO_ENABLED=0 go build -o /platform ./platform/main.go

#RUN apt-get update && apt-get install -y \
		#libkrb5-dev \
		#&& aptitude install -y krb5-user \
		#&& aptitude install -y openafs-krb5 \
    #&& rm -rf /var/lib/apt/lists/*


FROM golang:alpine

COPY --from=builder /platform .

COPY .platform.yml .

CMD [ "./platform", "run" ]
