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
	nex_datastore_swapdoodle "github.com/silver-volt4/swapdoodle/nex/datastore/nex_datastore_swapdoodle"
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
	datastore.GetNotificationURL = nex_datastore_swapdoodle.GetNotificationURL
	datastore.GetNewArrivedNotificationsV1 = nex_datastore_swapdoodle.GetNewArrivedNotificationsV1
	globals.HppServer.RegisterServiceProtocol(datastore)

	// Register Common DataStore protocol
	commonDataStoreProtocol := datastorecommon.NewCommonProtocol(datastore)
	globals.DatastoreCommon = commonDataStoreProtocol

	dsm := common_globals.NewDataStoreManager(
		nil,
		database.Postgres,
	)
	dsm.SetS3Config(globals.S3.Bucket, globals.S3.KeyBase, globals.S3.Presigner)
	dsm.VerifyObjectAccessPermission = datastore_db.VerifyReadAccessByDataIdAndPID

	globals.DatastoreCommon.SetManager(dsm)
	globals.DatastoreCommon.OnAfterCompletePostObject = nex_datastore_swapdoodle.OnAfterCompletePostObject

	// The datastore DB schema is created immediately after SetManager is called.
	// Since notifications have not been implemented yet, I'm keeping the old DB around
	database.InitNotificationTable()
}
