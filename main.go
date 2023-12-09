package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/openfga/go-sdk/client"
	"github.com/openfga/language/pkg/go/transformer"
)

var model = `
model
  schema 1.1

type user
type bu
  relations
    define member: [user]

type document
  relations
    define reader: [user]
    define writer: [user]
    define owner: [user, bu#member]
	define can_read: reader or writer or owner
`

func main() {

	storeId := "01HH7JXP37M6YFCRQVHN89DFFP"
	modelId := "01HH7JXP3W7CBCC2T1YVF75VW6"
	createUsers := false

	fgaClient, err := client.NewSdkClient(&client.ClientConfiguration{
		ApiScheme: "http",
		ApiHost:   "localhost:8080", // required, define without the scheme (e.g. api.fga.example instead of https://api.fga.example)
	})
	if err != nil {
		log.Printf("Failed to create client: %s", err)
		return
	}
	if storeId == "" {
		resp, err := fgaClient.CreateStore(context.Background()).Body(client.ClientCreateStoreRequest{Name: "FGA Demo"}).Execute()
		if err != nil {
			log.Printf("Failed to create store: %s", err)
			return
		}
		log.Printf("Store = %s", *resp.Id)
		storeId = *resp.Id
	}

	fgaClient.SetStoreId(storeId)

	if modelId == "" {
		jsonm, err := convert(model)

		if err != nil {
			log.Printf("Erro in convert: %s", err)
			return
		}

		var body client.ClientWriteAuthorizationModelRequest
		if err := json.Unmarshal([]byte(jsonm), &body); err != nil {
			log.Printf("error: %v", err)
			return
		}

		model, err := fgaClient.WriteAuthorizationModel(context.Background()).Body(body).Execute()

		if err != nil {
			log.Printf("Erro in wr model: %s", err)
			return
		}

		log.Printf("Model ID = %s", *model.AuthorizationModelId)
		modelId = *model.AuthorizationModelId
	}

	if createUsers {
		createAllUsers(fgaClient, modelId)
	}
	//createUsers(fgaClient, storeId, modelId, "user:u1", "member", "bu:bu1")
	//createUsers(fgaClient, storeId, modelId, "user:u2", "member", "bu:bu1")
	//createUsers(fgaClient, storeId, modelId, "bu:bu1#member", "owner", "document:doc1")
	//checks(fgaClient, modelId, "user:u1", "member", "bu:bu1")
	//checks(fgaClient, modelId, "user:u1", "can_read", "document:doc1")
	//checks(fgaClient, modelId, "user:u2", "can_read", "document:doc1")

	checks(fgaClient, modelId, "user:u_0_0", "can_read", "document:doc_0_0", false)

}

func createAllUsers(fgaclient *client.OpenFgaClient, modelId string) {

	// Create docs
	for j := 0; j < 100; j++ {
		tuples := []client.ClientTupleKey{}
		for i := 0; i < 100; i++ {
			tuple := client.ClientTupleKey{
				User:     fmt.Sprintf("bu:bu%d#member", j),
				Relation: "owner",
				Object:   fmt.Sprintf("document:doc_%d_%d", j, i),
			}
			tuples = append(tuples, tuple)
		}

		body := client.ClientWriteRequest{
			Writes: &tuples,
		}
		options := client.ClientWriteOptions{
			AuthorizationModelId: &modelId,
		}

		data, err := fgaclient.Write(context.Background()).Body(body).Options(options).Execute()
		if err != nil {
			log.Printf("Error add: %s", err)
		} else {
			for _, wr := range data.Writes {
				if wr.Status != client.SUCCESS {
					log.Printf("%s", wr.Status)
				}
			}

		}
	}

	// Assign BU0 to BU99 to users
	// user:u_0_xxx belongs to BU0
	for j := 0; j < 100; j++ {
		tuples := []client.ClientTupleKey{}
		for i := 0; i < 100; i++ {
			tuple := client.ClientTupleKey{
				User:     fmt.Sprintf("user:u_%d_%d", j, i),
				Relation: "member",
				Object:   fmt.Sprintf("bu:bu%d", j),
			}
			tuples = append(tuples, tuple)
		}

		body := client.ClientWriteRequest{
			Writes: &tuples,
		}
		options := client.ClientWriteOptions{
			AuthorizationModelId: &modelId,
		}

		data, err := fgaclient.Write(context.Background()).Body(body).Options(options).Execute()
		if err != nil {
			log.Printf("Error add: %s", err)
		} else {
			for _, wr := range data.Writes {
				if wr.Status != client.SUCCESS {
					log.Printf("%s", wr.Status)
				}
			}

		}
	}

}
func convert(dslString string) (string, error) {
	// Transform from DSL to a JSON string
	jsonStringModel, err := transformer.TransformDSLToJSON(dslString)
	if err != nil {
		log.Printf("Err translating: %s", err)
	}
	return jsonStringModel, err
}

func checks(fgaclient *client.OpenFgaClient, modelid string, user string, relation string, object string, silent bool) bool {
	body := client.ClientCheckRequest{
		User:     user,
		Relation: relation,
		Object:   object,
	}
	options := client.ClientCheckOptions{
		AuthorizationModelId: &modelid,
	}

	data, err := fgaclient.Check(context.Background()).Body(body).Options(options).Execute()
	if err != nil {
		log.Printf("check err: %s", err)
		return false
	}
	if !silent {
		log.Printf("res=%t", data.GetAllowed())
	}
	return data.GetAllowed()
}

func createUsers(fgaclient *client.OpenFgaClient, modelid string, user, relation, object string) {

	body := client.ClientWriteRequest{
		Writes: &[]client.ClientTupleKey{{
			User:     user,
			Relation: relation,
			Object:   object,
		}},
	}
	options := client.ClientWriteOptions{
		AuthorizationModelId: &modelid,
	}

	data, err := fgaclient.Write(context.Background()).Body(body).Options(options).Execute()
	if err != nil {
		log.Printf("Error add: %s", err)
	} else {
		for _, wr := range data.Writes {
			if wr.Status != client.SUCCESS {
				log.Printf("%s", wr.Status)
			}
		}

	}

}
