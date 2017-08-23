package systree

import (
	"encoding/json"

	"github.com/troian/surgemq/message"
	"github.com/troian/surgemq/types"
)

type server struct {
	version      string
	upTime       *dynamicValueUpTime
	currTime     *dynamicValueCurrentTime
	capabilities struct {
		SupportedVersions             []message.ProtocolVersion
		MaxQoS                        string
		MaxConnections                uint64
		MaximumPacketSize             uint32
		ServerKeepAlive               uint16
		ReceiveMaximum                uint16
		RetainAvailable               bool
		WildcardSubscriptionAvailable bool
		SubscriptionIDAvailable       bool
		SharedSubscriptionAvailable   bool
	}
}

func newServer(topicPrefix string, dynRetains, staticRetains *[]types.RetainObject) server {
	b := server{
		upTime:   newDynamicValueUpTime(topicPrefix + "/uptime"),
		currTime: newDynamicValueCurrentTime(topicPrefix + "/datetime"),
		version:  "1.0.0",
	}

	m, _ := message.NewMessage(message.ProtocolV311, message.PUBLISH)
	msg, _ := m.(*message.PublishMessage)
	msg.SetQoS(message.QoS0)               // nolint: errcheck
	msg.SetTopic(topicPrefix + "/version") // nolint: errcheck
	msg.SetPayload([]byte(b.version))

	*dynRetains = append(*dynRetains, b.upTime)
	*dynRetains = append(*dynRetains, b.currTime)
	*staticRetains = append(*staticRetains, msg)

	m, _ = message.NewMessage(message.ProtocolV311, message.PUBLISH)
	msg, _ = m.(*message.PublishMessage)
	msg.SetQoS(message.QoS0)                    // nolint: errcheck
	msg.SetTopic(topicPrefix + "/capabilities") // nolint: errcheck

	if data, err := json.Marshal(&b.capabilities); err == nil {
		msg.SetPayload(data)
	} else {
		msg.SetPayload([]byte(err.Error()))
	}

	*staticRetains = append(*staticRetains, msg)

	return b
}
