
# Pulumi Azure Go Quickstart

This project demonstrates a complete serverless web application on Azure using Pulumi and Go. It features:

- **Azure Container App** - API backend running a containerized Go application
- **Azure Storage Static Website** - Frontend served from blob storage
- **Docker Image** - Built and pushed to Docker Hub+
- **Custom Command in pulumi** - To push the Api url into the static page before deployment

## Architecture

The final application consists of:
1. A Go API that returns the current time
2. A static HTML frontend that fetches and displays the time

## Prerequisites

- Docker to startup the devcontainer

## Setup

### 1. Configure Local State Backend

Set up Pulumi to use local file-based state management instead of Pulumi Cloud:

```bash
mkdir -p .state
pulumi login file://.state/
```

### 2. Select or Create Stack

```bash
pulumi stack select dev
# If 'dev' stack doesn't exist, create it.
# Enter your stack passphrase when prompted
```

### 3. Configure You name so your resources are unique

Set up Docker Hub credentials for image pushing:

```bash
# Set your name for resource naming
pulumi config set tech:myname <your-name>
```

### 4. Fix Docker Credential Helper (Dev Container Only)

⚠️ **Warning**: Only run this in isolated environments like dev containers! This will modify your Docker configuration.

```bash
# Create backup of existing Docker config (if it exists)
if [ -f ~/.docker/config.json ]; then
    cp ~/.docker/config.json ~/.docker/config.json.backup.$(date +%Y%m%d_%H%M%S)
    echo "✅ Backed up existing Docker config"
fi

# Create directory and empty config to disable credential helpers
mkdir -p ~/.docker
echo '{}' > ~/.docker/config.json
echo "✅ Docker credential helper disabled"
```

### 5. Import Shared Resource Group

Import the existing Azure resource group that your service principal manages:

```bash
pulumi import azure-native:resources:ResourceGroup my-resource-group \
  /subscriptions/$(pulumi config get azure-native:subscriptionId)/resourcegroups/$(pulumi config get azure-native:resourceGroupName)
```

## Operations

### Preview Changes

```bash
pulumi preview
```

### Deploy Infrastructure

```bash
pulumi up
```

### View Outputs

After deployment, Pulumi will output:
- `apiUrl` - Your API endpoint URL
- `primaryWebEndpoint` - Your static website URL

### Destroy Resources

To destroy all resources except the protected resource group:

```bash
pulumi destroy --exclude-protected
```

## Development

### Project Structure

```
├── main.go              # Pulumi infrastructure code
├── app/                 # API application
│   ├── Dockerfile       # Container definition
│   ├── main.go         # Go API server
│   └── go.mod          # API dependencies
├── www/                 # Frontend assets
│   ├── index.html      # Static HTML page
│   └── favicon.ico     # Website favicon (optional)
└── README.md           # This file
```
