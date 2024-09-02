package metadata

import (
	dtwebhook "github.com/Dynatrace/dynatrace-operator/pkg/webhook"
	corev1 "k8s.io/api/core/v1"
)

func (mut *Mutator) mutateUserContainers(request *dtwebhook.BaseRequest) {
	newContainers := request.NewContainers(mut.IsContainerInjected)
	for i := range newContainers {
		container := newContainers[i]
		setupVolumeMountsForUserContainer(container)
	}
}

func (mut *Mutator) reinvokeUserContainers(request *dtwebhook.BaseRequest) bool {
	var updated bool

	newContainers := request.NewContainers(mut.IsContainerInjected)

	if len(newContainers) == 0 {
		return false
	}

	for i := range newContainers {
		container := newContainers[i]
		setupVolumeMountsForUserContainer(container)

		updated = true
	}

	return updated
}

func updateInstallContainer(installContainer *corev1.Container, workload *workloadInfo, clusterName string) {
	addInjectedEnv(installContainer)
	addClusterNameEnv(installContainer, clusterName)
	addWorkloadInfoEnvs(installContainer, workload)
	addWorkloadEnrichmentInstallVolumeMount(installContainer)
}
