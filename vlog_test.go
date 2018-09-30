/*
@Time : 2018/9/30 14:13
@Author : chang
@File : vlog_test
@Software: GoLand
@Info:
*/
package vlogs

import (
	"testing"
	"time"
)

func TestSetOutputWithFileWithTimeFormat(t *testing.T) {
	SetOutputWithFileWithTimeFormat("test", "C:\\Users\\chang\\Documents\\server.log", "", 1, "2006-01-02-15-04 ")
	for true {
		time.Sleep(time.Second)
		Debug(0, time.Now().String())
	}
	time.Sleep(time.Hour * 24)
}
