FROM mcr.microsoft.com/devcontainers/go:1-1.23-bookworm

# Install Azure CLI
RUN curl -sL https://aka.ms/InstallAzureCLIDeb | bash

# Install pulumi
RUN curl -fsSL https://get.pulumi.com | sh && \
    mkdir /usr/bin/.pulumi && \
    mv /root/.pulumi/* /usr/bin/.pulumi/ 
ENV PATH "$PATH:/usr/bin/.pulumi/bin/"