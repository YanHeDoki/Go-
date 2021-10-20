package mylogger

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

//filelog 基本对象设置
type filelog struct {
	Level       LogLevel
	filePath    string   //日志文件路径
	fileName    string   //日志文件名
	fileObj     *os.File //日志文件
	errFileObj  *os.File //错误日志文件
	maxFileSize int64    //最大文件大小
	checkmodel  string   //切割模式，输入time以小时分隔，输入size以文件大小分隔
	logH        int      //按几小时切割
}

//获取现在的时间的小时
var logtime int

//构造方法
func NewFilelog(lv, fp, fn, model string, logH int, maxsize int64) *filelog {

	loglevel, err := paserLoglevel(lv)
	if err != nil {
		panic(err)

	}
	f1 := &filelog{
		Level:       loglevel,
		filePath:    fp,
		fileName:    fn,
		maxFileSize: maxsize,
		logH:        logH,
		checkmodel:  model,
	}
	err2 := f1.initFile()
	if err2 != nil {
		panic(err2)
	}

	return f1
}

//创建日志的文件
func (f *filelog) initFile() error {
	logtime = time.Now().Hour()
	fullFileName := path.Join(f.filePath, f.fileName)
	fileObj, err := os.OpenFile(fullFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("open log file err:", err)
		return err
	}
	//err日志定义为xxxErr.txt
	// errfilname := strings.Split(f.fileName, ".")[0]
	// errfilname = errfilname + "Err" + ".log"
	// errfullFileName := path.Join(f.filePath, errfilname)

	errFilName := strings.Split(fullFileName, ".")[0]
	errfileObj, err := os.OpenFile(errFilName+"Err.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("open errlog file err:", err)
		return err
	}
	f.errFileObj = errfileObj
	f.fileObj = fileObj
	//f.Close()
	return nil

}

//比较方法，用于比较需要打印的和传入构造的等级
func (f *filelog) enable(loglevel LogLevel) bool {
	return f.Level <= loglevel
}

//校验切割模式方法
func (f *filelog) chckModel(mode string) bool {
	if mode == "time" {
		return f.checkHour(f.logH)
	} else if mode == "size" {
		return f.checkSize(f.fileObj)
	} else {
		fmt.Println("切割方式选择出错,未能切割")
	}
	return false
}

//时间切割对比方法
func (f *filelog) checkHour(t int) bool {

	return logtime+t == time.Now().Hour()
}

//大小切割对比方法
func (f *filelog) checkSize(file *os.File) bool {
	fi, err := file.Stat()
	if err != nil {
		fmt.Println("file get stat", err)
		return false
	}

	return fi.Size() >= f.maxFileSize
}

//切割文件方法 复用
func (f *filelog) splitFile(file *os.File) (*os.File, error) {
	//需要切割的文件
	nowStr := time.Now().Format("20060102150405")
	fi, err := file.Stat()
	if err != nil {
		fmt.Println("file getSata err:", err)
		return nil, err
	}
	logName := path.Join(f.filePath, fi.Name())
	nowFileName := fmt.Sprintf("%s.bak%s", logName, nowStr)

	//关闭当前的日志文件
	file.Close()

	//备份一下 rename
	os.Rename(logName, nowFileName)

	//打开新的文件
	f2, err := os.OpenFile(logName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("open new logFile err:", err)
		return nil, err
	}
	logtime = time.Now().Hour()
	//打开的新的日志对象赋值给F.fileobj
	return f2, nil
}

//log 日志打印方法
func (f *filelog) log(lv LogLevel, format string, arg ...interface{}) {
	if f.enable(lv) {
		msg := fmt.Sprintf(format, arg...)
		t := time.Now()
		funcname, filename, lineNo := getInfo(3)

		if f.chckModel(f.checkmodel) {
			newfile, err := f.splitFile(f.fileObj)
			if err != nil {
				return
			}
			f.fileObj = newfile
		}
		fmt.Fprintf(f.fileObj, "[%s] [%s][%s: %s: %d] %s \n", t.Format("2006-01-02 15:04:05"), paserLogString(lv), filename, funcname, lineNo, msg)
		if lv >= ERROR {

			if f.checkSize(f.errFileObj) {
				f2, err := f.splitFile(f.errFileObj)
				if err != nil {
					return
				}
				f.errFileObj = f2
			}
			fmt.Fprintf(f.errFileObj, "[%s] [%s][%s: %s: %d] %s \n", t.Format("2006-01-02 15:04:05"), paserLogString(lv), filename, funcname, lineNo, msg)
		}
	}
}

//更改切割模式
func (f *filelog) SetChckModel(model string) {
	f.checkmodel = model
}

//日志等级调用打印方法

func (f *filelog) Debug(format string, arg ...interface{}) {

	f.log(DEBUG, format, arg...)

}

func (f *filelog) Info(format string, arg ...interface{}) {

	f.log(INFO, format, arg...)
}

func (f *filelog) Warning(format string, arg ...interface{}) {

	f.log(WARNING, format, arg...)

}

func (f *filelog) Error(format string, arg ...interface{}) {

	f.log(ERROR, format, arg...)

}

func (f *filelog) Fater(format string, arg ...interface{}) {

	f.log(FATAL, format, arg...)

}

//关闭文件
func (f *filelog) Close() {
	f.fileObj.Close()
	f.errFileObj.Close()
}
