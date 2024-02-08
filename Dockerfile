# Use an official Node runtime as a parent image
FROM node:latest

# Set working directory
WORKDIR /app

# Install Hardhat globally
RUN npm install -g hardhat

# Install any other dependencies required for your Go server
# Install Go
# Install any dependencies for your Go server

# Copy the Go server code into the container
COPY ./ /app/go_server

# Install dependencies and build your Go server (commands depend on your Go project structure)
WORKDIR /app/go_server
RUN go mod tidy && go build -o krnl_node .

# Expose necessary ports for Hardhat node and Go server
EXPOSE 8080

# For running both Hardhat node and Go server
CMD ["sh", "-c", "npx hardhat node & ./go_server/krnl_node"]
