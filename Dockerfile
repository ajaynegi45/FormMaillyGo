# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

# Disable CGO, target Linux
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /src

# Cache modules
COPY go.mod go.sum ./
RUN go mod download

# Copy application sources
COPY . .

# Build the Lambda bootstrap binary
RUN go build -o bootstrap aws_lambda.go

# Stage 2: Package for AWS Lambda custom runtime
FROM public.ecr.aws/lambda/provided.al2023

# Set working directory inside the Lambda container
WORKDIR /var/task

# Copy the compiled bootstrap and static assets
COPY --from=builder /src/bootstrap   ./bootstrap
COPY --from=builder /src/public      ./public

# Ensure the binary is executable
RUN chmod +x ./bootstrap

# Command Lambda will invoke
CMD ["bootstrap"]
