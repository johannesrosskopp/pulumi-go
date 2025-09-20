#!/bin/bash

# Make pulumi run locally, better not commit this to source control even though it might be safe
mkdir -p .state
pulumi login file://.state/

pulumi config set azure-native:clientId '709fa6a3-e047-47b2-bca8-7268ce9cc4c1'
pulumi config set azure-native:tenantId 'e467b6d8-cf62-4e59-9a87-758ef858aeb6'
pulumi config set azure-native:subscriptionId 'cedbad7b-9e41-4314-b928-8c44b741a871'
pulumi config set azure-native:location 'westeurope'
pulumi config set azure-native:resourceGroupName 'pulumi-tech'
