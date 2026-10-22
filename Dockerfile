FROM --platform=$BUILDPLATFORM golang:1.27.3-alpine AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags='-s -w' -o /out/acg-api .

FROM alpine:3.22

WORKDIR /app

COPY --from=build /out/acg-api /usr/local/bin/acg-api

EXPOSE 8080

ENTRYPOINT ["acg-api"]
