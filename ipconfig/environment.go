package ipconfig

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/nepomuceno/cloud-ipam/model"
)

const ENVIRONMENT_PARTITION_KEY = "EnvironmentDefinition"

func (client *IpConfigClient) ListEnvironments() ([]model.EnvironmentDefinition, error) {
	result := make([]model.EnvironmentDefinition, 0)
	entities := client.NewListEntitiesPager(&aztables.ListEntitiesOptions{
		Filter: to.Ptr("PartitionKey eq '" + ENVIRONMENT_PARTITION_KEY + "'"),
	})
	for entities.More() {
		pageResp, err := entities.NextPage(client.ctx)
		if err != nil {
			return result, err
		}
		for _, entity := range pageResp.Entities {
			var definition model.EnvironmentDefinition
			err = json.Unmarshal(entity, &definition)
			if err != nil {
				return result, err
			}
			result = append(result, definition)
		}
	}
	return result, nil
}

func (client *IpConfigClient) AddOrUpdateEnvironment(definition model.EnvironmentDefinition) error {
	entity, err := json.Marshal(definition)
	if err != nil {
		return err
	}
	_, err = client.UpsertEntity(client.ctx, entity, &aztables.UpsertEntityOptions{
		UpdateMode: aztables.UpdateModeMerge,
	})
	return err
}

func (client *IpConfigClient) DeleteEnvironment(id string) error {
	// Check if there is any subscription in this environment
	entities := client.NewListEntitiesPager(&aztables.ListEntitiesOptions{
		Filter: to.Ptr("PartitionKey eq '" + id + "'"),
	})
	for entities.More() {
		pageResp, err := entities.NextPage(client.ctx)
		if err != nil {
			return err
		}
		if len(pageResp.Entities) > 0 {
			return fmt.Errorf("environment %s is not empty", id)
		}
	}
	_, err := client.DeleteEntity(client.ctx, ENVIRONMENT_PARTITION_KEY, id, nil)
	return err
}

func (client *IpConfigClient) GetEnvironment(id string) (model.EnvironmentDefinition, error) {
	entity, err := client.GetEntity(client.ctx, ENVIRONMENT_PARTITION_KEY, id, nil)
	if err != nil {
		return model.EnvironmentDefinition{}, err
	}
	var definition model.EnvironmentDefinition
	err = json.Unmarshal(entity.Value, &definition)
	if err != nil {
		return model.EnvironmentDefinition{}, err
	}
	return definition, nil
}
