package globals

import (
	pb "github.com/PretendoNetwork/grpc/go/account"
	"github.com/PretendoNetwork/nex-go/v2"
	datastorecommon "github.com/PretendoNetwork/nex-protocols-common-go/v2/datastore"
	"github.com/PretendoNetwork/plogger-go"
	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const HPP_ACCESS_KEY = "76f26496"
const S3_KEY_DATASTORE = "ds"

var MinIOClient *minio.Client

var LibraryVersion = nex.NewLibraryVersion(3, 8, 3)
var Logger *plogger.Logger
var KerberosPassword = "" // randomized in init()
var HppServer *nex.HPPServer
var DatastoreCommon *datastorecommon.CommonProtocol
var GRPCAccountClientConnection *grpc.ClientConn
var GRPCAccountClient pb.AccountClient
var GRPCAccountCommonMetadata metadata.MD
var S3BucketName string
