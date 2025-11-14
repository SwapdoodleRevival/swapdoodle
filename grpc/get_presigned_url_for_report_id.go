package grpc

import (
	"context"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/datastore/database"
	"github.com/silver-volt4/swapdoodle/globals"
)

const charset = "0123456789BCDFGHJKLMNPRTVWXY"

var key uint8

func (s *gRPCSwapdoodleServer) GetPresignedUrlForReportId(context context.Context, in *GetPresignedUrlForReportIdRequest) (*GetPresignedUrlForReportIdResponse, error) {
	noteId := in.GetNoteId()

	var id uint64 = 0

	for _, char := range noteId {
		index := strings.IndexRune(charset, char)
		if index == -1 {
			return nil, fmt.Errorf("invalid NoteID string")
		}
		id = id*28 + uint64(index)
	}

	id ^= 0xDEAD9ED5

	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, id)

	checksum := key
	for _, b := range bs {
		checksum ^= b
	}

	if checksum != 0 {
		return nil, fmt.Errorf("invalid NoteID checksum")
	}

	id &= 0xFF_00_FF_FF_FF
	id |= (id & 0xFF_00_00_00_00) >> 8
	id &= 0xFF_FF_FF_FF

	metaInfo, _, err := database.GetAccessObjectInfoByDataID(globals.DatastoreManager, types.UInt64(id))
	if err != nil {
		return nil, err
	}

	version, err := database.GetObjectLatestVersionNumber(globals.DatastoreManager, types.UInt64(id))
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("objects/%020d_%010d.bin", metaInfo.DataID, version)
	getData, e := globals.DatastoreManager.S3.PresignGet(key, time.Minute*15)
	if e != nil {
		return nil, e
	}

	return &GetPresignedUrlForReportIdResponse{
		Url: getData.URL.String(),
	}, nil
}

func init() {
	key = 0
	for _, char := range globals.HPP_ACCESS_KEY {
		key ^= uint8(char)
		key = (key << 4) | (key >> 4)
	}
}
