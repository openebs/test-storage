package handler

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/litmuschaos/litmus/litmus-portal/graphql-server/graph/model"
	"github.com/litmuschaos/litmus/litmus-portal/graphql-server/pkg/analytics"
	"github.com/litmuschaos/litmus/litmus-portal/graphql-server/pkg/analytics/ops/prometheus"
	dbOperations "github.com/litmuschaos/litmus/litmus-portal/graphql-server/pkg/database/mongodb/operations"
	dbSchema "github.com/litmuschaos/litmus/litmus-portal/graphql-server/pkg/database/mongodb/schema"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"strconv"
	"time"
)

func CreateDataSource(datasource *model.DSInput) (*model.DSResponse, error) {

	datasourceStatus := prometheus.TSDBHealthCheck(datasource.DsURL, datasource.DsType)

	if datasourceStatus == "Active" {

		newDS := dbSchema.DataSource{
			DsID:              uuid.New().String(),
			DsName:            datasource.DsName,
			DsType:            datasource.DsType,
			DsURL:             datasource.DsURL,
			AccessType:        datasource.AccessType,
			AuthType:          datasource.AuthType,
			BasicAuthUsername: datasource.BasicAuthUsername,
			BasicAuthPassword: datasource.BasicAuthPassword,
			ScrapeInterval:    datasource.ScrapeInterval,
			QueryTimeout:      datasource.QueryTimeout,
			HTTPMethod:        datasource.HTTPMethod,
			ProjectID:         *datasource.ProjectID,
			CreatedAt:         strconv.FormatInt(time.Now().Unix(), 10),
			UpdatedAt:         strconv.FormatInt(time.Now().Unix(), 10),
		}

		err := dbOperations.InsertDataSource(newDS)
		if err != nil {
			return nil, err
		}

		var newDSResponse = model.DSResponse{}
		_ = copier.Copy(&newDSResponse, &newDS)

		return &newDSResponse, nil
	} else {
		return nil, errors.New("Datasource: " + datasourceStatus)
	}
}

func CreateDashboard(dashboard *model.CreateDBInput) (string, error) {

	newDashboard := dbSchema.DashBoard{
		DbID:        uuid.New().String(),
		DbName:      dashboard.DbName,
		DbType:      dashboard.DbType,
		DsID:        dashboard.DsID,
		EndTime:     dashboard.EndTime,
		StartTime:   dashboard.StartTime,
		RefreshRate: dashboard.RefreshRate,
		ClusterID:   dashboard.ClusterID,
		ProjectID:   dashboard.ProjectID,
		IsRemoved:   false,
		CreatedAt:   strconv.FormatInt(time.Now().Unix(), 10),
		UpdatedAt:   strconv.FormatInt(time.Now().Unix(), 10),
	}

	var (
		newPanelGroups = make([]dbSchema.PanelGroup, len(dashboard.PanelGroups))
		newPanels      []*dbSchema.Panel
	)

	for i, panel_group := range dashboard.PanelGroups {

		panelGroupID := uuid.New().String()
		newPanelGroups[i].PanelGroupID = panelGroupID
		newPanelGroups[i].PanelGroupName = panel_group.PanelGroupName

		for _, panel := range panel_group.Panels {
			var newPromQueries []*dbSchema.PromQuery
			copier.Copy(&newPromQueries, &panel.PromQueries)

			var newPanelOptions dbSchema.PanelOption
			copier.Copy(&newPanelOptions, &panel.PanelOptions)

			newPanel := dbSchema.Panel{
				PanelID:      uuid.New().String(),
				PanelOptions: &newPanelOptions,
				PanelName:    panel.PanelName,
				PanelGroupID: panelGroupID,
				PromQueries:  newPromQueries,
				IsRemoved:    false,
				XAxisDown:    panel.XAxisDown,
				YAxisLeft:    panel.YAxisLeft,
				YAxisRight:   panel.YAxisRight,
				Unit:         panel.Unit,
				CreatedAt:    strconv.FormatInt(time.Now().Unix(), 10),
				UpdatedAt:    strconv.FormatInt(time.Now().Unix(), 10),
			}

			newPanels = append(newPanels, &newPanel)
		}
	}
	err := dbOperations.InsertPanel(newPanels)
	if err != nil {
		return "", fmt.Errorf("error on inserting panel data", err)
	}
	log.Print("sucessfully inserted prom query into promquery-collection")

	newDashboard.PanelGroups = newPanelGroups

	err = dbOperations.InsertDashBoard(newDashboard)
	if err != nil {
		return "", fmt.Errorf("error on inserting panel data", err)
	}
	log.Print("sucessfully inserted dashboard into dashboard-collection")

	return "sucessfully created", nil
}

func UpdateDataSource(datasource model.DSInput) (*model.DSResponse, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	if datasource.DsID == nil || *datasource.DsID == "" {
		return nil, errors.New("Datasource ID is nil or empty")
	}

	query := bson.D{{"ds_id", datasource.DsID}}

	update := bson.D{{"$set", bson.D{{"ds_name", datasource.DsName},
		{"ds_url", datasource.DsURL}, {"access_type", datasource.AccessType},
		{"auth_type", datasource.AuthType}, {"basic_auth_username", datasource.BasicAuthUsername},
		{"basic_auth_password", datasource.BasicAuthPassword}, {"scrape_interval", datasource.ScrapeInterval},
		{"query_timeout", datasource.QueryTimeout}, {"http_method", datasource.HTTPMethod},
		{"updated_at", timestamp}}}}

	err := dbOperations.UpdateDataSource(query, update)
	if err != nil {
		return nil, err
	}

	return &model.DSResponse{
		DsID:              datasource.DsID,
		DsName:            &datasource.DsName,
		DsType:            &datasource.DsType,
		DsURL:             &datasource.DsURL,
		AccessType:        &datasource.AccessType,
		AuthType:          &datasource.AuthType,
		BasicAuthPassword: datasource.BasicAuthPassword,
		BasicAuthUsername: datasource.BasicAuthUsername,
		ScrapeInterval:    &datasource.ScrapeInterval,
		QueryTimeout:      &datasource.QueryTimeout,
		HTTPMethod:        &datasource.HTTPMethod,
		UpdatedAt:         &timestamp,
	}, nil
}

func UpdateDashBoard(dashboard *model.UpdataDBInput) (string, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	if dashboard.DbID == "" && dashboard.DsID == "" {
		return "", errors.New("DashBoard ID or Datasource is nil or empty")
	}

	query := bson.D{{"db_id", dashboard.DbID}}

	update := bson.D{{"$set", bson.D{{"ds_id", dashboard.DsID}, {"db_name", dashboard.DbName},
		{"db_type", dashboard.DbType}, {"end_time", dashboard.EndTime},
		{"start_time", dashboard.StartTime}, {"refresh_rate", dashboard.RefreshRate},
		{"panel_groups", dashboard.PanelGroups}, {"updated_at", timestamp}}}}

	err := dbOperations.UpdateDashboard(query, update)
	if err != nil {
		return "", err
	}

	return "successfully updated", nil
}

func UpdatePanel(panels []*model.Panel) (string, error) {

	for _, panel := range panels {
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)

		if *panel.PanelID == "" && *panel.PanelGroupID == "" {
			return "", errors.New("DashBoard ID or Datasource is nil or empty")
		}

		var newPanelOption dbSchema.PanelOption
		copier.Copy(&newPanelOption, &panel.PanelOptions)

		var newPromQueries []dbSchema.PromQuery
		copier.Copy(&newPromQueries, panel.PromQueries)

		query := bson.D{{"panel_id", panel.PanelID}}

		update := bson.D{{"$set", bson.D{{"panel_name", panel.PanelName},
			{"panel_group_id", panel.PanelGroupID}, {"panel_options", newPanelOption},
			{"prom_queries", newPromQueries}, {"updated_at", timestamp},
			{"y_axis_left", panel.YAxisLeft}, {"y_axis_right", panel.YAxisRight},
			{"x_axis_down", panel.XAxisDown}, {"unit", panel.Unit}}}}

		err := dbOperations.UpdatePanel(query, update)
		if err != nil {
			return "", err
		}
	}

	return "successfully updated", nil
}

func DeleteDashboard(db_id *string) (bool, error) {

	dashboardQuery := bson.M{"db_id": db_id, "is_removed": false}
	dashboard, err := dbOperations.GetDashboard(dashboardQuery)
	if err != nil {
		return false, fmt.Errorf("failed to list dashboard, error: %v", err)
	}

	for _, panelGroup := range dashboard.PanelGroups {
		listPanelQuery := bson.M{"panel_group_id": panelGroup.PanelGroupID, "is_removed": false}
		panels, err := dbOperations.ListPanel(listPanelQuery)
		if err != nil {
			return false, fmt.Errorf("failed to list Panel, error: %v", err)
		}

		for _, panel := range panels {
			time := strconv.FormatInt(time.Now().Unix(), 10)

			query := bson.D{{"panel_id", panel.PanelID}, {"is_removed", false}}
			update := bson.D{{"$set", bson.D{{"is_removed", true}, {"updated_at", time}}}}

			err := dbOperations.UpdatePanel(query, update)
			if err != nil {
				return false, fmt.Errorf("failed to delete panel, error: %v", err)
			}
		}
	}

	time := strconv.FormatInt(time.Now().Unix(), 10)

	query := bson.D{{"db_id", db_id}}
	update := bson.D{{"$set", bson.D{{"is_removed", true}, {"updated_at", time}}}}

	err = dbOperations.UpdateDashboard(query, update)
	if err != nil {
		return false, fmt.Errorf("failed to delete dashboard, error: %v", err)
	}

	return true, nil
}

func DeleteDataSource(input model.DeleteDSInput) (bool, error) {

	time := strconv.FormatInt(time.Now().Unix(), 10)

	listDBQuery := bson.M{"ds_id": input.DsID, "is_removed": false}
	dashboards, err := dbOperations.ListDashboard(listDBQuery)
	if err != nil {
		return false, fmt.Errorf("failed to list dashboard, error: %v", err)
	}

	if input.ForceDelete == true {
		for _, dashboard := range dashboards {

			for _, panelGroup := range dashboard.PanelGroups {
				listPanelQuery := bson.M{"panel_group_id": panelGroup.PanelGroupID, "is_removed": false}
				panels, err := dbOperations.ListPanel(listPanelQuery)
				if err != nil {
					return false, fmt.Errorf("failed to list Panel, error: %v", err)
				}

				for _, panel := range panels {
					query := bson.D{{"panel_id", panel.PanelID}, {"is_removed", false}}
					update := bson.D{{"$set", bson.D{{"is_removed", true}, {"updated_at", time}}}}

					err := dbOperations.UpdatePanel(query, update)
					if err != nil {
						return false, fmt.Errorf("failed to delete panel, error: %v", err)
					}
				}
			}
			updateDBQuery := bson.D{{"db_id", dashboard.DbID}}
			update := bson.D{{"$set", bson.D{{"is_removed", true}, {"updated_at", time}}}}

			err = dbOperations.UpdateDashboard(updateDBQuery, update)
			if err != nil {
				return false, fmt.Errorf("failed to delete dashboard, error: %v", err)
			}
		}

	} else if len(dashboards) > 0 {
		var db_names []string
		for _, dashboard := range dashboards {
			db_names = append(db_names, dashboard.DbName)
		}

		return false, fmt.Errorf("failed to delete datasource, dashboard(s) are attached to the datasource: %v", db_names)
	}

	updateDSQuery := bson.D{{"ds_id", input.DsID}}
	update := bson.D{{"$set", bson.D{{"is_removed", true}, {"updated_at", time}}}}

	err = dbOperations.UpdateDataSource(updateDSQuery, update)
	if err != nil {
		return false, fmt.Errorf("failed to delete datasource, error: %v", err)
	}

	return true, nil
}

func QueryListDataSource(projectID string) ([]*model.DSResponse, error) {
	query := bson.M{"project_id": projectID, "is_removed": false}

	datasource, err := dbOperations.ListDataSource(query)
	if err != nil {
		return nil, err
	}

	var newDatasources []*model.DSResponse
	copier.Copy(&newDatasources, &datasource)

	for _, datasource := range newDatasources {
		datasource.HealthStatus = prometheus.TSDBHealthCheck(*datasource.DsURL, *datasource.DsType)
	}

	return newDatasources, nil
}

func GetPromQuery(promInput *model.PromInput) ([]*model.PromResponse, error) {
	var newPromResponse []*model.PromResponse
	for _, v := range promInput.Queries {
		newPromQuery := analytics.PromQuery{
			Queryid:    v.Queryid,
			Query:      v.Query,
			Legend:     v.Legend,
			Resolution: v.Resolution,
			Minstep:    v.Minstep,
			URL:        promInput.URL,
			Start:      promInput.Start,
			End:        promInput.End,
		}

		response, err := prometheus.Query(newPromQuery)
		if err != nil {
			return nil, err
		}
		newPromResponse = append(newPromResponse, &response)
	}
	return newPromResponse, nil
}

func QueryListDashboard(projectID string) ([]*model.ListDashboardReponse, error) {
	query := bson.M{"project_id": projectID, "is_removed": false}

	dashboards, err := dbOperations.ListDashboard(query)
	if err != nil {
		return nil, fmt.Errorf("error on query from dashboard collection by projectid", err)
	}

	var newListDashboard []*model.ListDashboardReponse
	copier.Copy(&newListDashboard, &dashboards)

	for _, dashboard := range newListDashboard {
		datasource, err := dbOperations.GetDataSourceByID(dashboard.DsID)
		if err != nil {
			return nil, fmt.Errorf("error on querying from datasource collection", err)
		}

		dashboard.DsType = &datasource.DsType
		dashboard.DsName = &datasource.DsName

		cluster, err := dbOperations.GetCluster(dashboard.ClusterID)
		if err != nil {
			return nil, fmt.Errorf("error on querying from cluster collection", err)
		}

		dashboard.ClusterName = &cluster.ClusterName

		for _, panelgroup := range dashboard.PanelGroups {
			query := bson.M{"panel_group_id": panelgroup.PanelGroupID, "is_removed": false}
			panels, err := dbOperations.ListPanel(query)
			if err != nil {
				return nil, fmt.Errorf("error on querying from promquery collection", err)
			}

			var tempPanels []*model.PanelResponse
			copier.Copy(&tempPanels, &panels)

			panelgroup.Panels = tempPanels
		}
	}

	return newListDashboard, nil
}
