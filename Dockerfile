# Use a minimal base image to reduce the size of the final image
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the pre-built binary from the local directory to the container
COPY finsight .

# Entry point
ENTRYPOINT ["./finsight"]