package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/openfga/go-sdk/client"
)

func BenchmarkCheck(b *testing.B) {
	storeId := "01HH7JXP37M6YFCRQVHN89DFFP"
	modelId := "01HH7JXP3W7CBCC2T1YVF75VW6"

	fgaClient, err := client.NewSdkClient(&client.ClientConfiguration{
		ApiScheme:            "http",
		ApiHost:              "localhost:8080", // required, define without the scheme (e.g. api.fga.example instead of https://api.fga.example)
		AuthorizationModelId: &modelId,
		StoreId:              storeId,
	})
	if err != nil {
		log.Printf("Failed to create client: %s", err)
		return
	}
	b.ResetTimer()
	for i := 0; i < 100; i++ {
		//bu := fmt.Sprintf("bu:bu%d", j)
		for j := 0; j < 100; j++ {
			user := fmt.Sprintf("user:u_%d_%d", j, i)
			checks(fgaClient, modelId, user, "can_read", "document:doc_0_0", true)
		}
	}
}
