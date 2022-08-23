package ipconfig

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

type IpConfigClient struct {
	*aztables.Client
	ctx  context.Context
	cred azcore.TokenCredential
}

func GetClient(tableName string, storageName string, ctx context.Context) (*IpConfigClient, error) {
	if storageName == "" {
		return nil, fmt.Errorf("storage-account-name is required")
	}
	serviceURL := fmt.Sprintf("https://%s.table.core.windows.net", storageName)
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}
	serviceClient, err := aztables.NewServiceClient(serviceURL, cred, nil)
	if err != nil {
		return nil, err
	}
	err = createIfDoesntExistTable(tableName, serviceClient, ctx)
	if err != nil {
		return nil, err
	}
	serviceURL = fmt.Sprintf("https://%s.table.core.windows.net/%s", storageName, tableName)
	client, err := aztables.NewClient(serviceURL, cred, nil)
	if err != nil {
		return nil, err
	}
	return &IpConfigClient{client, ctx, cred}, err
}

func (client *IpConfigClient) RemoveTable() error {
	_, err := client.Delete(client.ctx, nil)
	return err
}
