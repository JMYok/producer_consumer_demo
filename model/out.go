package model

import "fmt"

func init() {
	NewOut()
}

// Out 输出
type Out struct {
	data chan interface{}
}

// 单例模式
var OutObj *Out

// NewOut 初始化
func NewOut() *Out {
	if OutObj == nil {
		OutObj = &Out{
			data: make(chan interface{}, 65535), // 这里必须设置缓冲区
		}
	}
	return OutObj
}

// Println out 的写入方法
func (o *Out) Println(i interface{}) {
	OutObj.data <- i
}

// OutPut 将 out 内所有数据全部输出
func (o *Out) OutPut() {
	for {
		select {
		case i := <-o.data:
			fmt.Println(i)
		default:
			return
		}
	}
}
