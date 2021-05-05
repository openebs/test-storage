package workflow

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/litmuschaos/litmus/litmus-portal/graphql-server/pkg/database/mongodb"
)

var (
	workflowCollection *mongo.Collection
	backgroundContext  = context.Background()
)

func init() {
	workflowCollection = mongodb.Database.Collection("workflow-collection")
	_, err := workflowCollection.Indexes().CreateMany(backgroundContext, []mongo.IndexModel{
		{
			Keys: bson.M{
				"workflow_id": 1,
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{
				"workflow_name": 1,
			},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		log.Fatal("Error Creating Index for Workflow Collection: ", err)
	}
}

// UpdateWorkflowRun takes workflowID and wfRun parameters to update the workflow run details in the database
func UpdateWorkflowRun(workflowID string, wfRun ChaosWorkflowRun) (int, error) {
	ctx, cancel := context.WithTimeout(backgroundContext, 10*time.Second)
	defer cancel()
	count, err := workflowCollection.CountDocuments(ctx, bson.M{"workflow_id": workflowID, "workflow_runs.workflow_run_id": wfRun.WorkflowRunID})
	if err != nil {
		return 0, err
	}

	updateCount := 1
	if count == 0 {
		filter := bson.M{"workflow_id": workflowID}
		update := bson.M{"$push": bson.M{"workflow_runs": wfRun}}
		updateResp, err := workflowCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return 0, err
		}
		if updateResp.MatchedCount == 0 {
			return 0, errors.New("workflow not found")
		}
	} else if count == 1 {
		filter := bson.M{"workflow_id": workflowID, "workflow_runs": bson.M{"$elemMatch": bson.M{"workflow_run_id": wfRun.WorkflowRunID, "completed": false}}}
		update := bson.M{"$set": bson.M{"workflow_runs.$.last_updated": wfRun.LastUpdated, "workflow_runs.$.execution_data": wfRun.ExecutionData, "workflow_runs.$.completed": wfRun.Completed}}
		updateResp, err := workflowCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return 0, err
		}
		updateCount = int(updateResp.MatchedCount)
	}

	return updateCount, nil
}

// GetWorkflows takes a query parameter to retrieve the workflow details from the database
func GetWorkflows(query bson.D) ([]ChaosWorkFlowInput, error) {
	ctx, cancel := context.WithTimeout(backgroundContext, 10*time.Second)
	defer cancel()
	cursor, err := workflowCollection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	var workflows []ChaosWorkFlowInput
	err = cursor.All(ctx, &workflows)
	if err != nil {
		return nil, err
	}

	return workflows, nil
}

// GetWorkflowsByClusterID takes a clusterID parameter to retrieve the workflow details from the database
func GetWorkflowsByClusterID(clusterID string) ([]ChaosWorkFlowInput, error) {
	query := bson.D{{Key: "cluster_id", Value: clusterID}}
	ctx, cancel := context.WithTimeout(backgroundContext, 10*time.Second)
	defer cancel()
	cursor, err := workflowCollection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	var workflows []ChaosWorkFlowInput
	err = cursor.All(ctx, &workflows)
	if err != nil {
		return nil, err
	}
	return workflows, nil
}

// InsertChaosWorkflow takes details of a workflow and inserts into the database collection
func InsertChaosWorkflow(chaosWorkflow ChaosWorkFlowInput) error {
	ctx, cancel := context.WithTimeout(backgroundContext, 10*time.Second)
	defer cancel()
	_, err := workflowCollection.InsertOne(ctx, chaosWorkflow)
	if err != nil {
		return err
	}

	return nil
}

// UpdateChaosWorkflow takes query and update parameters to update the workflow details in the database
func UpdateChaosWorkflow(query bson.D, update bson.D) error {
	ctx, cancel := context.WithTimeout(backgroundContext, 10*time.Second)
	defer cancel()
	_, err := workflowCollection.UpdateOne(ctx, query, update)
	if err != nil {
		return err
	}

	return nil
}
