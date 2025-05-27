package datastore_db

import (
	"database/sql"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/silver-volt4/swapdoodle/database"
	"github.com/silver-volt4/swapdoodle/globals"
)

func VerifyReadAccessByDataIdAndPID(requesterPID types.PID, metaInfo datastore_types.DataStoreMetaInfo, objectAccessPassword, requesterAccessPassword types.UInt64) *nex.Error {
	nexError := IsObjectAvailable(metaInfo.DataID)
	if nexError != nil {
		return nexError
	}

	var t bool
	err := database.Postgres.QueryRow(`SELECT 1 FROM datastore.objects WHERE data_id=$1 AND (owner = $2 OR $2 IN (SELECT UNNEST(access_permission_recipients)))`, metaInfo.DataID, requesterPID).Scan(&t)

	if err != nil {
		if err == sql.ErrNoRows {
			return nex.NewError(nex.ResultCodes.DataStore.OperationNotAllowed, "Access denied")
		}

		globals.Logger.Error(err.Error())

		// TODO - Send more specific errors?
		return nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	return nil
}
