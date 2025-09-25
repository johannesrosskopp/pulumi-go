package main

import (
	"fmt"

	"github.com/pulumi/pulumi-command/sdk/go/command/local"

	resources "github.com/pulumi/pulumi-azure-native-sdk/resources/v3"
	storage "github.com/pulumi/pulumi-azure-native-sdk/storage/v3"
	web "github.com/pulumi/pulumi-azure-native-sdk/web/v3"
	"github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
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

		apiUrl, err := createApi(ctx, rg)
		if err != nil {
			return err
		}

		err = createStaticWebsite(ctx, StaticWebsiteArgs{
			rg:     rg,
			apiUrl: apiUrl,
		})
		if err != nil {
			return err
		}

		return nil
	})
}

func createApi(ctx *pulumi.Context, rg *resources.ResourceGroup) (pulumi.StringOutput, error) {
	techconfig := config.New(ctx, "tech")
	myname := techconfig.Require("myname")

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
		return pulumi.StringOutput{}, err
	}

	plan, err := web.NewAppServicePlan(ctx, "app-plan", &web.AppServicePlanArgs{
		ResourceGroupName: rg.Name,
		Location:          rg.Location,
		Name:              pulumi.String(fmt.Sprintf("basic-plan-%s", myname)),
		Sku: &web.SkuDescriptionArgs{
			Name: pulumi.String("B1"),
			Tier: pulumi.String("Basic"),
		},
		// Linux is true, why oh why, no one knows
		Reserved: pulumi.BoolPtr(true),
	})
	if err != nil {
		return pulumi.StringOutput{}, err
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
			Cors: &web.CorsSettingsArgs{
				AllowedOrigins: pulumi.StringArray{
					pulumi.String("*"), // Allow all origins - you can restrict this later
				},
			},
		},
	})
	if err != nil {
		return pulumi.StringOutput{}, err
	}

	apiUrl := pulumi.Sprintf("https://%s", webapp.DefaultHostName)

	ctx.Export("apiUrl", apiUrl)

	return apiUrl, nil
}

type StaticWebsiteArgs struct {
	rg     *resources.ResourceGroup
	apiUrl pulumi.StringOutput
}

func createStaticWebsite(ctx *pulumi.Context, args StaticWebsiteArgs) error {
	techconfig := config.New(ctx, "tech")
	myname := techconfig.Require("myname")

	fmt.Printf("Myname: %s\n", myname)

	storageAccount, err := storage.NewStorageAccount(ctx, "storageaccount", &storage.StorageAccountArgs{
		ResourceGroupName: args.rg.Name,
		Location:          args.rg.Location,
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
		ResourceGroupName: args.rg.Name,
		IndexDocument:     pulumi.String(staticWebsiteIndexDocument),
	})
	if err != nil {
		return err
	}

	// Read the HTML file to use as a trigger for rebuilding when it changes
	htmlAsset := pulumi.NewFileAsset("./www/index.html")

	staticPageHTML, err := local.NewCommand(ctx, "my-bucket", &local.CommandArgs{
		Update: pulumi.String("ls -la .www/; sed -e \"s|API_URL|$ENV_API_URL|g\" ./www/index.html"),
		Environment: pulumi.StringMap{
			"ENV_API_URL": args.apiUrl,
		},
		Triggers: pulumi.Array{
			htmlAsset,
		},
	})
	if err != nil {
		return err
	}

	staticPageHTML.Stdout.ApplyT(func(staticPageHTML string) error {
		fmt.Printf("Static page HTML: %s\n", staticPageHTML)
		_, err = storage.NewBlob(ctx, "blobResource", &storage.BlobArgs{
			AccountName:       storageAccount.Name,
			ContainerName:     staticWebsite.ContainerName,
			ResourceGroupName: args.rg.Name,
			BlobName:          pulumi.String(staticWebsiteIndexDocument),
			ContentType:       pulumi.String("text/html"),
			Source:            pulumi.NewStringAsset(staticPageHTML),
			Type:              storage.BlobTypeBlock,
		})
		if err != nil {
			return err
		}

		// Upload favicon if it exists
		_, err = storage.NewBlob(ctx, "faviconResource", &storage.BlobArgs{
			AccountName:       storageAccount.Name,
			ContainerName:     staticWebsite.ContainerName,
			ResourceGroupName: args.rg.Name,
			BlobName:          pulumi.String("favicon.ico"),
			ContentType:       pulumi.String("image/x-icon"),
			Source:            pulumi.NewFileAsset("./www/favicon.ico"),
			Type:              storage.BlobTypeBlock,
		})
		if err != nil {
			// If favicon doesn't exist, that's OK - just continue
			fmt.Printf("Note: favicon.ico not found in ./www/ - skipping favicon upload\n")
		}

		return nil
	})

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
}
