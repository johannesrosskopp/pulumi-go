
```bash
pulumi import azure-native:resources:ResourceGroup my-resource-group /subscriptions/$(pulumi config get azure-native:subscriptionId)/resourcegroups/$(pulumi config get azure-native:resourceGroupName)
```

Destroy everything but keep the shared (protected resource group)

pulumi destroy --exclude-protected