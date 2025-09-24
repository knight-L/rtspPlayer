package main

import (
	"time"

	"github.com/deepch/vdk/av"
)

// StreamChannelMake check stream exist
func (obj *StorageST) StreamChannelMake(val StreamST) StreamST {
	//make client's
	val.clients = make(map[string]ClientST)
	//make last ack
	val.ack = time.Now().Add(-255 * time.Hour)
	//make signals buffer chain
	val.signals = make(chan int, 100)
	return val
}

// StreamChannelUnlock unlock status to no lock
func (obj *StorageST) StreamChannelUnlock(streamID string) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if streamTmp, ok := obj.Streams[streamID]; ok {
		streamTmp.runLock = false
		obj.Streams[streamID] = streamTmp
	}
}

// StreamChannelControl get stream
func (obj *StorageST) StreamChannelControl(key string) (*StreamST, error) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if streamTmp, ok := obj.Streams[key]; ok {
		return &streamTmp, nil
	}
	return nil, ErrorStreamNotFound
}

// StreamChannelExist check stream exist
func (obj *StorageST) StreamChannelExist(streamID string) bool {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if streamTmp, ok := obj.Streams[streamID]; ok {
		streamTmp.ack = time.Now()
		obj.Streams[streamID] = streamTmp
		return ok
	}
	return false
}

// StreamChannelReload reload stream
func (obj *StorageST) StreamChannelReload(uuid string) error {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		tmp.signals <- SignalStreamRestart
		return nil
	}
	return ErrorStreamNotFound
}

// StreamInfo return stream info
func (obj *StorageST) StreamChannelInfo(uuid string) (*StreamST, error) {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		return &tmp, nil
	}
	return nil, ErrorStreamNotFound
}

// StreamChannelCodecs get stream codec storage or wait
func (obj *StorageST) StreamChannelCodecs(streamID string) ([]av.CodecData, error) {
	for i := 0; i < 100; i++ {
		ret, err := (func() ([]av.CodecData, error) {
			obj.mutex.RLock()
			defer obj.mutex.RUnlock()
			tmp, ok := obj.Streams[streamID]
			if !ok {
				return nil, ErrorStreamNotFound
			}
			return tmp.codecs, nil
		})()

		if ret != nil || err != nil {
			return ret, err
		}

		time.Sleep(50 * time.Millisecond)
	}
	return nil, ErrorStreamChannelCodecNotFound
}

// StreamChannelStatus change stream status
func (obj *StorageST) StreamChannelStatus(key string, val int) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[key]; ok {
		tmp.Status = val
		obj.Streams[key] = tmp
	}
}

// StreamChannelCast broadcast stream
func (obj *StorageST) StreamChannelCast(key string, val *av.Packet) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[key]; ok {
		if len(tmp.clients) > 0 {
			for _, i2 := range tmp.clients {
				if i2.mode == RTSP {
					continue
				}
				if len(i2.outgoingAVPacket) < 1000 {
					i2.outgoingAVPacket <- val
				} else if len(i2.signals) < 10 {
					i2.signals <- SignalStreamStop
				}
			}
			tmp.ack = time.Now()
			obj.Streams[key] = tmp
		}
	}
}

// StreamChannelCastProxy broadcast stream
func (obj *StorageST) StreamChannelCastProxy(key string, val *[]byte) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[key]; ok {
		if len(tmp.clients) > 0 {
			for _, i2 := range tmp.clients {
				if i2.mode != RTSP {
					continue
				}
				if len(i2.outgoingRTPPacket) < 1000 {
					i2.outgoingRTPPacket <- val
				} else if len(i2.signals) < 10 {
					i2.signals <- SignalStreamStop
				}
			}
			tmp.ack = time.Now()
			obj.Streams[key] = tmp
		}
	}
}

// StreamChannelCodecsUpdate update stream codec storage
func (obj *StorageST) StreamChannelCodecsUpdate(streamID string, val []av.CodecData, sdp []byte) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[streamID]; ok {
		tmp.codecs = val
		tmp.sdp = sdp
		obj.Streams[streamID] = tmp
	}
}

// StreamChannelSDP codec storage or wait
func (obj *StorageST) StreamChannelSDP(streamID string) ([]byte, error) {
	for i := 0; i < 100; i++ {
		obj.mutex.RLock()
		tmp, ok := obj.Streams[streamID]
		obj.mutex.RUnlock()
		if !ok {
			return nil, ErrorStreamNotFound
		}

		if len(tmp.sdp) > 0 {
			return tmp.sdp, nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return nil, ErrorStreamNotFound
}

func (obj *StorageST) StreamChannelAdd(uuid string, val StreamST) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if _, ok := obj.Streams[uuid]; ok {
		return ErrorStreamChannelAlreadyExists
	}
	val = obj.StreamChannelMake(val)
	obj.Streams[uuid] = val
	if !val.OnDemand {
		val.runLock = true
		go StreamServerRunStreamDo(uuid)
	}
	return nil
}
