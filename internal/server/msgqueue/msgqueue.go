package msgqueue

import (
	"github.com/zyedidia/generic/queue"
	"net"
	"sync"
)

type Message struct {
	Addr *net.UDPAddr
	Msg  []byte
}

func New() *MsgQueue {
	return &MsgQueue{Queue: queue.New[Message]()}
}

// MsgQueue - TODO: Need to use best practice (channels instead of mutexes)
type MsgQueue struct {
	sync.Mutex
	Queue *queue.Queue[Message]
}
