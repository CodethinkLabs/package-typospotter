FROM golang:1.16.3@sha256:f7d3519759ba6988a2b73b5874b17c5958ac7d0aa48a8b1d84d66ef25fa345f1 as build
WORKDIR /src
COPY . /src
RUN go build -o package-typospotter ./cmd/squatter-spotter/main.go


FROM gcr.io/distroless/base:nonroot@sha256:bc84925113289d139a9ef2f309f0dd7ac46ea7b786f172ba9084ffdb4cbd9490

COPY --from=build /src/package-typospotter /usr/local/bin/package-typospotter

ENTRYPOINT ["/usr/local/bin/package-typospotter"]