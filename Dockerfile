ARG GO_VERSION=1.23

FROM golang:${GO_VERSION}-alpine as builder

ENV USER=nonroot
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /app

# Copy application code
COPY . .

RUN rm go.work

ARG SERVICE_NAME
ARG LIBRARIES
RUN go work init ./cmd/${SERVICE_NAME} ${LIBRARIES}

RUN go mod download

RUN go mod verify


ARG PORT
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o main \
    -ldflags="-X main.port=${PORT}" ./cmd/${SERVICE_NAME}

# Deploy stage
FROM alpine:latest

WORKDIR /root/

# Copy binary
COPY --from=builder /app/main .

COPY --from=builder /app/cmd/ui/templates ./templates

RUN ls 

ARG SERVICE_NAME
RUN echo "SERVICE_NAME is: ${SERVICE_NAME}" && \
    if [ "${SERVICE_NAME}" != "ui" ]; then \
        echo "Removing templates as SERVICE_NAME is not ui"; \
        rm -rf ./templates; \
    else \
        echo "Keeping templates as SERVICE_NAME is ui"; \
    fi

RUN ls

ARG PORT
EXPOSE ${PORT}

CMD ["./main"]
