// Package security provides license validation and security features
package security

// Premium feature constants
const (
	// Integration features
	FeatureDatabricksIntegration = "databricks_integration"
	FeatureDBTIntegration        = "dbt_integration"
	FeatureCloudStorageS3        = "cloud_storage_s3"
	FeatureCloudStorageAzure     = "cloud_storage_azure"
	FeatureCloudStorageGCP       = "cloud_storage_gcp"
	FeatureDataCatalog           = "data_catalog"
	FeatureWorkflowOrchestration = "workflow_orchestration"

	// License plans
	PlanCommunity  = "community"
	PlanStarter    = "starter"
	PlanPro        = "pro"
	PlanEnterprise = "enterprise"
)

// PlanFeatures maps license plans to their included features
var PlanFeatures = map[string][]string{
	PlanCommunity: {},
	PlanStarter: {
		FeatureDatabricksIntegration,
		FeatureCloudStorageS3,
	},
	PlanPro: {
		FeatureDatabricksIntegration,
		FeatureDBTIntegration,
		FeatureCloudStorageS3,
		FeatureCloudStorageAzure,
		FeatureCloudStorageGCP,
		FeatureDataCatalog,
		FeatureWorkflowOrchestration,
	},
	// Enterprise tier is contact-only through nessi.dev
	PlanEnterprise: {
		FeatureDatabricksIntegration,
		FeatureDBTIntegration,
		FeatureCloudStorageS3,
		FeatureCloudStorageAzure,
		FeatureCloudStorageGCP,
		FeatureDataCatalog,
		FeatureWorkflowOrchestration,
	},
}

// FeatureDescriptions provides human-readable descriptions for premium features
var FeatureDescriptions = map[string]string{
	FeatureDatabricksIntegration: "Integration with Databricks for data catalog and Delta Lake",
	FeatureDBTIntegration:        "Integration with dbt for data transformation workflows",
	FeatureCloudStorageS3:        "Amazon S3 cloud storage integration",
	FeatureCloudStorageAzure:     "Azure Blob Storage integration",
	FeatureCloudStorageGCP:       "Google Cloud Storage integration",
	FeatureDataCatalog:           "Data catalog publishing and integration",
	FeatureWorkflowOrchestration: "Workflow orchestration with Airflow, Prefect, and Dagster",
}

// PlanPricing provides pricing information for each plan
var PlanPricing = map[string]struct {
	Monthly float64
	Annual  float64
}{
	PlanCommunity: {
		Monthly: 0,
		Annual:  0,
	},
	PlanStarter: {
		Monthly: 49.99,
		Annual:  499.90, // 2 months free
	},
	PlanPro: {
		Monthly: 99.99,
		Annual:  999.90, // 2 months free
	},
	// Enterprise tier is contact-only through nessi.dev
	PlanEnterprise: {
		Monthly: 0, // Contact for pricing
		Annual:  0, // Contact for pricing
	},
}
