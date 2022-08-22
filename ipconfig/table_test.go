package ipconfig_test

import (
	"context"
	"time"

	"github.com/mepomuceno/cloud-ipam/ipconfig"
)

func getTestClient(storageAccountName string) (*ipconfig.IpConfigClient, error) {
	now := time.Now()
	tableName := "TEST" + now.Format("20060102150405")
	return ipconfig.GetClient(tableName, storageAccountName, context.Background())
}
