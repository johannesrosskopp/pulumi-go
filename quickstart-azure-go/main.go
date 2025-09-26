package main

import (
	"fmt"

	resources "github.com/pulumi/pulumi-azure-native-sdk/resources/v3"
	storage "github.com/pulumi/pulumi-azure-native-sdk/storage/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		rg, err := resources.NewResourceGroup(ctx, "my-resource-group", &resources.ResourceGroupArgs{
			Location:          pulumi.String("westeurope"),
			ResourceGroupName: pulumi.String("pulumi-tech"),
		}, pulumi.Protect(true))
		if err != nil {
			return err
		}

		// Hmm, maybe that better goes into a config? Try using a namespace, e.g. tech
		myname := XXXXXX

		fmt.Printf("Myname: %s\n", myname)

		storageAccount, err := storage.NewStorageAccount(ctx, "storageaccount", &storage.StorageAccountArgs{
			ResourceGroupName: rg.Name,
			Location:          rg.Location,
			Sku: &storage.SkuArgs{
				Name: pulumi.String("Standard_LRS"),
			},
			Kind:        pulumi.String("StorageV2"),
			AccountName: pulumi.String(fmt.Sprintf("pulumistorage%s", myname)),
		})
		if err != nil {
			return err
		}
		staticWebsiteIndexDocument := "index.html"
		staticWebsite, err := storage.NewStorageAccountStaticWebsite(ctx, "storageAccountStaticWebsiteResource", &storage.StorageAccountStaticWebsiteArgs{
			AccountName:       storageAccount.Name,
			ResourceGroupName: rg.Name,
			IndexDocument:     pulumi.String(staticWebsiteIndexDocument),
		})
		if err != nil {
			return err
		}

		filePath := fmt.Sprintf("./www/%s", staticWebsiteIndexDocument)
		_, err = storage.NewBlob(ctx, "blobResource", &storage.BlobArgs{
			AccountName:       storageAccount.Name,
			ContainerName:     staticWebsite.ContainerName,
			ResourceGroupName: rg.Name,
			BlobName:          pulumi.String(staticWebsiteIndexDocument),
			ContentType:       pulumi.String("text/html"),
			Source:            pulumi.NewFileAsset(filePath),
			Type:              storage.BlobTypeBlock,
		})
		if err != nil {
			return err
		}

		// Q: What happens if we print a string.Output directly?
		// A: ...
		fmt.Printf("Wrong: You can access your website at the following URL: %s\n", storageAccount.PrimaryEndpoints.Web())

		ctx.Export("primaryWebEndpoint", storageAccount.PrimaryEndpoints.Web())

		return nil
	})
}
