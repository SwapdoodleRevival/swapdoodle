package nex

import (
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastorecommon "github.com/PretendoNetwork/nex-protocols-common-go/v2/datastore"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	securecommon "github.com/PretendoNetwork/nex-protocols-common-go/v2/secure-connection"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
	secure "github.com/PretendoNetwork/nex-protocols-go/v2/secure-connection"
	"github.com/silver-volt4/swapdoodle/database"
	datastore_db "github.com/silver-volt4/swapdoodle/database/datastore"
	"github.com/silver-volt4/swapdoodle/globals"
)

func registerProtocols() {
	secureProtocol := secure.NewProtocol()
	globals.HppServer.RegisterServiceProtocol(secureProtocol)

	commonSecureProtocol := securecommon.NewCommonProtocol(secureProtocol)
	commonSecureProtocol.CreateReportDBRecord = func(pid types.PID, reportID types.UInt32, reportData types.QBuffer) error {
		return nil
	}

	// Register DataStore protocol
	datastore := datastore.NewProtocol()
	globals.HppServer.RegisterServiceProtocol(datastore)

	// Register Common DataStore protocol
	commonDataStoreProtocol := datastorecommon.NewCommonProtocol(datastore)
	globals.DatastoreCommon = commonDataStoreProtocol

	dsm := common_globals.NewDataStoreManager(
		nil,
		database.Postgres,
	)

	manager := globals.NewS3Manager(globals.MinIOClient)
	dsm.SetS3Config(globals.S3BucketName, globals.S3_KEY_DATASTORE, manager)
	dsm.VerifyObjectAccessPermission = datastore_db.VerifyReadAccessByDataIdAndPID

	globals.DatastoreCommon.SetManager(dsm)
}
