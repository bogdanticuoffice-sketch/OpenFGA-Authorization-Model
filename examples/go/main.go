package main

import (
	"context"
	"fmt"
	"log"

	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
)

func main() {
	ctx := context.Background()

	// Store ID is not set at construction — created dynamically below.
	fgaClient, err := client.NewSdkClient(&client.ClientConfiguration{
		ApiUrl: "http://localhost:8080",
	})
	if err != nil {
		log.Fatalf("Failed to create OpenFGA client: %v", err)
	}

	storeID := createStore(ctx, fgaClient)
	fgaClient.SetStoreId(storeID)

	modelID := createAuthorizationModel(ctx, fgaClient)
	fgaClient.SetAuthorizationModelId(modelID)

	createRelationships(ctx, fgaClient)
	checkAccess(ctx, fgaClient)
	listPermissions(ctx, fgaClient)
}

func createStore(ctx context.Context, fgaClient *client.OpenFgaClient) string {
	resp, err := fgaClient.CreateStore(ctx).Body(client.ClientCreateStoreRequest{
		Name: "authorization-store",
	}).Execute()
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}
	fmt.Printf("Created store: %s\n", resp.Id)
	return resp.Id
}

// createAuthorizationModel writes a minimal RBAC model (user / organization / project)
// using TypeDefinitions. For larger models, prefer loading from a .fga file using
// the openfga/language package and its transformer.
func createAuthorizationModel(ctx context.Context, fgaClient *client.OpenFgaClient) string {
	schemaVersion := "1.1"
	thisUserset := openfga.Userset{This: &map[string]interface{}{}}

	resp, err := fgaClient.WriteAuthorizationModel(ctx).Body(client.ClientWriteAuthorizationModelRequest{
		SchemaVersion: schemaVersion,
		TypeDefinitions: []openfga.TypeDefinition{
			{
				Type: "user",
			},
			{
				Type: "organization",
				Relations: &map[string]openfga.Userset{
					"admin":  thisUserset,
					"member": thisUserset,
				},
				Metadata: &openfga.Metadata{
					Relations: &map[string]openfga.RelationMetadata{
						"admin": {
							DirectlyRelatedUserTypes: &[]openfga.RelationReference{
								{Type: "user"},
							},
						},
						"member": {
							DirectlyRelatedUserTypes: &[]openfga.RelationReference{
								{Type: "user"},
							},
						},
					},
				},
			},
			{
				Type: "project",
				Relations: &map[string]openfga.Userset{
					"organization": thisUserset,
					"owner":        thisUserset,
					"editor":       thisUserset,
					"viewer":       thisUserset,
				},
				Metadata: &openfga.Metadata{
					Relations: &map[string]openfga.RelationMetadata{
						"organization": {
							DirectlyRelatedUserTypes: &[]openfga.RelationReference{
								{Type: "organization"},
							},
						},
						"owner": {
							DirectlyRelatedUserTypes: &[]openfga.RelationReference{
								{Type: "user"},
							},
						},
						"editor": {
							DirectlyRelatedUserTypes: &[]openfga.RelationReference{
								{Type: "user"},
							},
						},
						"viewer": {
							DirectlyRelatedUserTypes: &[]openfga.RelationReference{
								{Type: "user"},
							},
						},
					},
				},
			},
		},
	}).Execute()
	if err != nil {
		log.Fatalf("Failed to write authorization model: %v", err)
	}
	fmt.Printf("Authorization model ID: %s\n", resp.AuthorizationModelId)
	return resp.AuthorizationModelId
}

func createRelationships(ctx context.Context, fgaClient *client.OpenFgaClient) {
	err := fgaClient.WriteTuples(ctx).Body([]client.ClientTupleKey{
		{
			User:     "user:alice",
			Relation: "admin",
			Object:   "organization:acme",
		},
		{
			User:     "user:bob",
			Relation: "member",
			Object:   "organization:acme",
		},
		{
			User:     "organization:acme",
			Relation: "organization",
			Object:   "project:api",
		},
		{
			User:     "user:alice",
			Relation: "owner",
			Object:   "project:api",
		},
	}).Execute()
	if err != nil {
		log.Fatalf("Failed to write relationships: %v", err)
	}
	fmt.Println("Relationships created successfully")
}

func checkAccess(ctx context.Context, fgaClient *client.OpenFgaClient) {
	resp, err := fgaClient.Check(ctx).Body(client.ClientCheckRequest{
		User:     "user:alice",
		Relation: "admin",
		Object:   "organization:acme",
	}).Execute()
	if err != nil {
		log.Fatalf("Failed to check access: %v", err)
	}
	fmt.Printf("Alice is admin of acme: %v\n", resp.Allowed)
}

func listPermissions(ctx context.Context, fgaClient *client.OpenFgaClient) {
	resp, err := fgaClient.ListObjects(ctx).Body(client.ClientListObjectsRequest{
		User:     "user:alice",
		Relation: "admin",
		Type:     "organization",
	}).Execute()
	if err != nil {
		log.Fatalf("Failed to list objects: %v", err)
	}
	fmt.Printf("Alice can admin: %v\n", resp.Objects)
}
