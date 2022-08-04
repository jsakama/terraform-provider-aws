package rds

import "time"

const (
	ClusterRoleStatusActive  = "ACTIVE"
	ClusterRoleStatusDeleted = "DELETED"
	ClusterRoleStatusPending = "PENDING"
)

const (
	storageTypeStandard = "standard"
	storageTypeGP2      = "gp2"
	storageTypeIO1      = "io1"
)

func StorageType_Values() []string {
	return []string{
		storageTypeStandard,
		storageTypeGP2,
		storageTypeIO1,
	}
}

// https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/accessing-monitoring.html#Overview.DBInstance.Status.
const (
	InstanceStatusAvailable                     = "available"
	InstanceStatusBackingUp                     = "backing-up"
	InstanceStatusConfiguringEnhancedMonitoring = "configuring-enhanced-monitoring"
	InstanceStatusConfiguringIAMDatabaseAuth    = "configuring-iam-database-auth"
	InstanceStatusConfiguringLogExports         = "configuring-log-exports"
	InstanceStatusCreating                      = "creating"
	InstanceStatusDeleting                      = "deleting"
	InstanceStatusIncompatibleParameters        = "incompatible-parameters"
	InstanceStatusIncompatibleRestore           = "incompatible-restore"
	InstanceStatusMaintenance                   = "maintenance"
	InstanceStatusModifying                     = "modifying"
	InstanceStatusRebooting                     = "rebooting"
	InstanceStatusRenaming                      = "renaming"
	InstanceStatusResettingMasterCredentials    = "resetting-master-credentials"
	InstanceStatusStarting                      = "starting"
	InstanceStatusStopping                      = "stopping"
	InstanceStatusStorageFull                   = "storage-full"
	InstanceStatusStorageOptimization           = "storage-optimization"
	InstanceStatusUpgrading                     = "upgrading"
)

const (
	InstanceAutomatedBackupStatusPending     = "pending"
	InstanceAutomatedBackupStatusReplicating = "replicating"
	InstanceAutomatedBackupStatusRetained    = "retained"
)

const (
	EventSubscriptionStatusActive    = "active"
	EventSubscriptionStatusCreating  = "creating"
	EventSubscriptionStatusDeleting  = "deleting"
	EventSubscriptionStatusModifying = "modifying"
)

const (
	EngineAurora           = "aurora"
	EngineAuroraMySQL      = "aurora-mysql"
	EngineAuroraPostgreSQL = "aurora-postgresql"
	EngineMySQL            = "mysql"
	EnginePostgres         = "postgres"
)

func Engine_Values() []string {
	return []string{
		EngineAurora,
		EngineAuroraMySQL,
		EngineAuroraPostgreSQL,
		EngineMySQL,
		EnginePostgres,
	}
}

const (
	EngineModeGlobal        = "global"
	EngineModeMultiMaster   = "multimaster"
	EngineModeParallelQuery = "parallelquery"
	EngineModeProvisioned   = "provisioned"
	EngineModeServerless    = "serverless"
)

func EngineMode_Values() []string {
	return []string{
		EngineModeGlobal,
		EngineModeMultiMaster,
		EngineModeParallelQuery,
		EngineModeProvisioned,
		EngineModeServerless,
	}
}

const (
	ExportableLogTypeAgent      = "agent"
	ExportableLogTypeAlert      = "alert"
	ExportableLogTypeAudit      = "audit"
	ExportableLogTypeError      = "error"
	ExportableLogTypeGeneral    = "general"
	ExportableLogTypeListener   = "listener"
	ExportableLogTypeOEMAgent   = "oemagent"
	ExportableLogTypePostgreSQL = "postgresql"
	ExportableLogTypeSlowQuery  = "slowquery"
	ExportableLogTypeTrace      = "trace"
	ExportableLogTypeUpgrade    = "upgrade"
)

func ClusterExportableLogType_Values() []string {
	return []string{
		ExportableLogTypeAudit,
		ExportableLogTypeError,
		ExportableLogTypeGeneral,
		ExportableLogTypePostgreSQL,
		ExportableLogTypeSlowQuery,
	}
}

func InstanceExportableLogType_Values() []string {
	return []string{
		ExportableLogTypeAgent,
		ExportableLogTypeAlert,
		ExportableLogTypeAudit,
		ExportableLogTypeError,
		ExportableLogTypeGeneral,
		ExportableLogTypeListener,
		ExportableLogTypeOEMAgent,
		ExportableLogTypePostgreSQL,
		ExportableLogTypeSlowQuery,
		ExportableLogTypeTrace,
		ExportableLogTypeUpgrade,
	}
}

const (
	RestoreTypeCopyOnWrite = "copy-on-write"
	RestoreTypeFullCopy    = "full-copy"
)

func RestoreType_Values() []string {
	return []string{
		RestoreTypeCopyOnWrite,
		RestoreTypeFullCopy,
	}
}

const (
	TimeoutActionForceApplyCapacityChange = "ForceApplyCapacityChange"
	TimeoutActionRollbackCapacityChange   = "RollbackCapacityChange"
)

func TimeoutAction_Values() []string {
	return []string{
		TimeoutActionForceApplyCapacityChange,
		TimeoutActionRollbackCapacityChange,
	}
}

const (
	propagationTimeout = 2 * time.Minute
)
