package nex_datastore_swapdoodle

import (
	"fmt"
	"strings"
	"time"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/silver-volt4/swapdoodle/globals"
)

func GetNotificationURL(err error, packet nex.PacketInterface, callID uint32, param datastore_types.DataStoreGetNotificationURLParam) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	key := globals.S3GetNotificationKey(packet.Sender().PID())

	get, err := globals.S3.PresignGet(key, time.Hour*24*7)
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.OperationNotAllowed, "change_error")
	}

	resStream := nex.NewByteStreamOut(globals.HppServer.LibraryVersions(), globals.HppServer.ByteStreamSettings())

	urlInfo := datastore_types.NewDataStoreReqGetNotificationURLInfo()

	urlInfo.URL = types.NewString(fmt.Sprintf("%s://%s/", get.URL.Scheme, get.URL.Host))
	urlInfo.Key = types.NewString(strings.TrimPrefix(get.URL.Path, "/"))
	urlInfo.Query = types.NewString(fmt.Sprintf("?%s", get.URL.Query().Encode()))
	urlInfo.RootCACert = types.NewBuffer(nil)

	urlInfo.WriteTo(resStream)

	res := nex.NewRMCSuccess(globals.HppServer, resStream.Bytes())
	res.ProtocolID = datastore.ProtocolID
	res.MethodID = datastore.MethodGetNotificationURL
	res.CallID = callID

	return res, nil
}
