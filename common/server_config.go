package common

// billing类型定义
const (
	// 经典
	BillTypeCommon = iota
	// 怀旧
	BillTypeHuaiJiu
)

// ServerConfig billing服务配置
type ServerConfig struct {
	IP               string   //billing服务器listen的IP
	Port             int      //billing服务器listen的端口
	DbHost           string   `json:"db_host" yaml:"db_host"`                         //数据库主机名或者IP
	DbPort           int      `json:"db_port" yaml:"db_port"`                         //数据库端口
	DbUser           string   `json:"db_user" yaml:"db_user"`                         //数据库用户名
	DbPassword       string   `json:"db_password" yaml:"db_password"`                 //数据库密码
	DbName           string   `json:"db_name" yaml:"db_name"`                         //数据库名
	AllowOldPassword bool     `json:"allow_old_password" yaml:"allow_old_password"`   //是否启用oldPassword(除非报这个错误,否则不建议开启)
	AutoReg          bool     `json:"auto_reg" yaml:"auto_reg"`                       //是否开启自动注册
	AllowIps         []string `json:"allow_ips" yaml:"allow_ips"`                     //允许连接billing的服务端IP,为空则表示不限制
	PointFix         int      `json:"point_fix" yaml:"point_fix"`                     //用于查询点数时少1或者多1点的修正(修正值一般为0或者1)
	MaxClientCount   int      `json:"max_client_count" yaml:"max_client_count"`       //最多允许进入的用户数量(0表示无限制)
	PcMaxClientCount int      `json:"pc_max_client_count" yaml:"pc_max_client_count"` //每台电脑最多允许进入的用户数量(0表示无限制)
	BillType         int      `json:"bill_type" yaml:"bill_type"`                     //billing类型 0经典 1怀旧
}
