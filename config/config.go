/*
 * @Author: sujingwei 348149047@qq.com
 * @Date: 2024-03-12 19:13:05
 * @LastEditors: sujingwei 348149047@qq.com
 * @LastEditTime: 2024-03-12 20:37:16
 * @FilePath: \go-simple-framework\config\config.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package config

import (
	"github.com/sujingwei/go-simple-framework/configuration"
	webframework "github.com/sujingwei/go-simple-framework/web-framework"
)

type AppConfigure struct {
	App configuration.App      `yaml: "app"`
	Web webframework.WebConfig `yaml: "web`
}
