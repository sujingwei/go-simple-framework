/*
 * @Author: sujingwei 348149047@qq.com
 * @Date: 2024-03-13 20:16:51
 * @LastEditors: sujingwei 348149047@qq.com
 * @LastEditTime: 2024-03-14 16:44:07
 * @FilePath: \go-simple-framework\web-framework\web_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package webframework

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestUUID(t *testing.T) {
	uid := strings.ReplaceAll(uuid.New().String(), "-", "")
	fmt.Printf("%T, %v, %s\n", uid, uid, uid)
}
