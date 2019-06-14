package server

import (
	"billing/bhandler"
	"billing/tools"
	"database/sql"
	"fmt"
	"net"
)

func processBillingData(request *bhandler.BillingData, db *sql.DB, conn *net.TCPConn,
	handlers []bhandler.BillingHandler) error {
	var response *bhandler.BillingData = nil
	for _, handler := range handlers {
		if request.OpType == handler.GetType() {
			response = handler.GetResponse(request)
			break
		}
	}
	if response != nil {
		//响应
		responseData := response.PackData()
		_, err := conn.Write(responseData)
		if err != nil {
			return err
		}
	} else {
		//无法处理当前请求类型
		tools.ShowErrorInfoStr(fmt.Sprintf("unknown BillingData \n\tOpType: 0x%X \n\tOpData: %v",
			int(request.OpType), request.OpData))
	}
	return nil
}
