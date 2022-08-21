package ipconfig

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

func isTableAlreadyExists(err error) bool {
	return err.(*azcore.ResponseError).ErrorCode == "TableAlreadyExists"
}

func createIfDoesntExistTable(tableName string, serviceClient *aztables.ServiceClient, ctx context.Context) error {
	_, err := serviceClient.CreateTable(ctx, tableName, nil)
	if err != nil {
		if !isTableAlreadyExists(err) {
			return err
		}
	}
	return nil
}
