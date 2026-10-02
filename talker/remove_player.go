package talker

import (
	"encoding/binary"
	"errors"
	"fmt"

	"mkw-server/core"
	"mkw-server/logging"
	"mkw-server/util"
)

// Id: RemoveAid (0x02)
type RemoveAidRequest struct {
	ip   uint32
	port uint16
}

const RemoveAidRequestLength = 6

func unpackRemoveAidRequest(msg []byte) (*RemoveAidRequest, error) {
	if len(msg) != RemoveAidRequestLength {
		return nil, fmt.Errorf("RemoveAidRequest isn't 6 bytes (%d)", len(msg))
	}

	return &RemoveAidRequest{
		ip:   binary.BigEndian.Uint32(msg[0:4]),
		port: binary.BigEndian.Uint16(msg[4:6]),
	}, nil
}

func handleRemoveAidRequest(req *RemoveAidRequest) error {
	if req == nil {
		return errors.New("RemoveAidRequest is nil")
	}

	logging.Log("Handling RemoveAid. Attempting to remove aid %s", util.FormatIPPort(req.ip, req.port))

	if !core.RoomInitialized() {
		return errors.New("Room isn't initialized, this shouldn't happen at this point")
	}

	addr := util.CreateUDPAddr(req.ip, req.port)
	if addr == nil {
		return errors.New("CreateUDPAddr failed in handleRemoveAidRequest")
	}

	err := core.RemoveAidFromRoom(addr.String())
	if err != nil {
		return err
	}

	return nil
}
