package dc1server

const (

	/**
	 * 设备上线1  {"action":"activate=","uuid":"activate=e28","auth":"","params":{"device_type":"PLUG_DC1_7","mac":"A4:7B:9D:06:A0:E6"}}
	 */
	ACTIVATE = "activate="
	/**
	 * 设备上线2
	 */
	IDENTIFY = "identify"
	// 不知道作用的请求
	DATETIME = "datetime"
	/**
	 * 每增加50kwh，自动上报
	 */
	DETAL_KWH = "kWh+"
	/**
	 * 查询设备状态
	 */
	DATAPOINT = "datapoint"
	/**
	 * 设置设备开关
	 */
	SET_DATAPOINT = "datapoint="
	//电量增加
	KWH = "kWh+"

	CODE_SUCCESS = 200

	// dc1命令下发队列
	Dc1CommandSendQueue = "dc1:command:send"

	// dc1定时任务下发队列
	Dc1CommandPlanQueue = "dc1:command:plan"

	Dc1StatusReplyQueue = "dc1:status:reply"
)
