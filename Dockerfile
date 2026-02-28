FROM golang:1 AS build

WORKDIR /app

COPY . .

RUN go build -o bob ./cmd/bob

FROM scratch

WORKDIR /cli

COPY --from=build /app/bob bob

ENTRYPOINT ["./bob"]
