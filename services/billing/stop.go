package billing

import (
	"github.com/liuguangw/billing_go/common"
)

// Stop 发送停止命令到server
func (s *Server) Stop() error {
	packet := &common.BillingPacket{
		MsgID:  [2]byte{0, 0},
		OpData: []byte("stop"),
	}
	if _, err := s.sendPacketToServer(packet); err != nil {
		return err
	}
	return nil
}
