package main

import (
	"fmt"

	resources "github.com/pulumi/pulumi-azure-native-sdk/resources/v3"
	storage "github.com/pulumi/pulumi-azure-native-sdk/storage/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
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
		// myname := XXXXXX
		// Solution:
		techconfig := config.New(ctx, "tech")
		myname := techconfig.Require("myname")

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
		// staticWebsiteIndexDocumentPath := fmt.Sprintf("./www/%s", staticWebsiteIndexDocument)
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

		// Upload favicon if it exists
		_, err = storage.NewBlob(ctx, "faviconResource", &storage.BlobArgs{
			AccountName:       storageAccount.Name,
			ContainerName:     staticWebsite.ContainerName,
			ResourceGroupName: rg.Name,
			BlobName:          pulumi.String("favicon.ico"),
			ContentType:       pulumi.String("image/x-icon"),
			Source:            pulumi.NewFileAsset("./www/favicon.ico"),
			Type:              storage.BlobTypeBlock,
		})
		if err != nil {
			// If favicon doesn't exist, that's OK - just continue
			fmt.Printf("Note: favicon.ico not found in ./www/ - skipping favicon upload\n")
		}

		// Q: What happens if we print a string.Output directly?
		// A: ...
		fmt.Printf("Wrong: You can access your website at the following URL: %s\n", storageAccount.PrimaryEndpoints.Web())
		// Solution:
		storageAccount.PrimaryEndpoints.Web().ApplyT(func(url string) error {
			fmt.Printf("S1: You can access your website at the following URL: %s\n", url)
			return nil
		})

		ctx.Export("primaryWebEndpoint", storageAccount.PrimaryEndpoints.Web())

		return nil
	})
}
