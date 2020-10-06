package types

//SubscriberConfigurationVars contains the required configurable parameters for subscriber installation
type SubscriberConfigurationVars struct {
	PortalNamespace         string
	PortalScope             string
	GQLServerURI            string
	SubscriberImage         string
	ArgoServerImage         string
	WorkflowControllerImage string
	ChaosOperatorImage      string
	WorkflowExecutorImage   string
	ChaosRunnerImage        string
}
