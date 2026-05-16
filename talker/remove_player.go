package talker

import (
	"encoding/binary"
	"errors"
	"fmt"

	"mkw-server/core"
	"mkw-server/logging"
	"mkw-server/util"
)

// Id: RemovePlayer (0x02)
type RemovePlayerRequest struct {
	ip   uint32
	port uint16
}

func unpackRemovePlayerRequest(msg []byte) (*RemovePlayerRequest, error) {
	if len(msg) != 6 {
		return nil, fmt.Errorf("RemovePlayerRequest isn't 6 bytes (%d)", len(msg))
	}

	return &RemovePlayerRequest{
		ip:   binary.BigEndian.Uint32(msg[0:4]),
		port: binary.BigEndian.Uint16(msg[4:6]),
	}, nil
}

func handleRemovePlayerRequest(req *RemovePlayerRequest) error {
	if req == nil {
		return errors.New("LeaveMessage is nil")
	}

	logging.Log("Handling RemovePlayer. Attempting to remove player %s", util.FormatIPPort(req.ip, req.port))

	if !core.RoomInitialized() {
		return errors.New("Room isn't initialized, this shouldn't happen at this point")
	}

	addr := util.CreateUDPAddr(req.ip, req.port)
	if addr == nil {
		return errors.New("CreateUDPAddr failed in handleRemovePlayerRequest")
	}

	err := core.RemovePlayerFromRoom(addr.String())
	if err != nil {
		return err
	}

	return nil
}
