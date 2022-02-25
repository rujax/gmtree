FROM golang:1.16-alpine As build

WORKDIR /src/
COPY . .
RUN go mod download
RUN go build -o ./gmtree ./gmtree.go

FROM alpine

COPY --from=build /src/gmtree /usr/local/bin

ENTRYPOINT ["gmtree"]