FROM golang:1.22

ENV PATH="/go/bin:${PATH}"
ENV GOTOOLCHAIN=auto

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

RUN curl -sSfL https://github.com/Orange-OpenSource/hurl/releases/download/8.0.1/hurl_8.0.1_amd64.deb -o /tmp/hurl.deb \
    && apt-get update \
    && apt-get install -y --no-install-recommends /tmp/hurl.deb \
    && rm -f /tmp/hurl.deb \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

CMD ["sh", "-c", "sqlc generate && go run ."]
