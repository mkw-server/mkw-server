package talker

import (
	"encoding/binary"
	"errors"
	"io"
	"net"

	"mkw-server/core"
	"mkw-server/logging"
	"mkw-server/util"
)

// WFCTalker talks over TCP to wfc-server, mainly handling adding/removing players from the room
// its messy currently since it does a lot and the messages aren't the best (not finalized)
type WFCTalker struct {
	conn net.Conn // TCP connection to wfc-server
	port uint16   // room's udp port, used as an id (a bad one)
}

var wfcTalker *WFCTalker

type RequestToMKWServer uint8

const (
	AddPlayer    = 0x01 // Request from wfc-server to add a player
	RemovePlayer = 0x02 // Request from wfc-server to remove a player
)

func (reqType RequestToMKWServer) valid() bool {
	return reqType == AddPlayer || reqType == RemovePlayer
}

type ResponseFromMKWServer uint8

const (
	Ready        = 0x00 // Response to wfc-server informing room is ready for players
	JoinAccepted = 0x01 // Response to wfc-server confirming a successful AddPlayer
	Log          = 0xff // Informs wfc-server of a mkw-server log
)

func NewWFCTalker(port uint16, serverAddress string) error {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		logging.Log("Failed to connect to WFC server: %v", err)
		return err
	}

	if !core.RoomInitialized() {
		logging.Log("Room pointer is nil when creating WFC Talker")
		return err
	}

	wfcTalker = &WFCTalker{
		conn: conn,
		port: port,
	}
	logging.Log("Talker created! (Hello wfc-server!)")

	// immediately tell wfc-server the room address
	sendReady()

	return nil
}

func sendReady() {
	buf := make([]byte, 3)
	buf[0] = Ready
	binary.BigEndian.PutUint16(buf[1:], wfcTalker.port)
	logging.Log("Sending ResponseType.Ready to wfc-server")
	SendToWFC(buf)
}

func Start() {
	go func() {
		header := make([]byte, 1)
		for {
			_, err := io.ReadFull(wfcTalker.conn, header)
			if err != nil {
				logging.Log("Error reading header from wfc-server %v", err)
				return
			}
			reqType := RequestToMKWServer(header[0])
			if !reqType.valid() {
				logging.Log("Invalid request type sent over %d", reqType)
				continue
			}
			handleRequest(reqType)
		}
	}()
}

func handleRequest(reqType RequestToMKWServer) {
	// First byte indicates request type
	switch reqType {
	case AddPlayer:
		logging.Log("Received AddPlayer")
		msg := make([]byte, AddPlayerRequestLength)
		_, err := io.ReadFull(wfcTalker.conn, msg)
		if err != nil {
			logging.Log("Error reading AddPlayer data from wfc-server: %v", err)
			return
		}

		req, err := unpackAddPlayerRequest(msg)
		if err != nil {
			logging.Log(err.Error())
			return
		}

		err = handleAddPlayerRequest(req)
		if err != nil {
			logging.Log("Failed to handle AddPlayer for reason %s:", err.Error())
		}
		logging.Log("Successfully handled AddPlayer for player %s", util.FormatIPPort(req.ip, req.port))

	case RemovePlayer:
		logging.Log("Received RemovePlayer")
		msg := make([]byte, RemovePlayerRequestLength)
		_, err := io.ReadFull(wfcTalker.conn, msg)
		if err != nil {
			logging.Log("Error reading RemovePlayer data from wfc-server: %v", err)
			return
		}

		req, err := unpackRemovePlayerRequest(msg)
		if err != nil {
			logging.Log("Unable to unpack RemovePlayer for reason %s", err.Error())
			return
		}

		err = handleRemovePlayerRequest(req)
		if err != nil {
			logging.Log("Failed to handle RemovePlayer for reason %s:", err.Error())
		}
		logging.Log("Successfully handled RemovePlayer for player %s", util.FormatIPPort(req.ip, req.port))

	default:
		logging.Log("Received Unknown Request %d from wfc-server", reqType)
	}
}

func SendToWFC(msg []byte) error {
	// Messages need to be wrapped around delimiters to prevent
	// messages from being streamed together.

	// messages are prefixed with 0xbb, 0xef, 0xdc, 0xc8
	// and suffixed with 0xce, 0xf9, 0xd3, 0xaa

	out := append([]byte{0xbb, 0xef, 0xdc, 0xc8}, msg...)
	out = append(out, []byte{0xce, 0xf9, 0xd3, 0xaa}...)
	_, err := wfcTalker.conn.Write(out)
	return err
}

func SendMessageToWFC(message string) error {
	if wfcTalker == nil {
		return errors.New("Can't send to wfc-server, talker is nil")
	}

	if wfcTalker.conn == nil {
		return errors.New("Can't send to wfc-server, talker.conn is nil")
	}

	// wfc-server knows to log mkw-server logs with Log (0xff)
	data := append([]byte{Log}, message...)

	err := SendToWFC(data)
	if err != nil {
		return err
	}

	return nil
}

func SendPacketDataToWFC(data []byte) error {
	_, err := wfcTalker.conn.Write(data)
	if err != nil {
		logging.Log("Failed to send packet data to WFC: %v", err)
		return err
	}
	logging.Log("Sent %d bytes of packet data to WFC", len(data))
	return nil
}

func Close() {
	wfcTalker.conn.Close()
}

func TalkerInitialized() bool {
	return wfcTalker != nil
}
