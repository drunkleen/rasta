#!/bin/bash
clear

set -e  # Exit immediately if a command exits with a non-zero status


BASEDIR=$(pwd)
# Set environment variables if needed
#export APP_ENV=${APP_ENV:-development}
#export PORT=${PORT:-8080}
#export DATABASE_URL=${DATABASE_URL:-"postgres://localhost:5432/rasta?sslmode=disable"}

LOGFILE="$BASEDIR/logs/startup.log"
mkdir -p "$(dirname "$LOGFILE")"  # Ensure the logs directory exists

echo "Starting Rasta application..." | tee -a "$LOGFILE"

# Check if the -debug argument is passed
if [[ "$*" == *"-debug"* ]]; then
    export RASTA_DEV_MODE=true
    echo "Debug mode enabled. RASTA_DEBUG_MODE is set to true." | tee -a "$LOGFILE"
fi

# Function to install Go tools if not present
install_tool() {
    TOOL_NAME=$1
    INSTALL_CMD=$2
    if ! command -v "$TOOL_NAME" > /dev/null; then
        echo "$TOOL_NAME not found. Installing..." | tee -a "$LOGFILE"
        eval "$INSTALL_CMD" 2>> "$LOGFILE"
        if [ $? -ne 0 ]; then
            echo "Failed to install $TOOL_NAME. Please install it manually." | tee -a "$LOGFILE"
            exit 1
        fi
    else
        echo "$TOOL_NAME is already installed." | tee -a "$LOGFILE"
    fi
}

# Install migrate if not installed
#install_tool "migrate" "go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"

# Run database migrations
#echo "Running database migrations..." | tee -a "$LOGFILE"
#migrate -path pkg/database/migrations -database "$DATABASE_URL" up 2>> "$LOGFILE"

# Install swag if not installed
install_tool "swag" "go install github.com/swaggo/swag/cmd/swag@latest"

# Generating Docs
echo "Generating documentation..." | tee -a "$LOGFILE"
swag init -g ./cmd/rasta/main.go -o ./docs/swagger 2>> "$LOGFILE"

# Start the application
echo "Starting the application on port $PORT..." | tee -a "$LOGFILE"
go run cmd/rasta/main.go 2>> "$LOGFILE"
