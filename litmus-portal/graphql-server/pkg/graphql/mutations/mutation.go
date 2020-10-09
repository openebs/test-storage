package mutations

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/copier"

	"github.com/google/uuid"
	"github.com/litmuschaos/litmus/litmus-portal/graphql-server/graph/model"
	"github.com/litmuschaos/litmus/litmus-portal/graphql-server/pkg/cluster"
	store "github.com/litmuschaos/litmus/litmus-portal/graphql-server/pkg/data-store"
	database "github.com/litmuschaos/litmus/litmus-portal/graphql-server/pkg/database/mongodb"
	"github.com/litmuschaos/litmus/litmus-portal/graphql-server/pkg/graphql/subscriptions"
	"github.com/litmuschaos/litmus/litmus-portal/graphql-server/utils"
	"github.com/pkg/errors"
	"github.com/tidwall/sjson"
	"go.mongodb.org/mongo-driver/bson"
)

//ClusterRegister creates an entry for a new cluster in DB and generates the url used to apply manifest
func ClusterRegister(input model.ClusterInput) (*model.ClusterRegResponse, error) {
	clusterID := uuid.New().String()

	token, err := cluster.ClusterCreateJWT(clusterID)
	if err != nil {
		return &model.ClusterRegResponse{}, err
	}

	newCluster := database.Cluster{
		ClusterID:    clusterID,
		ClusterName:  input.ClusterName,
		Description:  input.Description,
		ProjectID:    input.ProjectID,
		AccessKey:    utils.RandomString(32),
		ClusterType:  input.ClusterType,
		PlatformName: input.PlatformName,
		CreatedAt:    strconv.FormatInt(time.Now().Unix(), 10),
		UpdatedAt:    strconv.FormatInt(time.Now().Unix(), 10),
		Token:        token,
	}

	err = database.InsertCluster(newCluster)
	if err != nil {
		return &model.ClusterRegResponse{}, err
	}

	log.Print("NEW CLUSTER REGISTERED : ID-", clusterID, " PID-", input.ProjectID)

	return &model.ClusterRegResponse{
		ClusterID:   newCluster.ClusterID,
		Token:       token,
		ClusterName: newCluster.ClusterName,
	}, nil
}

//ConfirmClusterRegistration takes the cluster_id and access_key from the subscriber and validates it, if validated generates and sends new access_key
func ConfirmClusterRegistration(identity model.ClusterIdentity, r store.StateData) (*model.ClusterConfirmResponse, error) {
	cluster, err := database.GetCluster(identity.ClusterID)
	if err != nil {
		return &model.ClusterConfirmResponse{IsClusterConfirmed: false}, err
	}

	if cluster.AccessKey == identity.AccessKey {
		newKey := utils.RandomString(32)
		time := strconv.FormatInt(time.Now().Unix(), 10)
		query := bson.D{{"cluster_id", identity.ClusterID}}
		update := bson.D{{"$unset", bson.D{{"token", ""}}}, {"$set", bson.D{{"access_key", newKey}, {"is_registered", true}, {"updated_at", time}}}}

		err = database.UpdateCluster(query, update)
		if err != nil {
			return &model.ClusterConfirmResponse{IsClusterConfirmed: false}, err
		}

		cluster.IsRegistered = true
		cluster.AccessKey = ""

		newCluster := model.Cluster{}
		copier.Copy(&newCluster, &cluster)

		log.Print("CLUSTER Confirmed : ID-", cluster.ClusterID, " PID-", cluster.ProjectID)
		subscriptions.SendClusterEvent("cluster-registration", "New Cluster", "New Cluster registration", newCluster, r)

		return &model.ClusterConfirmResponse{IsClusterConfirmed: true, NewClusterKey: &newKey, ClusterID: &cluster.ClusterID}, err
	}
	return &model.ClusterConfirmResponse{IsClusterConfirmed: false}, err
}

//NewEvent takes a event from a subscriber, validates identity and broadcasts the event to the users
func NewEvent(clusterEvent model.ClusterEventInput, r store.StateData) (string, error) {
	cluster, err := database.GetCluster(clusterEvent.ClusterID)
	if err != nil {
		return "", err
	}

	if cluster.AccessKey == clusterEvent.AccessKey && cluster.IsRegistered {
		log.Print("CLUSTER EVENT : ID-", cluster.ClusterID, " PID-", cluster.ProjectID)

		newCluster := model.Cluster{}
		copier.Copy(&newCluster, &cluster)

		subscriptions.SendClusterEvent("cluster-event", clusterEvent.EventName, clusterEvent.Description, newCluster, r)
		return "Event Published", nil
	}

	return "", errors.New("ERROR WITH CLUSTER EVENT")
}

// WorkFlowRunHandler Updates or Inserts a new Workflow Run into the DB
func WorkFlowRunHandler(input model.WorkflowRunInput, r store.StateData) (string, error) {
	cluster, err := cluster.VerifyCluster(*input.ClusterID)
	if err != nil {
		log.Print("ERROR", err)
		return "", err
	}

	//err = database.UpsertWorkflowRun(database.WorkflowRun(newWorkflowRun))
	err = database.UpsertWorkflowRun(input.WorkflowID, database.WorkflowRun{
		WorkflowRunID: input.WorkflowRunID,
		LastUpdated:   strconv.FormatInt(time.Now().Unix(), 10),
		ExecutionData: input.ExecutionData,
	})
	if err != nil {
		log.Print("ERROR", err)
		return "", err
	}

	subscriptions.SendWorkflowEvent(model.WorkflowRun{
		ClusterID:     cluster.ClusterID,
		ClusterName:   cluster.ClusterName,
		ProjectID:     cluster.ProjectID,
		LastUpdated:   strconv.FormatInt(time.Now().Unix(), 10),
		WorkflowRunID: input.WorkflowRunID,
		WorkflowName:  input.WorkflowName,
		ExecutionData: input.ExecutionData,
		WorkflowID:    input.WorkflowID,
	}, r)

	return "Workflow Run Accepted", nil
}

// LogsHandler receives logs from the workflow-agent and publishes to frontend clients
func LogsHandler(podLog model.PodLog, r store.StateData) (string, error) {
	_, err := cluster.VerifyCluster(*podLog.ClusterID)
	if err != nil {
		log.Print("ERROR", err)
		return "", err
	}
	if reqChan, ok := r.WorkflowLog[podLog.RequestID]; ok {
		resp := model.PodLogResponse{
			PodName:       podLog.PodName,
			WorkflowRunID: podLog.WorkflowRunID,
			PodType:       podLog.PodType,
			Log:           podLog.Log,
		}
		reqChan <- &resp
		close(reqChan)
		return "LOGS SENT SUCCESSFULLY", nil
	}
	return "LOG REQUEST CANCELLED", nil
}

func CreateChaosWorkflow(input *model.ChaosWorkFlowInput, r store.StateData) (*model.ChaosWorkFlowResponse, error) {

	var Weightages []*database.WeightagesInput
	copier.Copy(&Weightages, &input.Weightages)

	workflow_id := uuid.New().String()

	var workflow map[string]interface{}
	err := json.Unmarshal([]byte(input.WorkflowManifest), &workflow)
	if err != nil {
		return nil, err
	}

	newWorkflowManifest, _ := sjson.Set(input.WorkflowManifest, "metadata.labels.workflow_id", workflow_id)
	if strings.ToLower(workflow["kind"].(string)) == "cronworkflow" {
		newWorkflowManifest, _ = sjson.Set(input.WorkflowManifest, "spec.workflowMetadata.labels.workflow_id", workflow_id)
	}

	newChaosWorkflow := database.ChaosWorkFlowInput{
		WorkflowID:          workflow_id,
		WorkflowManifest:    newWorkflowManifest,
		CronSyntax:          input.CronSyntax,
		WorkflowName:        input.WorkflowName,
		WorkflowDescription: input.WorkflowDescription,
		IsCustomWorkflow:    input.IsCustomWorkflow,
		ProjectID:           input.ProjectID,
		ClusterID:           input.ClusterID,
		Weightages:          Weightages,
		CreatedAt:           strconv.FormatInt(time.Now().Unix(), 10),
		UpdatedAt:           strconv.FormatInt(time.Now().Unix(), 10),
		WorkflowRuns:        []*database.WorkflowRun{},
	}

	err = database.InsertChaosWorkflow(newChaosWorkflow)
	if err != nil {
		return nil, err
	}

	subscriptions.SendWorkflowRequest(&newChaosWorkflow, r)

	return &model.ChaosWorkFlowResponse{
		WorkflowID:          workflow_id,
		CronSyntax:          input.CronSyntax,
		WorkflowName:        input.WorkflowName,
		WorkflowDescription: input.WorkflowDescription,
		IsCustomWorkflow:    input.IsCustomWorkflow,
	}, nil
}
