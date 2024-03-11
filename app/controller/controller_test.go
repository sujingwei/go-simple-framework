/*
 * @Author: sujingwei 348149047@qq.com
 * @Date: 2024-03-11 14:22:48
 * @LastEditors: sujingwei 348149047@qq.com
 * @LastEditTime: 2024-03-11 15:49:10
 * @FilePath: \go-simple-framework\app\controller\controller_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package controller

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestController(t *testing.T) {
	controller := new(WelcomeController)

	if t := reflect.TypeOf(controller); t != nil {
		v := reflect.ValueOf(controller)
		if v.NumMethod() > 0 {
			for i := 0; i < v.NumMethod(); i++ {
				method := v.Method(i)
				fmt.Printf("method.Name[%s]: %+s\t\t", t.Elem().Name(), t.Method(i).Name)
				// fmt.Printf("参数个数：%d\t", method.Type().NumIn())
				if method.Type().NumIn() == 1 {
					if _, ok := method.Interface().(func(*gin.Context)); ok {
						fmt.Printf("isGin: %v\t", ok)
					} else {
						fmt.Printf("isGin: %v\t", ok)
					}
				}
				fmt.Println("\n---------------------------------------------------------------------------------------------------------")
			}
		}
	}

}
