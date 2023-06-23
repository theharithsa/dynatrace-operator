package version

import (
	"context"
	"testing"

	dynatracev1beta1 "github.com/Dynatrace/dynatrace-operator/src/api/v1beta1"
	"github.com/Dynatrace/dynatrace-operator/src/dtclient"
	"github.com/Dynatrace/dynatrace-operator/src/kubeobjects/address"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCodeModulesUpdater(t *testing.T) {
	testImage := dtclient.LatestImageInfo{
		Source: "some.registry.com",
		Tag:    "1.2.3",
	}
	t.Run("Getters work as expected", func(t *testing.T) {
		dynakube := &dynatracev1beta1.DynaKube{
			Spec: dynatracev1beta1.DynaKubeSpec{
				OneAgent: dynatracev1beta1.OneAgentSpec{
					ApplicationMonitoring: &dynatracev1beta1.ApplicationMonitoringSpec{
						Version:      testImage.Tag,
						UseCSIDriver: address.Of(true),
						AppInjectionSpec: dynatracev1beta1.AppInjectionSpec{
							CodeModulesImage: testImage.String(),
						},
					},
				},
			},
		}
		mockClient := &dtclient.MockDynatraceClient{}
		mockCodeModulesImageInfo(mockClient, testImage)
		updater := newCodeModulesUpdater(dynakube, mockClient)

		assert.Equal(t, "codemodules", updater.Name())
		assert.True(t, updater.IsEnabled())
		assert.Equal(t, dynakube.Spec.OneAgent.ApplicationMonitoring.CodeModulesImage, updater.CustomImage())
		assert.Equal(t, dynakube.Spec.OneAgent.ApplicationMonitoring.Version, updater.CustomVersion())
		assert.True(t, updater.IsAutoUpdateEnabled())
		imageInfo, err := updater.LatestImageInfo()
		require.NoError(t, err)
		assert.Equal(t, testImage, *imageInfo)
	})
}

func TestCodeModulesUseDefault(t *testing.T) {
	ctx := context.TODO()
	testVersion := "1.2.3"
	t.Run("Set according to version field, unset previous status", func(t *testing.T) {
		dynakube := &dynatracev1beta1.DynaKube{
			Spec: dynatracev1beta1.DynaKubeSpec{
				OneAgent: dynatracev1beta1.OneAgentSpec{
					ApplicationMonitoring: &dynatracev1beta1.ApplicationMonitoringSpec{
						Version: testVersion,
					},
				},
			},
			Status: dynatracev1beta1.EnvironmentStatus{
				CodeModules: oldCodeModulesStatus(),
			},
		}
		mockClient := &dtclient.MockDynatraceClient{}
		updater := newCodeModulesUpdater(dynakube, mockClient)

		err := updater.UseTenantRegistry(ctx, nil)
		require.NoError(t, err)
		assertDefaultCodeModulesStatus(t, testVersion, dynakube.Status.CodeModules)
	})
	t.Run("Set according to default, unset previous status", func(t *testing.T) {
		dynakube := &dynatracev1beta1.DynaKube{
			Spec: dynatracev1beta1.DynaKubeSpec{
				OneAgent: dynatracev1beta1.OneAgentSpec{
					ApplicationMonitoring: &dynatracev1beta1.ApplicationMonitoringSpec{},
				},
			},
			Status: dynatracev1beta1.EnvironmentStatus{
				CodeModules: oldCodeModulesStatus(),
			},
		}
		mockClient := &dtclient.MockDynatraceClient{}
		mockLatestAgentVersion(mockClient, testVersion)
		updater := newCodeModulesUpdater(dynakube, mockClient)

		err := updater.UseTenantRegistry(ctx, nil)
		require.NoError(t, err)
		assertDefaultCodeModulesStatus(t, testVersion, dynakube.Status.CodeModules)
	})
}

func oldCodeModulesStatus() dynatracev1beta1.CodeModulesStatus {
	return dynatracev1beta1.CodeModulesStatus{
		VersionStatus: dynatracev1beta1.VersionStatus{
			ImageID: "prev",
		},
	}
}

func assertDefaultCodeModulesStatus(t *testing.T, expectedVersion string, codeModulesStatus dynatracev1beta1.CodeModulesStatus) {
	assert.Equal(t, expectedVersion, codeModulesStatus.Version)
	assert.Empty(t, codeModulesStatus.ImageID)
}
