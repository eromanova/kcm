// Copyright 2024
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clusterdeployment

const (
	// Common
	EnvVarClusterDeploymentName     = "CLUSTER_DEPLOYMENT_NAME"
	EnvVarClusterDeploymentPrefix   = "CLUSTER_DEPLOYMENT_PREFIX"
	EnvVarClusterDeploymentTemplate = "CLUSTER_DEPLOYMENT_TEMPLATE"
	EnvVarControlPlaneNumber        = "CONTROL_PLANE_NUMBER"
	EnvVarWorkerNumber              = "WORKER_NUMBER"
	EnvVarNamespace                 = "NAMESPACE"
	// EnvVarNoCleanup disables After* cleanup in provider specs to allow for
	// debugging of test failures.
	EnvVarNoCleanup             = "NO_CLEANUP"
	EnvVarManagementClusterName = "MANAGEMENT_CLUSTER_NAME"

	// AWS
	EnvVarAWSAccessKeyID     = "AWS_ACCESS_KEY_ID"
	EnvVarAWSSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	EnvVarAWSVPCID           = "AWS_VPC_ID"
	EnvVarAWSSubnets         = "AWS_SUBNETS"
	EnvVarAWSInstanceType    = "AWS_INSTANCE_TYPE"
	EnvVarAWSSecurityGroupID = "AWS_SG_ID"
	EnvVarPublicIP           = "AWS_PUBLIC_IP"

	// VSphere
	EnvVarVSphereUser     = "VSPHERE_USER"
	EnvVarVSpherePassword = "VSPHERE_PASSWORD"

	// Azure
	EnvVarAzureClientSecret = "AZURE_CLIENT_SECRET"
	EnvVarAzureClientID     = "AZURE_CLIENT_ID"
	EnvVarAzureTenantID     = "AZURE_TENANT_ID"
	EnvVarAzureSubscription = "AZURE_SUBSCRIPTION"
	EnvVarAzureRegion       = "AZURE_REGION"

	// Adopted
	EnvVarAdoptedKubeconfigPath = "KUBECONFIG_DATA_PATH"
	EnvVarAdoptedCredential     = "ADOPTED_CREDENTIAL"

	// Remote
	EnvVarPrivateSSHKeyB64 = "PRIVATE_SSH_KEY_B64"
)
