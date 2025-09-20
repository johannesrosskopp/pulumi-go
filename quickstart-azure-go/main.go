package main

import (
	"fmt"

	resources "github.com/pulumi/pulumi-azure-native-sdk/resources/v3"
	storage "github.com/pulumi/pulumi-azure-native-sdk/storage/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

    myname := "johannes"

		rg, err := resources.NewResourceGroup(ctx, "my-resource-group", &resources.ResourceGroupArgs{
			Location:          pulumi.String("westeurope"),
			ResourceGroupName: pulumi.String("pulumi-tech"),
		}, pulumi.Protect(true))
		if err != nil {
			return err
		}

		storageAccount, err := storage.NewStorageAccount(ctx, "storageaccount", &storage.StorageAccountArgs{
			ResourceGroupName: rg.Name,
			Location:          rg.Location,
			Sku: &storage.SkuArgs{
				Name: pulumi.String("Standard_LRS"),
			},
			Kind: pulumi.String("StorageV2"),
      AccountName: pulumi.String(fmt.Sprintf("pulumistorage%s", myname)),
		})
		if err != nil {
			return err
		}

		staticWebsiteIndexDocument := "index.html"
		staticWebsiteIndexDocumentPath := fmt.Sprintf("./www/%s", staticWebsiteIndexDocument)
		staticWebsite, err := storage.NewStorageAccountStaticWebsite(ctx, "storageAccountStaticWebsiteResource", &storage.StorageAccountStaticWebsiteArgs{
			AccountName:       storageAccount.Name,
			ResourceGroupName: rg.Name,
			IndexDocument:     pulumi.String(staticWebsiteIndexDocument),
      
		})
		if err != nil {
			return err
		}

		_, err = storage.NewBlob(ctx, "blobResource", &storage.BlobArgs{
			AccountName:       storageAccount.Name,
			ContainerName:     staticWebsite.ContainerName,
			ResourceGroupName: rg.Name,
			// AccessTier:        storage.BlobAccessTierHot,
			BlobName:    pulumi.String(staticWebsiteIndexDocument),
			ContentType: pulumi.String("text/html"),
			Source:      pulumi.NewFileAsset(staticWebsiteIndexDocumentPath),
			Type:        storage.BlobTypeBlock,
		})
		if err != nil {
			return err
		}

    ctx.Export("primaryWebEndpoint", storageAccount.PrimaryEndpoints.Web())

		return nil
	})
}
