package sfu

import (
	"sync"
	"testing"

	"github.com/pion/ion-sfu/pkg/buffer"
)

func TestWebRTCReceiverOnExtPktHandlerConcurrentAccess(t *testing.T) {
	var receiver WebRTCReceiver
	var wg sync.WaitGroup

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 2000; j++ {
				receiver.OnExtPktHandler(0, func(*buffer.ExtPacket) {})
				receiver.OnExtPktHandler(0, nil)
			}
		}()
	}

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 4000; j++ {
				if handler := receiver.onExtPktHandler[0].Load(); handler != nil && handler.fn != nil {
					handler.fn(&buffer.ExtPacket{})
				}
			}
		}()
	}

	wg.Wait()
}
