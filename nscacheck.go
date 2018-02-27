package nscacheck

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Syncbak-Git/nsca"
)

// NSCASettings is configuration to contact Shinken NSCA receiver
type NSCASettings struct {
	Server     string `json:"server"`     // Shinken receiver URL
	Port       int    `json:"port"`       // Shinken receiver Port
	Encryption int    `json:"encryption"` // Encryption - examples: 0 (none) or 1 (XOR)
	Key        string `json:"key"`        // NSCA key (password)
}

func NSCASettingsFromJSON(jsonStr string) (NSCASettings, error) {
	var nscaSettings NSCASettings
	err := json.Unmarshal([]byte(jsonStr), &nscaSettings)
	return nscaSettings, err
}

func Send(nscaSettings NSCASettings, host string, service string, state int, msg string) error {
	serverInfo := nsca.ServerInfo{
		Host:             nscaSettings.Server,
		Port:             fmt.Sprintf("%d", nscaSettings.Port),
		EncryptionMethod: nscaSettings.Encryption,
		Password:         nscaSettings.Key,
		Timeout:          5 * time.Second,
	}
	nscaServer := nsca.NSCAServer{}
	if err := nscaServer.Connect(serverInfo); err != nil {
		return err
	}
	defer nscaServer.Close()

	nscaMsg := nsca.Message{
		Host:    host,
		Service: service,
		State:   int16(state),
		Message: msg,
	}
	return nscaServer.Send(&nscaMsg)
}
