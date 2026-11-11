FROM golang:1.22

ENV PATH="/go/bin:${PATH}"
ENV GOTOOLCHAIN=auto

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install github.com/a-h/templ/cmd/templ@v0.2.793

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

CMD ["sh", "-c", "sqlc generate && templ generate && go build -o /tmp/servidor . && /tmp/servidor"]
