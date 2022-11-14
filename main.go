package gotracy

/*
#cgo CXXFLAGS: -DTRACY_ENABLE -D_WIN32_WINNT=0x0602 -DWINVER=0x0602
#cgo LDFLAGS:  -L./ -luser32 -lkernel32 -lws2_32 -lImageHlp -lDbghelp -Wl,-rpath=./
#include "gotracy.h"
#include <stdlib.h>
#include <stdio.h>
*/
import "C"

import (
	"runtime"
	"strconv"
	"sync"
	"time"
)

func TracySetThreadName(name string) {
	runtime.LockOSThread()

	cgoname := C.CString(name)
	C.GoTracySetThreadName(cgoname)
}

var tracyZoneBeginMutex sync.Mutex

func TracyZoneBegin(name string, color uint32) int {

	tracyZoneBeginMutex.Lock()
	defer tracyZoneBeginMutex.Unlock()

	pc, filename, line, _ := runtime.Caller(1)
	funcname := runtime.FuncForPC(pc).Name()

	cgoname := C.CString(name)
	cgofunc := C.CString(funcname)
	cgofilename := C.CString(filename)

	ret := C.GoTracyZoneBegin(cgoname, cgofunc, cgofilename, C.uint(line), C.uint(color))
	//defer C.free(unsafe.Pointer(cgoname))
	//defer C.free(unsafe.Pointer(cgofunc))
	//defer C.free(unsafe.Pointer(cgofilename))

	return int(ret)
}

func TracyZoneEnd(c int) {
	str := strconv.Itoa(c)
	cgostr := C.CString(str)
	C.GoTracyZoneEnd(cgostr)
}

func TracyZoneValue(c int, value int64) {
	str := strconv.Itoa(c)
	cgostr := C.CString(str)
	C.GoTracyZoneValue(cgostr, C.uint64_t(value))
}

func TracyZoneText(c int, text string) {
	str := strconv.Itoa(c)
	cgostr := C.CString(str)
	cgotext := C.CString(text)
	C.GoTracyZoneText(cgostr, cgotext)
}

func TracyMessageL(msg string) {
	cgomsg := C.CString(msg)
	C.GoTracyMessageL(cgomsg)
}

func TracyMessageLC(msg string, color uint32) {
	cgomsg := C.CString(msg)
	C.GoTracyMessageLC(cgomsg, C.uint(color))
}

func TracyFrameMark() {
	C.GoTracyFrameMark()
}

func TracyFrameMarkName(name string) {
	cgoname := C.CString(name)
	C.GoTracyFrameMarkName(cgoname)
}

func TracyFrameMarkStart(name string) {
	cgoname := C.CString(name)
	C.GoTracyFrameMarkStart(cgoname)
}

func TracyFrameMarkEnd(name string) {
	cgoname := C.CString(name)
	C.GoTracyFrameMarkEnd(cgoname)
}

func TracyPlotFloat(name string, val float32) {
	cgoname := C.CString(name)
	C.GoTracyPlotFloat(cgoname, C.float(val))
}

func TracyPlotDouble(name string, val float64) {
	cgoname := C.CString(name)
	C.GoTracyPlotDoublet(cgoname, C.double(val))
}

func TracyPlotInt(name string, val int) {
	cgoname := C.CString(name)
	C.GoTracyPlotInt(cgoname, C.int(val))
}

func TracyMessageAppinfo(name string) {
	cgoname := C.CString(name)
	C.GoTracyMessageAppinfo(cgoname)
}

var mutex sync.Mutex

func otherThread() {

	TracySetThreadName("otherThread")

	for {

		ido := TracyZoneBegin("TEST", 0xF0F0FA)
		time.Sleep(time.Millisecond * 1500)
		TracyZoneValue(ido, 1000)
		time.Sleep(time.Millisecond * 500)
		TracyZoneValue(ido, 500)
		TracyZoneEnd(ido)
		TracyMessageL("Id from TEST: " + strconv.Itoa(ido))
		TracyMessageLC("MESSAGE FROM TEST ZONE", 0xFF0F0F)
		time.Sleep(time.Second * 3)
		TracyFrameMarkName("oThread")
		//TracyFrameMarkEnd("oThread")
	}
}

/*
func main() {

	TracySetThreadName("mainThread")
	fmt.Println("TEST")
	i := 0
	go otherThread()

	for {
		//ZoneScoped()
		//TracyFrameMarkStart("mThread")
		id := TracyZoneBegin("BLABLA", 0xFF00FF)

		time.Sleep(time.Nanosecond * 100)

		id2 := TracyZoneBegin("SUB BLABLA", 0xFF00FF)
		TracyZoneValue(id2, 100)
		TracyZoneText(id2, "To jest ważna informacja...")

		TracyPlotDouble("sin", math.Sin(float64(i)/10))
		//TracyMessageLC("To jest ważna informacja: " + strconv.Itoa(id), 0xFF3344)

		time.Sleep(time.Nanosecond * 100)
		TracyZoneValue(id2, 100)
		TracyZoneEnd(id2)
		//FrameMark()
		i++
		TracyZoneEnd(id)

		TracyMessageL("Id from BlaBla: " + strconv.Itoa(id))

		if i > 1000 {
			break
		}

		time.Sleep(time.Second * 1)

		gostr := fmt.Sprintf("Loop %d", i)

		log.Println(gostr)
		//TracyMessageL("TEST" + gostr)
		//TracyFrameMarkEnd("mThread")
		TracyFrameMark()
		//go otherThread()

	}

}*/
