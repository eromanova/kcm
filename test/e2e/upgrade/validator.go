// Copyright 2025
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

package upgrade

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	hcv2 "github.com/fluxcd/helm-controller/api/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/util/node"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	kcmv1 "github.com/K0rdent/kcm/api/v1alpha1"
)

type validationFunc func(context.Context, crclient.Client, crclient.Client, string, string, string) error

type ClusterUpgrade struct {
	mgmtClient  crclient.Client
	childClient crclient.Client
	namespace   string
	name        string
	newTemplate string
	validator   Validator
}

type Validator struct {
	validationFuncs []validationFunc
}

func NewClusterUpgrade(mgmtClient, childClient crclient.Client, namespace, name, newTemplate string, validator Validator) ClusterUpgrade {
	return ClusterUpgrade{
		mgmtClient:  mgmtClient,
		childClient: childClient,
		namespace:   namespace,
		name:        name,
		newTemplate: newTemplate,
		validator:   validator,
	}
}

func NewDefaultClusterValidator() Validator {
	return Validator{validationFuncs: []validationFunc{
		validateHelmRelease,
	}}
}

func NewK0sClusterValidator() Validator {
	return Validator{validationFuncs: []validationFunc{
		validateHelmRelease,
		validateK0sVersion,
	}}
}

func (v *ClusterUpgrade) Validate(ctx context.Context) {
	Eventually(func() bool {
		for _, validate := range v.validator.validationFuncs {
			if err := validate(ctx, v.mgmtClient, v.childClient, v.namespace, v.name, v.newTemplate); err != nil {
				_, _ = fmt.Fprintf(GinkgoWriter, "[%s/%s] upgrade validation error: %v\n", v.namespace, v.name, err)
				return false
			}
		}
		return true
	}, 20*time.Minute, 20*time.Second).Should(BeTrue())
}

func validateHelmRelease(ctx context.Context, mgmtClient, _ crclient.Client, namespace, name, newTemplate string) error {
	hr := &hcv2.HelmRelease{}
	err := mgmtClient.Get(ctx, types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, hr)
	if err != nil {
		return fmt.Errorf("failed to get %s/%s HelmRelease: %v", namespace, name, err)
	}

	template := &kcmv1.ClusterTemplate{}
	if err := mgmtClient.Get(ctx, types.NamespacedName{
		Namespace: namespace,
		Name:      newTemplate,
	}, template); err != nil {
		return err
	}
	if hr.Spec.ChartRef.Name != template.Status.ChartRef.Name {
		return fmt.Errorf("waiting for chartName to be updated in %s/%s HelmRelease", namespace, name)
	}
	readyCondition := apimeta.FindStatusCondition(hr.GetConditions(), kcmv1.ReadyCondition)
	if readyCondition == nil {
		return fmt.Errorf("waiting for %s/%s HelmRelease to have Ready condition", namespace, name)
	}
	if readyCondition.ObservedGeneration != hr.Generation {
		return errors.New("waiting for status.observedGeneration to be updated")
	}
	if readyCondition.Status != metav1.ConditionTrue {
		return errors.New("waiting for Ready condition to have status: true")
	}
	if readyCondition.Reason != hcv2.UpgradeSucceededReason {
		return errors.New("waiting for Ready condition to have `UpgradeSucceeded` reason")
	}
	return nil
}

func validateK0sVersion(ctx context.Context, mgmtClient, childClient crclient.Client, namespace, name, newTemplate string) error {
	template := &kcmv1.ClusterTemplate{}
	if err := mgmtClient.Get(ctx, types.NamespacedName{
		Namespace: namespace,
		Name:      newTemplate,
	}, template); err != nil {
		return err
	}

	k0sVersion, err := getK0sVersion(template)
	if err != nil {
		return fmt.Errorf("couldn't get k0s version from the ClusterTemplate: %v", err)
	}

	nodes := &corev1.NodeList{}
	err = childClient.List(ctx, nodes)
	if err != nil {
		return fmt.Errorf("failed to list nodes of the %s/%s cluster: %v", namespace, name, err)
	}
	for _, n := range nodes.Items {
		if n.Status.NodeInfo.KubeletVersion != k0sVersion {
			return fmt.Errorf("the k0s version of the node %s is not yet updated: expected %s, current %s", n.Name, k0sVersion, n.Status.NodeInfo.KubeletVersion)
		}
		if !node.IsNodeReady(&n) {
			return fmt.Errorf("node %s is not yet Ready", n.Name)
		}
	}
	return nil
}

func getK0sVersion(template *kcmv1.ClusterTemplate) (string, error) {
	var data map[string]any
	if err := json.Unmarshal(template.Status.Config.Raw, &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal the config of the %s template: %v", template.Name, err)
	}

	if k0s, ok := data["k0s"].(map[string]any); ok {
		if version, ok := k0s["version"].(string); ok {
			return version, nil
		}
	}
	return "", fmt.Errorf("k0s version is not defined in the config of the %s template", template.Name)
}
