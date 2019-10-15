package cpool

type ICPoolMannger interface {

	// 获取一个链接
	Get(servId string) (c interface{}, err error)

	// 注册服务
	Register(servId string, options ServOptions) error

	// 将链接放回
	Release(servId string, c interface{})
}

type IServ interface {

	// 初始化链接池
	init(options ServOptions) error

	// 返回配置
	options() (options ServOptions)

	// 维持心跳
	heartbeat(c interface{}) error

	// 连接
	connect() (c interface{}, err error)

	// 关闭链接
	destroy(c interface{}) error
}
