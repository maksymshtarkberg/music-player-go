package messaging

import (
	"log"

	"github.com/nats-io/nats.go"
)

var nc *nats.Conn

func InitNATS(url string) {
	var err error
	nc, err = nats.Connect(url)
	if err != nil {
		log.Fatal(err)
	}
}

func GetNATSConn() *nats.Conn {
	return nc
}

func Publish(subject string, message []byte) error {
	return nc.Publish(subject, message)
}

func Subscribe(subject string, handler func(msg *nats.Msg)) {
	nc.Subscribe(subject, handler)
}

func StartSubscribers() {
	Subscribe("user.register", HandleRegistration)

	Subscribe("user.login", HandleLogin)

}
