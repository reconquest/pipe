package pipe

import "sync/atomic"

var (
	flowID   int64 = 0
	streamID int64
)

func getFlowID() int {
	id := atomic.AddInt64(&flowID, 1)
	return int(id)
}

func getStreamID() int {
	id := atomic.AddInt64(&streamID, 1)
	return int(id)
}
