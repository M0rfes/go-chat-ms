ARG GO_VERSION=1.23

# Use Debian as the base image (glibc)
FROM golang:${GO_VERSION}-bullseye as builder

# Install librdkafka-dev (pre-built for glibc)
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    librdkafka-dev \
    build-essential

# ENV USER=nonroot
# ENV UID=10001
# RUN useradd \
#     --disabled-password \
#     --gecos "" \
#     --home "/na" \
#     --shell "/sbin/nologin" \
#     --no-create-home \
#     --uid "${UID}" \
#     "${USER}"

WORKDIR /app

# Copy application code
COPY . .

RUN rm -rf go.work go.work.sum

ARG SERVICE_NAME
ARG PORT
ARG LIBRARIES

# Initialize go.work inside the container (absolute paths are safer)
RUN go work init && \
    go work use ./cmd/${SERVICE_NAME} && \
    for lib in ${LIBRARIES}; do \
        go work use ${lib}; \
    done

# Download modules before building (this improves caching)
RUN go mod download

# Verify modules
RUN go mod verify

# Build with cross-compilation support for librdkafka and updated import path
RUN CGO_ENABLED=1 \
    go build -a -o main \
    -ldflags="-X main.port=${PORT}" \
    ./cmd/${SERVICE_NAME}

# Intermediate stage for copying templates (only if SERVICE_NAME is ui)
FROM debian:bullseye-slim as templates-stage

ARG SERVICE_NAME
WORKDIR /templates

# Only copy templates if SERVICE_NAME matches
COPY --from=builder /app/cmd/ui/templates /templates

# Deploy stage (final stage)
FROM debian:bullseye-slim

WORKDIR /root/

# Copy binary from the builder stage
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