package talker

import (
	"encoding/binary"
	"errors"
	"fmt"

	"mkw-server/core"
	"mkw-server/logging"
	"mkw-server/util"
)

// Id: RequestFromWFCServer.AddPlayer (0x01)
type AddPlayerRequest struct {
	ip       uint32
	port     uint16
	aid      uint8
	isHost   bool
	searchId uint64
	hasGuest bool
}

// TODO: Shouldn't need to update this everytime a field is added, it should be automatic.
const AddPlayerRequestLength = 17

// Id: ResponseToWFCServer.JoinAccepted (0x01)
type JoinAcceptedResponse struct {
	searchId uint64
}

func unpackAddPlayerRequest(msg []byte) (*AddPlayerRequest, error) {
	if len(msg) != AddPlayerRequestLength {
		return nil, fmt.Errorf("Unable to unpack AddPlayerMessage! len(msg) (%d) != %d", len(msg), AddPlayerRequestLength)
	}

	return &AddPlayerRequest{
		ip:       binary.BigEndian.Uint32(msg[0:4]),
		port:     binary.BigEndian.Uint16(msg[4:6]),
		aid:      uint8(msg[6]),
		isHost:   msg[7] == 1,
		searchId: binary.BigEndian.Uint64(msg[8:16]),
		hasGuest: msg[16] == 1,
	}, nil
}

func packJoinAcceptedResponse(searchId uint64) []byte {
	b := make([]byte, 9)

	b[0] = JoinAccepted
	binary.BigEndian.PutUint64(b[1:], searchId)
	return b
}

// addr is the address of the client that wants to join the room
func handleAddPlayerRequest(req *AddPlayerRequest) error {
	if req == nil {
		return errors.New("AddPlayerRequest is nil!")
	}

	logging.Log("Handling AddPlayer. Attempting to add player %s", util.FormatIPPort(req.ip, req.port))

	if !core.RoomInitialized() {
		return errors.New("Room isn't initialized, this should not happen at this point")
	}

	addr := util.CreateUDPAddr(req.ip, req.port)
	if addr == nil {
		return errors.New("CreateUDPAddr returned nil")
	}

	err := core.AddPlayerToRoom(addr.String(), req.aid, req.hasGuest)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	logging.Log("Sending wfc-server PlayerAdded")

	err = SendToWFC(packJoinAcceptedResponse(req.searchId))
	if err != nil {
		return fmt.Errorf("Failed to notify WFC of new player: %v", err)
	}
	return nil
}
