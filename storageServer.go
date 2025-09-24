package main

import (
	"github.com/sirupsen/logrus"
)

// ServerLogLevel read debug options
func (obj *StorageST) ServerLogLevel() logrus.Level {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.LogLevel
}

// ServerICEServers read ICE servers
func (obj *StorageST) ServerICEServers() []string {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	return obj.Server.ICEServers
}

// ServerICEServers read ICE username
func (obj *StorageST) ServerICEUsername() string {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	return obj.Server.ICEUsername
}

// ServerICEServers read ICE credential
func (obj *StorageST) ServerICECredential() string {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	return obj.Server.ICECredential
}

// ServerWebRTCPortMin read WebRTC Port Min
func (obj *StorageST) ServerWebRTCPortMin() uint16 {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	return obj.Server.WebRTCPortMin
}

// ServerWebRTCPortMax read WebRTC Port Max
func (obj *StorageST) ServerWebRTCPortMax() uint16 {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	return obj.Server.WebRTCPortMax
}
