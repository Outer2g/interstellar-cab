FROM golang:1.17-alpine AS build

WORKDIR /interstellar-cab

COPY Makefile ./
COPY go.mod ./
COPY go.sum ./

RUN go mod download

ADD cmd ./cmd
ADD pkg ./pkg

RUN go build -o app cmd/main.go

FROM golang:1.17-alpine

WORKDIR /

COPY --from=build /interstellar-cab/app /inter-app

CMD /inter-app
