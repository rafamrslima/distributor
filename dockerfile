FROM golang:1.25 AS builder
WORKDIR /src
COPY . .
RUN go build -o /distributor ./cmd/distributor

FROM gcr.io/distroless/base-debian12
COPY --from=builder /distributor /distributor
ENTRYPOINT ["/distributor"]