package project

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/litmuschaos/litmus/litmus-portal/graphql-server/graph/model"
	"github.com/litmuschaos/litmus/litmus-portal/graphql-server/pkg/database/mongodb"
)

var projectCollection *mongo.Collection

func init() {
	projectCollection = mongodb.Database.Collection("project")
}

// CreateProject ...
func CreateProject(ctx context.Context, project *Project) error {
	// ctx, _ := context.WithTimeout(backgroundContext, 10*time.Second)
	_, err := projectCollection.InsertOne(ctx, project)
	if err != nil {
		return errors.New("Error creating Project: " + err.Error())
	}

	return nil
}

// GetProject ...
func GetProject(ctx context.Context, query bson.D) (*Project, error) {
	// ctx, _ := context.WithTimeout(backgroundContext, 10*time.Second)
	var project = new(Project)
	err := projectCollection.FindOne(ctx, query).Decode(project)
	if err != nil {
		return nil, errors.New("Error getting project " + err.Error())
	}

	return project, err
}

// GetProjectsByUserID ...
func GetProjectsByUserID(ctx context.Context, userID string) ([]Project, error) {
	// ctx, _ := context.WithTimeout(backgroundContext, 10*time.Second)
	projects := []Project{}
	query := bson.M{"members": bson.M{"$elemMatch": bson.M{"user_id": userID, "invitation": bson.M{"$ne": DeclinedInvitation}}}}
	cursor, err := projectCollection.Find(ctx, query)
	if err != nil {
		return nil, errors.New("Error getting project with userID: " + userID + " error:" + err.Error())
	}
	err = cursor.All(ctx, &projects)
	if err != nil {
		return nil, errors.New("Error getting project with userID: " + userID + " error:" + err.Error())
	}

	return projects, err
}

// AddMember ...
func AddMember(ctx context.Context, projectID string, member *Member) error {

	query := bson.M{"_id": projectID}
	update := bson.M{"$push": bson.M{"members": member}}
	_, err := projectCollection.UpdateOne(ctx, query, update)
	if err != nil {
		return errors.New("Error updating project with projectID: " + projectID +  "error: " +  err.Error())
	}
	return nil
}

// RemoveInvitation :Removes member or cancels the invitation
func RemoveInvitation(ctx context.Context, projectID string, userID string, invitation Invitation) error {
	query := bson.M{"_id": projectID}
	update := bson.M{"$pull": bson.M{"members": bson.M{"user_id": userID}}}
	_, err := projectCollection.UpdateOne(ctx, query, update)
	if err != nil {
		if invitation == AcceptedInvitation {
			return errors.New("Error Removing the member with userID:" + userID + "from project with project id: " + projectID + err.Error())
		}
		return errors.New("Error Removing the member with userID:" + userID + "from project with project id: " + projectID + err.Error())
	}
	return nil
}

// UpdateInvite :Updates the status of sent invitation
func UpdateInvite(ctx context.Context, projectID, userID string, invitation Invitation, Role *model.MemberRole) error {
	options := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.user_id": userID},
		},
	})
	query := bson.M{"_id": projectID}
	var update bson.M

	switch invitation {
	case PendingInvitation:
		update = bson.M{"$set": bson.M{"members.$[elem].invitation": invitation, "members.$[elem].role": Role}}
	case DeclinedInvitation:
		update = bson.M{"$set": bson.M{"members.$[elem].invitation": invitation}}
	case AcceptedInvitation:
		update = bson.M{"$set": bson.M{"members.$[elem].invitation": invitation, "members.$[elem].joined_at": time.Now().Format(time.RFC1123Z)}}
	case ExitedProject:
		update = bson.M{"$set": bson.M{"members.$[elem].invitation": invitation}}
	}
	_, err := projectCollection.UpdateOne(ctx, query, update, options)
	if err != nil {
		return errors.New("Error updating project with projectID: " + projectID + " error: " + err.Error())
	}
	return nil
}

// UpdateProjectName :Updates Name of the project
func UpdateProjectName(ctx context.Context, projectID string, projectName string) error {
	query := bson.M{"_id": projectID}
	update := bson.M{"$set": bson.M{"name": projectName}}

	_, err := projectCollection.UpdateOne(ctx, query, update)
	if err != nil {
		return err
	}

	return nil
}
