package main

import (
	"time"

	"github.com/deepch/vdk/av"
)

// ClientAdd Add New Client to Translations
func (obj *StorageST) ClientAdd(streamID string, mode int) (string, chan *av.Packet, chan *[]byte, error) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	streamTmp, ok := obj.Streams[streamID]
	if !ok {
		return "", nil, nil, ErrorStreamNotFound
	}
	//Generate UUID client
	cid, err := generateUUID()
	if err != nil {
		return "", nil, nil, err
	}
	chAV := make(chan *av.Packet, 2000)
	chRTP := make(chan *[]byte, 2000)

	streamTmp.clients[cid] = ClientST{mode: mode, outgoingAVPacket: chAV, outgoingRTPPacket: chRTP, signals: make(chan int, 100)}
	streamTmp.ack = time.Now()
	obj.Streams[streamID] = streamTmp
	return cid, chAV, chRTP, nil

}

// ClientDelete Delete Client
func (obj *StorageST) ClientDelete(streamID string, cid string) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if _, ok := obj.Streams[streamID]; ok {
		delete(obj.Streams[streamID].clients, cid)
	}
}

// ClientHas check is client ext
func (obj *StorageST) ClientHas(streamID string) bool {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	streamTmp, ok := obj.Streams[streamID]
	if !ok {
		return false
	}
	if time.Since(streamTmp.ack).Seconds() > 30 {
		return false
	}
	return true
}
