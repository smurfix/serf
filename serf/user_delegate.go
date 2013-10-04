package serf

import (
	"log"
)

// NodeMeta is used to retrieve meta-data about the current node
func (s *Serf) NodeMeta(limit int) []byte {
	// Use the meta field for the role
	roleBytes := []byte(s.conf.Role)
	if len(roleBytes) > limit {
		roleBytes = roleBytes[:limit]
	}
	return roleBytes
}

// NotifyMsg is called when a user-data message is received.
// This should not block
func (s *Serf) NotifyMsg(buf []byte) {
	if len(buf) == 0 {
		return
	}

	rebroadcast := false
	msgType := messageType(buf[0])
	switch msgType {
	case leaveMsg:
		l := leave{}
		if err := decode(buf[1:], &l); err != nil {
			log.Printf("[ERR] Decoding leave message failed: %v", err)
		}

		rebroadcast = s.intendLeave(&l)
	default:
		log.Printf("[WARN] Received message of unknown type %d", msgType)
	}

	// Check if we should rebroadcast
	if rebroadcast {
		s.rebroadcast(buf)
	}
}

// GetBroadcasts is called when user data messages can be broadcast.
// It can return a list of buffers to send. Each buffer should assume an
// overhead as provided with a limit on the total byte size allowed.
func (s *Serf) GetBroadcasts(overhead, limit int) [][]byte {
	return s.broadcasts.GetBroadcasts(overhead, limit)
}

// LocalState is used for a TCP Push/Pull. This is sent to
// the remote side as well as membership information
func (s *Serf) LocalState() []byte {
	// Do not push any state
	return nil
}

// MergeRemoteState is invoked after a TCP Push/Pull. This is the
// state received from the remote side.
func (s *Serf) MergeRemoteState(buf []byte) {
	// Do not merge any remote state
}
