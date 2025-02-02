/*
SPDX-License-Identifier: Apache-2.0

Copyright Contributors to the Submariner project.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package scc

import (
	"github.com/pkg/errors"
	"github.com/submariner-io/subctl/pkg/scc"
	"github.com/submariner-io/submariner-operator/pkg/embeddedyamls"
	"k8s.io/client-go/dynamic"
)

func Ensure(dynClient dynamic.Interface, namespace string) (bool, error) {
	gatewaySaName, err := embeddedyamls.GetObjectName(embeddedyamls.Config_rbac_submariner_gateway_service_account_yaml)
	if err != nil {
		return false, errors.Wrap(err, "error parsing the gateway ServiceAccount resource")
	}

	routeAgentSaName, err := embeddedyamls.GetObjectName(embeddedyamls.Config_rbac_submariner_route_agent_service_account_yaml)
	if err != nil {
		return false, errors.Wrap(err, "error parsing the route agent ServiceAccount resource")
	}

	globalnetSaName, err := embeddedyamls.GetObjectName(embeddedyamls.Config_rbac_submariner_globalnet_service_account_yaml)
	if err != nil {
		return false, errors.Wrap(err, "error parsing the globalnet ServiceAccount resource")
	}

	npSyncerSaName, err := embeddedyamls.GetObjectName(embeddedyamls.Config_rbac_networkplugin_syncer_service_account_yaml)
	if err != nil {
		return false, errors.Wrap(err, "error parsing the networkplugin syncer ServiceAccount resource")
	}

	diagnoseSaName, err := embeddedyamls.GetObjectName(embeddedyamls.Config_rbac_submariner_diagnose_service_account_yaml)
	if err != nil {
		return false, errors.Wrap(err, "error parsing the diagnose ServiceAccount resource")
	}

	saNames := []string{
		gatewaySaName, routeAgentSaName, globalnetSaName, npSyncerSaName, diagnoseSaName,
	}
	updateScc := false

	for _, saName := range saNames {
		result, err := scc.Update(dynClient, namespace, saName)
		if err != nil {
			return false, errors.Wrap(err, "error updating the SCC resource")
		}

		updateScc = updateScc || result
	}

	return updateScc, nil
}
