# DOCKERFILE ------------------------------------------------------------------------
# -----------------------------------------------------------------------------------

    
# STAGE 01: Build the Go Binary *****************************************************

# 1. Tell Docker to use Golang as the base for building the app ---------------------
FROM golang:1.24.2 AS builder

# 2. Set the working directory inside the container to /app -------------------------
# All the next commands will run inside this folder.
WORKDIR /app

# 3. Copy go mod and sum files ------------------------------------------------------
# The files get copied from the computer into the container.
COPY go.mod go.sum ./
# The dependencies listed in go.mod get downloaded in the container so that
# the app can be compiled.
RUN go mod download

# 4. Copy the source code -----------------------------------------------------------
# All the source code gets copied from the computer into the container's /app folder.
COPY . .

# 5. Copy Docker-specific .env file
COPY docker/.env ./docker.env
ENV ENV_PATH=./docker.env

# 6. Build the Go binary ------------------------------------------------------------
# Build disabling C code (more portable) in Linux (Docker runs in Linux) and build
# a binary named server in the folder ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api



# STAGE 02: Create minimal image ***************************************************

# 1. Start new light image using Alpine Linux --------------------------------------
FROM alpine:latest

# 2. Set the working directory inside the new image --------------------------------
WORKDIR /root/

# 3. Copy binary from builder stage into new image ---------------------------------
COPY --from=builder /app/server .

# 4. Copy binary .env into new image -----------------------------------------------
COPY --from=builder /app/server .
COPY --from=builder /app/docker.env ./docker.env

# 5. Set ENV_PATH env variable for docker's .env file ------------------------------
ENV ENV_PATH=./docker.env

# 6. Set Port for app to listen on ------------------------------------------------- 
EXPOSE 8080

# 7. Start the API -----------------------------------------------------------------
# Run the Golang API by executing the server binary
CMD ["./server"]

