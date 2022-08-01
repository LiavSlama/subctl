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
	"context"
	"fmt"

	"github.com/submariner-io/admiral/pkg/reporter"
	"github.com/submariner-io/subctl/internal/constants"
	"github.com/submariner-io/subctl/internal/gvr"
	"github.com/submariner-io/subctl/pkg/cluster"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	mcsv1a1 "sigs.k8s.io/mcs-api/pkg/apis/v1alpha1"
)

func ServiceDiscovery(clusterInfo *cluster.Info, status reporter.Interface) bool {
	status.Start("Checking Lighthouse components configuration")
	defer status.End()

	tracker := reporter.NewTracker(status)


	

	checkServiceExport(clusterInfo, tracker)

	if tracker.HasFailures() {
		return false
	}

	status.Success("Lighthouse components are working properly")

	return true
}

// this function checks if all serviceExport have a matching serviceImport.
func checkServiceExport(clusterInfo *cluster.Info, status reporter.Interface) {
	foundAll := true
	serviceExportGVR := gvr.FromMetaGroupVersion(mcsv1a1.GroupVersion, "serviceexports")

	serviceExports, err := clusterInfo.ClientProducer.ForDynamic().Resource(serviceExportGVR).Namespace(corev1.NamespaceAll).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		status.Failure("Error listing ServiceExport resources: %v", err)
	}

	serviceImportsGVR := gvr.FromMetaGroupVersion(mcsv1a1.GroupVersion, "serviceimports")

	
	for i := range serviceExports.Items {
		ns := serviceExports.Items[i].GetNamespace()
		name := serviceExports.Items[i].GetName()
		ep := clusterInfo.ClientProducer.ForKubernetes().DiscoveryV1().EndpointSlices(ns)
		_, err = ep.Get(context.TODO(), fmt.Sprintf("%s-%s", name, clusterInfo.Submariner.Spec.ClusterID), metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				status.Failure("the endpointSlice for %v-%v serviceExport doesn't exist", ns, name)
			} else {
				status.Failure("Couldn't find an endPointSlice match for %v-%v serviceExport", ns, name)
			}
			
			foundAll = false
		}

		_, err := clusterInfo.ClientProducer.ForDynamic().Resource(serviceImportsGVR).
		Namespace(constants.OperatorNamespace).Get(context.TODO(),
		 fmt.Sprintf("%s-%s-%s", name, ns, clusterInfo.Submariner.Spec.ClusterID), metav1.GetOptions{})
		 
		if err != nil {
			if apierrors.IsNotFound(err) {
				status.Failure("serviceImport for %v-%v serviceExport doesn't exist", ns, name)
			} else {
				status.Failure("Couldn't find serviceImport match for %v-%v serviceExport", ns, name)
			}
			
			foundAll = false
		}
	}
	
	if foundAll {
		status.Success("Found all serviceExport's a matching serviceImport and endPointSlice in the cluster")
	}
}
