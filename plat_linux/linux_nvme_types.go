package plat_linux

// Reference: https://www.naraeon.net/windows-10-nvme-get-features-sample-cpp/
// Reference: 10.0.16299.0/shared/nvme.h

type STORAGE_PROPERTY_ID = uint32
type STORAGE_QUERY_TYPE = uint32
type STORAGE_PROTOCOL_TYPE = uint32

const (
	StorageDeviceProperty STORAGE_PROPERTY_ID = 0 + iota
	StorageAdapterProperty
	StorageDeviceIdProperty
	StorageDeviceUniqueIdProperty
	StorageDeviceWriteCacheProperty
	StorageMiniportProperty
	StorageAccessAlignmentProperty
	StorageDeviceSeekPenaltyProperty
	StorageDeviceTrimProperty
	StorageDeviceWriteAggregationProperty
	StorageDeviceDeviceTelemetryProperty
	StorageDeviceLBProvisioningProperty
	StorageDevicePowerProperty
	StorageDeviceCopyOffloadProperty
	StorageDeviceResiliencyProperty
	StorageDeviceMediumProductType
	StorageAdapterRpmbProperty
	StorageAdapterCryptoProperty
)

const (
	StorageDeviceIoCapabilityProperty STORAGE_PROPERTY_ID = 48 + iota
	StorageAdapterProtocolSpecificProperty
	StorageDeviceProtocolSpecificProperty
	StorageAdapterTemperatureProperty
	StorageDeviceTemperatureProperty
	StorageAdapterPhysicalTopologyProperty
	StorageDevicePhysicalTopologyProperty
	StorageDeviceAttributesProperty
	StorageDeviceManagementStatus
	StorageAdapterSerialNumberProperty
	StorageDeviceLocationProperty
	StorageDeviceNumaProperty
	StorageDeviceZonedDeviceProperty
	StorageDeviceUnsafeShutdownCount
	StorageDeviceEnduranceProperty
	StorageDeviceLedStateProperty
)

const (
	StorageDeviceSelfEncryptionProperty STORAGE_PROPERTY_ID = 64 + iota
	StorageFruIdProperty
)

const (
	PropertyStandardQuery   STORAGE_QUERY_TYPE = 0
	PropertyExistsQuery     STORAGE_QUERY_TYPE = 1
	PropertyMaskQuery       STORAGE_QUERY_TYPE = 2
	PropertyQueryMaxDefined STORAGE_QUERY_TYPE = 3
)

const (
	ProtocolTypeUnknown STORAGE_PROTOCOL_TYPE = 0x00 + iota
	ProtocolTypeScsi
	ProtocolTypeAta
	ProtocolTypeNvme
	ProtocolTypeSd
	ProtocolTypeUfs
	ProtocolTypeProprietary STORAGE_PROTOCOL_TYPE = 0x7E
	ProtocolTypeMaxReserved STORAGE_PROTOCOL_TYPE = 0x7F
)

const (
	NVMeDataTypeUnknown = 0 + iota
	NVMeDataTypeIdentify
	NVMeDataTypeLogPage
	NVMeDataTypeFeature
)

const (
	NVME_IDENTIFY_CNS_SPECIFIC_NAMESPACE = 0 + iota
	NVME_IDENTIFY_CNS_CONTROLLER
	NVME_IDENTIFY_CNS_ACTIVE_NAMESPACES
	NVME_IDENTIFY_CNS_DESCRIPTOR_NAMESPACE
	NVME_IDENTIFY_CNS_NVM_SET
	NVME_IDENTIFY_CNS_SPECIFIC_NAMESPACE_IO_COMMAND_SET
	NVME_IDENTIFY_CNS_SPECIFIC_CONTROLLER_IO_COMMAND_SET
	NVME_IDENTIFY_CNS_ACTIVE_NAMESPACE_LIST_IO_COMMAND_SET
	NVME_IDENTIFY_CNS_ALLOCATED_NAMESPACE_LIST
	NVME_IDENTIFY_CNS_ALLOCATED_NAMESPACE
	NVME_IDENTIFY_CNS_CONTROLLER_LIST_OF_NSID
	NVME_IDENTIFY_CNS_CONTROLLER_LIST_OF_NVM_SUBSYSTEM
	NVME_IDENTIFY_CNS_PRIMARY_CONTROLLER_CAPABILITIES
	NVME_IDENTIFY_CNS_SECONDARY_CONTROLLER_LIST
	NVME_IDENTIFY_CNS_NAMESPACE_GRANULARITY_LIST
	NVME_IDENTIFY_CNS_UUID_LIST
	NVME_IDENTIFY_CNS_DOMAIN_LIST
	NVME_IDENTIFY_CNS_ENDURANCE_GROUP_LIST
	NVME_IDENTIFY_CNS_ALLOCATED_NAMSPACE_LIST_IO_COMMAND_SET
	NVME_IDENTIFY_CNS_ALLOCATED_NAMESPACE_IO_COMMAND_SET
	NVME_IDENTIFY_CNS_IO_COMMAND_SET
)

type STORAGE_PROPERTY_QUERY struct {
	PropertyId STORAGE_PROPERTY_ID
	QueryType  STORAGE_QUERY_TYPE
	// AdditionParameters
}

type STORAGE_PROTOCOL_SPECIFIC_DATA struct {
	ProtocolType                 STORAGE_PROTOCOL_TYPE
	DataType                     uint32
	ProtocolDataRequestValue     uint32
	ProtocolDataRequestSubValue  uint32
	ProtocolDataOffset           uint32
	ProtocolDataLength           uint32
	FixedProtocolReturnData      uint32
	ProtocolDataRequestSubValue2 uint32
	ProtocolDataRequestSubValue3 uint32
	ProtocolDataRequestSubValue4 uint32
}

type StorageQueryWithBuffer struct {
	Query            STORAGE_PROPERTY_QUERY
	ProtocolSpecific STORAGE_PROTOCOL_SPECIFIC_DATA
	Buffer           [4096]byte
}
