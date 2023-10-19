package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/joho/godotenv"
)

const containerName = "files"

type Functions interface {
	GetFiles() ([]string, error)
	Upload([]byte, string) error
	Download(string) ([]byte, error)
}

type AzureFunctions struct {
	subscriptionId		string
	fileBaseUrl			string
}

func NewAzureService() *AzureFunctions {
	return &AzureFunctions{
		subscriptionId: 	loadEnvVar("AZURE_SUBSCRIPTION_ID"),
		fileBaseUrl: 		loadEnvVar("AZURE_FILE_BASE_URL"),
	}
}

func loadEnvVar(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	return os.Getenv(key)
}

func (az *AzureFunctions) GetFiles() ([]string, error) {
	ctx := context.Background()

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	client, err := azblob.NewClient(az.fileBaseUrl, credential, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// List files
	pager := client.NewListBlobsFlatPager(containerName, &azblob.ListBlobsFlatOptions{
		Include: azblob.ListBlobsInclude{ Snapshots: true, Versions: true },
	})

	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != err {
			log.Fatal(err)
			return nil, err
		}

		foundBlobs := []string{}
		for _, blob := range resp.Segment.BlobItems {
			foundBlobs = append(foundBlobs, *blob.Name)
		}

		return foundBlobs, nil
	}

	return nil, nil
}

func (az *AzureFunctions) Upload(data []byte, filename string) error {
	ctx := context.Background()

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	client, err := azblob.NewClient(az.fileBaseUrl, credential, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Upload
	_, err = client.UploadBuffer(ctx, containerName, filename, data, &azblob.UploadBufferOptions {})
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (az *AzureFunctions) Download(filename string) ([]byte, error) {
	ctx := context.Background()

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	client, err := azblob.NewClient(az.fileBaseUrl, credential, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Download
	resp, err := client.DownloadStream(ctx, containerName, filename, &azblob.DownloadStreamOptions{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	
	reader := resp.Body
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return data, nil
}