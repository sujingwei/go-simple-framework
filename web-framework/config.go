package webframework

/**
 * @description: Web配置
 * @return {*}
 */
type WebConfig struct {
	// 用于指定HTTP服务器监听的地址和端口。通常是一个字符串，例如:8080表示监听在本地所有IP地址的8080端口上。
	Addr string `json:"addr" yaml:"addr"`
	// 用于设置服务器读取请求的超时时间。如果在指定时间内还未接收到完整的请求数据，服务器将关闭连接。该值应该根据实际情况设置合理的数值。
	// 单位秒
	ReadTimeout int `json:"readTimeout" yaml:"readTimeout"`
	// 单位秒
	ReadHeaderTimeout int `json:"readHeaderTimeout" yaml:"readHeaderTimeout"`
	// 用于设置服务器写入响应的超时时间。如果在指定时间内未将完整的响应数据写入连接，服务器将关闭连接。同样，应该根据实际情况设置合理的数值。
	// 单位秒
	WriteTimeout int `json:"writeTimeout" yaml:"writeTimeout"`
	// 用于设置连接的空闲超时时间。当一个连接在指定时间内没有任何活动，服务器将关闭它。这个值对于节省系统资源和防止连接池溢出很有用。
	// 单位秒
	IdleTimeout int `json:"idleTimeout" yaml:"idleTimeout"`
	// 用于设置请求头的最大字节数。如果请求头超过该值，服务器将返回413 Request Header Fields Too Large响应。默认值为1MB
	MaxHeaderBytes int `json:"maxHeaderBytes" yaml:"maxHeaderBytes"`
}
