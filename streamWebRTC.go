package main

import (
	"time"

	webrtc "github.com/deepch/vdk/format/webrtcv3"
	"github.com/sirupsen/logrus"
)

// ServerStreamWebRTC stream video over WebRTC
func ServerStreamWebRTC(uuid string, data string) string {
	logger := log.WithFields(logrus.Fields{
		"module": "http_webrtc",
		"stream": uuid,
		"func":   "ServerStreamWebRTC",
	})

	if !Storage.StreamChannelExist(uuid) {
		var stream StreamST
		stream.URL = uuid
		err := Storage.StreamChannelAdd(uuid, stream)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"call": "StreamChannelExist",
			}).Errorln(err.Error())
			return ""
		}
	}

	codecs, err := Storage.StreamChannelCodecs(uuid)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"call": "StreamCodecs",
		}).Errorln(err.Error())
		return ""
	}
	muxerWebRTC := webrtc.NewMuxer(webrtc.Options{ICEServers: Storage.ServerICEServers(), ICEUsername: Storage.ServerICEUsername(), ICECredential: Storage.ServerICECredential(), PortMin: Storage.ServerWebRTCPortMin(), PortMax: Storage.ServerWebRTCPortMax()})
	answer, err := muxerWebRTC.WriteHeader(codecs, data)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"call": "WriteHeader",
		}).Errorln(err.Error())
		return ""
	}

	go func() {
		cid, ch, _, err := Storage.ClientAdd(uuid, WEBRTC)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"call": "ClientAdd",
			}).Errorln(err.Error())
			return
		}
		defer Storage.ClientDelete(uuid, cid)
		var videoStart bool
		noVideo := time.NewTimer(10 * time.Second)
		for {
			select {
			case <-noVideo.C:
				//				c.IndentedJSON(500, Message{Status: 0, Payload: ErrorStreamNoVideo.Error()})
				logger.WithFields(logrus.Fields{
					"call": "ErrorStreamNoVideo",
				}).Errorln(ErrorStreamNoVideo.Error())
				return
			case pck := <-ch:
				if pck.IsKeyFrame {
					noVideo.Reset(10 * time.Second)
					videoStart = true
				}
				if !videoStart {
					continue
				}
				err = muxerWebRTC.WritePacket(*pck)
				if err != nil {
					logger.WithFields(logrus.Fields{
						"call": "WritePacket",
					}).Errorln(err.Error())
					return
				}
			}
		}
	}()

	return answer
}
