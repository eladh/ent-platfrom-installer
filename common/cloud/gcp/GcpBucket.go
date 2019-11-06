package gcp

import (
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/option"
	"log"
	"os"
)

func CreateGcpBucket(projectId string, bucketName string, storageClass string, location string) {
	ctx := context.Background()

	// Creates a client
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Creates the new bucket with params
	if err := createBucketWithAttrs(client, projectId, bucketName, storageClass, location); err != nil {
		log.Fatalf("Failed to create bucket: %v", err)
	}
}

func createBucketWithAttrs(client *storage.Client, projectID, bucketName string, storageClass string, location string) error {
	ctx := context.Background()
	bucket := client.Bucket(bucketName)

	if err := bucket.Create(ctx, projectID, &storage.BucketAttrs{
		StorageClass: storageClass,
		Location:     location,
	}); err != nil {
		return err
	}

	return nil
}
