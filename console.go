package mylogger

import (
	"fmt"
	"time"
)

//日志结构体
type Consolelog struct {
	Consoleloglv LogLevel
}

//构造函数
//调用转换函数,把传入的string类型转成自定义的Level等级进行比较
func NewConsoleLog(lvstr string) *Consolelog {

	level, err := paserLoglevel(lvstr)
	if err != nil {
		fmt.Println("传入参数出错:", err)
		panic(err)
	}
	return &Consolelog{
		Consoleloglv: level,
	}
}

//比较方法，用于比较需要打印的和传入构造的等级
func (c *Consolelog) enable(level LogLevel) bool {
	return c.Consoleloglv <= level

}

//日志打印方法
func (c *Consolelog) log(lv LogLevel, format string, arg ...interface{}) {
	if c.enable(lv) {
		msg := fmt.Sprintf(format, arg...)
		t := time.Now()
		funcname, filename, lineNo := getInfo(3)
		fmt.Printf("[%s] [%s][%s: %s: %d] %s \n", t.Format("2006-01-02 15:04:05"), paserLogString(lv), filename, funcname, lineNo, msg)
	}
}

//日志等级调用打印方法
func (c *Consolelog) Debug(format string, arg ...interface{}) {
	if c.enable(DEBUG) {
		c.log(DEBUG, format, arg...)
	}
}

func (c *Consolelog) Info(format string, arg ...interface{}) {

	c.log(INFO, format, arg...)

}

func (c *Consolelog) Warning(format string, arg ...interface{}) {

	c.log(WARNING, format, arg...)

}

func (c *Consolelog) Error(format string, arg ...interface{}) {

	c.log(ERROR, format, arg...)

}

func (c *Consolelog) Fater(format string, arg ...interface{}) {

	c.log(FATAL, format, arg...)

}
