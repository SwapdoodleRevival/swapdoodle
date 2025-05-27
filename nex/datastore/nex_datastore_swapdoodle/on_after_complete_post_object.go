package nex_datastore_swapdoodle

import (
	"fmt"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/silver-volt4/swapdoodle/database"
	"github.com/silver-volt4/swapdoodle/globals"
)

func OnAfterCompletePostObject(packet nex.PacketInterface, param datastore_types.DataStoreCompletePostParam) {
	// Create notifications
	_, err := database.Postgres.Exec(`INSERT INTO datastore.notifications (data_id, recipient_pid)
		SELECT data_id, UNNEST(access_permission_recipients) as recipient_pid
		FROM datastore.objects
		WHERE data_id = $1`, param.DataID)

	if err != nil {
		globals.Logger.Error(err.Error())
		return
	}

	rows, err := database.Postgres.Query(`SELECT MAX(notification_id), recipient_pid
	FROM datastore.notifications
	WHERE recipient_pid IN (SELECT UNNEST(access_permission_recipients) FROM datastore.objects WHERE data_id = $1)
	GROUP BY recipient_pid`, param.DataID)

	if err != nil {
		globals.Logger.Error(err.Error())
		return
	}

	// Update "noti-files"
	for rows.Next() {
		var notificationId uint64
		var pid types.PID
		rows.Scan(&notificationId, &pid)

		bucket := globals.S3.Bucket
		key := fmt.Sprintf("%s/%s", globals.S3.KeyBase, globals.S3GetNotificationKey(pid))

		// TODO: Looks like we're not the only ones wondering about the meaning of that last number:
		// https://github.com/PretendoNetwork/pokemon-rumble-world-secure/blob/main/nex/datastore/complete_post_object_v1.go#L46
		// Let's leave it constant for now...
		globals.S3SetFileContent(bucket, key, fmt.Sprintf("%d,%d,%d", notificationId, pid, 1479103557))
	}
}
