package main

import (
	"fmt"
	"github.com/Azure/azure-automation-go-worker/internal/configuration"
	"github.com/Azure/azure-automation-go-worker/internal/jrds"
	"os"
)

func main() {
	fmt.Println("Sandbox starting")

	if len(os.Args) < 2 {
		panic("missing sandbox.exe parameter")
	}
	sandboxId := os.Args[1]
	fmt.Println("Arg %v ", os.Args)
	httpClient := jrds.NewSecureHttpClient(configuration.GetJrdsCertificatePath(), configuration.GetJrdsKeyPath())
	jrdsClient := jrds.NewJrdsClient(&httpClient, configuration.GetJrdsBaseUri(), configuration.GetAccountId(), configuration.GetHybridWorkerGroupName())

	jobActions := jrds.JobActions{}
	err := jrdsClient.GetJobActions(sandboxId, &jobActions)
	if err != nil {
		fmt.Printf("error getting jobaction %v", err)
	}

	var arr []jrds.MessageMetadata
	for _, jobaction := range jobActions.Value {
		arr = append(arr, *jobaction.MessageMetadata)
	}

	metadatas := jrds.MessageMetadatas{arr}
	err = jrdsClient.AcknowledgeJobAction(sandboxId, metadatas)
	if err != nil {
		fmt.Printf("error getting messageMetadata %v", err)
	}

	jobData := jrds.JobData{}
	err = jrdsClient.GetJobData(*jobActions.Value[0].JobId, &jobData)
	if err != nil {
		fmt.Printf("error getting jobData %v", err)
	}
	fmt.Printf("JobData%v", jobData)

	jobUpdatableData := jrds.JobUpdatableData{}
	err = jrdsClient.GetUpdatableJobData(*jobActions.Value[0].JobId, &jobUpdatableData)
	if err != nil {
		fmt.Printf("error getting jobupdatable Data %v", err)
	}

	jrdsClient.SetJobStatus(sandboxId, *jobActions.Value[0].JobId, 3, false, nil)
	err = jrdsClient.UnloadJob(*jobUpdatableData.SubscriptionId, sandboxId, *jobActions.Value[0].JobId, false, "")
	if err != nil {
		fmt.Printf("error unloading %v", err)
	}

	a := string("exception")
	jrdsClient.SetJobStatus(sandboxId, *jobActions.Value[0].JobId, 4, true, &a)
	fmt.Println("Exiting..")
}
