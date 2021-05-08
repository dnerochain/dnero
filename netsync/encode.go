package netsync

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/dnerochain/dnero/common"
	"github.com/dnerochain/dnero/dispatcher"
	"github.com/dnerochain/dnero/rlp"
)

// type MessageIDEnum uint8

// const (
// 	MessageIDInvRequest MessageIDEnum = iota
// 	MessageIDInvResponse
// 	MessageIDDataRequest
// 	MessageIDDataResponse
// )

func encodeMessage(message interface{}) (common.Bytes, error) {
	var buf bytes.Buffer
	var msgID common.MessageIDEnum
	switch message.(type) {
	case dispatcher.InventoryRequest:
		msgID = common.MessageIDInvRequest
	case dispatcher.InventoryResponse:
		msgID = common.MessageIDInvResponse
	case dispatcher.DataRequest:
		msgID = common.MessageIDDataRequest
	case dispatcher.DataResponse:
		msgID = common.MessageIDDataResponse
	default:
		return nil, errors.New("Unsupported message type")
	}
	err := rlp.Encode(&buf, msgID)
	if err != nil {
		return nil, err
	}
	err = rlp.Encode(&buf, message)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decodeMessage(raw common.Bytes) (interface{}, error) {
	if len(raw) <= 1 {
		return nil, fmt.Errorf("Invalid message size")
	}
	var msgID common.MessageIDEnum
	err := rlp.DecodeBytes(raw[:1], &msgID)
	if err != nil {
		return nil, err
	}
	if msgID == common.MessageIDInvRequest {
		data := dispatcher.InventoryRequest{}
		err = rlp.DecodeBytes(raw[1:], &data)
		return data, err
	} else if msgID == common.MessageIDInvResponse {
		data := dispatcher.InventoryResponse{}
		err = rlp.DecodeBytes(raw[1:], &data)
		return data, err
	} else if msgID == common.MessageIDDataRequest {
		data := dispatcher.DataRequest{}
		err = rlp.DecodeBytes(raw[1:], &data)
		return data, err
	} else if msgID == common.MessageIDDataResponse {
		data := dispatcher.DataResponse{}
		err = rlp.DecodeBytes(raw[1:], &data)
		return data, err
	} else {
		return nil, fmt.Errorf("Unknown message ID: %v", msgID)
	}
}
