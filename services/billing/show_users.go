package billing

import (
	"fmt"
	"github.com/liuguangw/billing_go/common"
)

// ShowUsers 通知服务器显示用户列表
func (s *Server) ShowUsers() error {
	packet := &common.BillingPacket{
		MsgID:  [2]byte{0, 0},
		OpData: []byte("show_users"),
	}
	response, err := s.sendPacketToServer(packet)
	if err != nil {
		return err
	}
	fmt.Println(string(response.OpData))
	return nil
}
