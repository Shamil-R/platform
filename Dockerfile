FROM golang:1.11.1 AS builder

RUN go get -d github.com/golang/dep/cmd/dep && \
    go install github.com/golang/dep/cmd/dep

WORKDIR $GOPATH/src/gitlab/nefco/platform

RUN apt-get update && apt-get install -y \
		libkrb5-dev \
		#&& aptitude install -y krb5-user \
		#&& apt-get install -y aptitude \
		#&& aptitude install -y openafs-krb5 \
    && rm -rf /var/lib/apt/lists/*

COPY Gopkg.toml Gopkg.lock ./



RUN dep ensure --vendor-only -v

COPY . ./



RUN go run ./platform/main.go codegen --clean

RUN go build -o /platform ./platform/main.go




FROM golang:1.11.1

COPY --from=builder /platform .

COPY .platform.yml .
COPY test.keytab .
RUN export KRB5_CLIENT_KTNAME=./test.keytab

CMD [ "./platform", "run" ]
