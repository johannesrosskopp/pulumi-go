package main

import (
	// "fmt"

	"github.com/pulumi/pulumi-command/sdk/go/command/local"

	resources "github.com/pulumi/pulumi-azure-native-sdk/resources/v3"
	// storage "github.com/pulumi/pulumi-azure-native-sdk/storage/v3"
	// web "github.com/pulumi/pulumi-azure-native-sdk/web/v3"
	// "github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	// "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// myname := "johannes"

		
		random, err := local.NewCommand(ctx, "my-bucket", &local.CommandArgs{
			Create: pulumi.String("echo 'New File' > new_file.txt; echo $(pwd)/new_file.txt"),
			Delete: pulumi.String("rm new_file.txt"),
		})
		if err != nil {
			return err
		}

		ctx.Export("output", random.Stdout)

		_, err = resources.NewResourceGroup(ctx, "my-resource-group", &resources.ResourceGroupArgs{
			Location:          pulumi.String("westeurope"),
			ResourceGroupName: pulumi.String("pulumi-tech"),
		}, pulumi.Protect(true))
		if err != nil {
			return err
		}
/*
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

		// Q: What happens if we print a string.Output directly?
		// A: ...
		fmt.Printf("Wrong: You can access your website at the following URL: %s\n", storageAccount.PrimaryEndpoints.Web())
		// Solution:
		storageAccount.PrimaryEndpoints.Web().ApplyT(func(url string) error {
			fmt.Printf("S1: You can access your website at the following URL: %s\n", url)
			return nil
		})

		// Setup Docker image: Configure Docker Hub credentials via config
		// pulumi config set dockerhub:username <your-username>
		// pulumi config set dockerhub:password <your-password> --secret
		dockerCfg := config.New(ctx, "dockerhub")
		dockerUsername := dockerCfg.Require("username")
		dockerPassword := dockerCfg.RequireSecret("password")

		image, err := docker.NewImage(ctx, "app-image", &docker.ImageArgs{
			Build: &docker.DockerBuildArgs{
				Context:    pulumi.String("./app"),
				Dockerfile: pulumi.String("./app/Dockerfile"),
			},
			ImageName: pulumi.String(fmt.Sprintf("%s/pulumi-tech-tmp:timestamp_%s_latest", dockerUsername, myname)),
			Registry: &docker.RegistryArgs{
				Server:   pulumi.String("docker.io"),
				Username: pulumi.String(dockerUsername),
				Password: dockerPassword,
			},
		})
		if err != nil {
			return err
		}

		plan, err := web.NewAppServicePlan(ctx, "app-plan", &web.AppServicePlanArgs{
			ResourceGroupName: rg.Name,
			Location:          rg.Location,
			Name:              pulumi.String("basic-plan"),
			Sku: &web.SkuDescriptionArgs{
				Name: pulumi.String("B1"),
				Tier: pulumi.String("Basic"),
			},
			// Kind: pulumi.String("app,linux,container"),
			// Linux is true, why oh why, no one knows
			Reserved: pulumi.BoolPtr(true),
		})
		if err != nil {
			return err
		}

		image.ImageName.ApplyT(func(imageName string) error {
			fmt.Printf("Image name: %s\n", imageName)
			return nil
		})

		webapp, err := web.NewWebApp(ctx, "webapp", &web.WebAppArgs{
			ResourceGroupName: rg.Name,
			Kind:              pulumi.String("app,linux,container"),
			Location:          rg.Location,
			Name:              pulumi.String(fmt.Sprintf("timeapp-%s", myname)),
			ServerFarmId:      plan.ID(),
			SiteConfig: &web.SiteConfigArgs{
				LinuxFxVersion: pulumi.Sprintf("DOCKER|%s", image.ImageName),
				AppSettings: web.NameValuePairArray{
					&web.NameValuePairArgs{
						Name:  pulumi.String("WEBSITES_ENABLE_APP_SERVICE_STORAGE"),
						Value: pulumi.String("false"),
					},
				},
			},
		})
		if err != nil {
			return err
		}

		// Run a docker build

		ctx.Export("primaryWebEndpoint", storageAccount.PrimaryEndpoints.Web())

		ctx.Export("apiUrl", pulumi.Sprintf("http://%s", webapp.DefaultHostName))*/

		return nil
	})
}
