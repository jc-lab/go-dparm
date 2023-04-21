package ata

import (
	"github.com/jc-lab/go-dparm/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLbaRegsSize(t *testing.T) {
	assert.Equal(t, 5, test.SizeOf(t, &LbaRegs{}))
}

func TestTfSize(t *testing.T) {
	assert.Equal(t, 15, test.SizeOf(t, &Tf{}))
}

func TestIdentityGeneralConfigurationSize(t *testing.T) {
	assert.Equal(t, 2, test.SizeOf(t, &IdentityGeneralConfiguration{}))
}

func TestIdentityTrustedComputingSize(t *testing.T) {
	assert.Equal(t, 2, test.SizeOf(t, &IdentityTrustedComputing{}))
}

func TestIdentityCapabilitiesSize(t *testing.T) {
	assert.Equal(t, 4, test.SizeOf(t, &IdentityCapabilities{}))
}

func TestIdentityAdditionalSupportedSize(t *testing.T) {
	assert.Equal(t, 2, test.SizeOf(t, &IdentityAdditionalSupported{}))
}

func TestIdentitySerialAtaCapabilitiesSize(t *testing.T) {
	assert.Equal(t, 4, test.SizeOf(t, &IdentitySerialAtaCapabilities{}))
}

func TestIdentitySerialAtaFeaturesSupportedSize(t *testing.T) {
	assert.Equal(t, 2, test.SizeOf(t, &IdentitySerialAtaFeaturesSupported{}))
}

func TestIdentitySerialAtaFeaturesEnabledSize(t *testing.T) {
	assert.Equal(t, 2, test.SizeOf(t, &IdentitySerialAtaFeaturesEnabled{}))
}

func TestIdentityCommandSetSupportSize(t *testing.T) {
	assert.Equal(t, 6, test.SizeOf(t, &IdentityCommandSetSupport{}))
}

func TestIdentityCommandSetActiveSize(t *testing.T) {
	assert.Equal(t, 6, test.SizeOf(t, &IdentityCommandSetActive{}))
}

func TestIdentityNormalSecurityEraseUnitSize(t *testing.T) {
	assert.Equal(t, 2, test.SizeOf(t, &IdentityNormalSecurityEraseUnit{}))
}

func TestIdentityEnhancedSecurityEraseUnitSize(t *testing.T) {
	assert.Equal(t, 2, test.SizeOf(t, &IdentityEnhancedSecurityEraseUnit{}))
}

func TestIdentityPhysicalLogicalSectorSizeSize(t *testing.T) {
	assert.Equal(t, 2, test.SizeOf(t, &IdentityPhysicalLogicalSectorSize{}))
}

func TestIdentityCommandSetSupportExtSize(t *testing.T) {
	assert.Equal(t, 2, test.SizeOf(t, &IdentityCommandSetSupportExt{}))
}

func TestIdentityCommandSetActiveExtSize(t *testing.T) {
	assert.Equal(t, 2, test.SizeOf(t, &IdentityCommandSetActiveExt{}))
}

func TestIdentitySecurityStatusSize(t *testing.T) {
	assert.Equal(t, 2, test.SizeOf(t, &IdentitySecurityStatus{}))
}

func TestIdentityCfgPowerMode1Size(t *testing.T) {
	assert.Equal(t, 2, test.SizeOf(t, &IdentityCfgPowerMode1{}))
}

func TestIdentityDataSetManagementFeatureSize(t *testing.T) {
	assert.Equal(t, 2, test.SizeOf(t, &IdentityDataSetManagementFeature{}))
}

func TestIdentitySctSommandTransportSize(t *testing.T) {
	assert.Equal(t, 2, test.SizeOf(t, &IdentitySctSommandTransport{}))
}

func TestIdentityBlockAlignmentSize(t *testing.T) {
	assert.Equal(t, 2, test.SizeOf(t, &IdentityBlockAlignment{}))
}

func TestIdentityNvCacheCapabilitiesSize(t *testing.T) {
	assert.Equal(t, 2, test.SizeOf(t, &IdentityNvCacheCapabilities{}))
}

func TestIdentityNvCacheOptionsSize(t *testing.T) {
	assert.Equal(t, 2, test.SizeOf(t, &IdentityNvCacheOptions{}))
}

func TestIdentityTransportMajorVersionSize(t *testing.T) {
	assert.Equal(t, 2, test.SizeOf(t, &IdentityTransportMajorVersion{}))
}

func TestIdentityDeviceDataSize(t *testing.T) {
	assert.Equal(t, 512, test.SizeOf(t, &IdentityDeviceData{}))
}
