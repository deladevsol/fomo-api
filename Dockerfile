FROM golang:1.26-bookworm AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /fomo-api ./cmd/fomo-api

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /fomo-api /fomo-api
EXPOSE 8787
ENTRYPOINT ["/fomo-api"]
CMD ["serve", "--listen", "0.0.0.0:8787"]
