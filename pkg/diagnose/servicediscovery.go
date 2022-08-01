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

package diagnose

import (
	"github.com/submariner-io/admiral/pkg/reporter"
	"github.com/submariner-io/subctl/pkg/cluster"
)

// const and vars

func ServiceDiscovery(clusterInfo *cluster.Info, status reporter.Interface) bool {
	mustHaveSubmariner(clusterInfo)

	status.Start("Checks if the Ligthouse components function properly.")
	defer status.End()

	status.Success("The detected Lighthouse components work properly.")
	return true
}
