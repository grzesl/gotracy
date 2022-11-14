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

var tracyStringsMap map[string]*C.char = make(map[string]*C.char)
var allocStringMutex sync.Mutex

func allocString(text string) *C.char {

	allocStringMutex.Lock()
	defer allocStringMutex.Unlock()

	val, ok := tracyStringsMap[text]
	if ok {
		return val
	}

	cgotext := C.CString(text)
	tracyStringsMap[text] = cgotext
	return cgotext
}

func TracySetThreadName(name string) {
	runtime.LockOSThread()

	C.GoTracySetThreadName(allocString(name))
}

var tracyZoneBeginMutex sync.Mutex

func TracyZoneBegin(name string, color uint32) int {

	tracyZoneBeginMutex.Lock()
	defer tracyZoneBeginMutex.Unlock()

	pc, filename, line, _ := runtime.Caller(1)
	funcname := runtime.FuncForPC(pc).Name()

	ret := C.GoTracyZoneBegin(allocString(name), allocString(funcname),
		allocString(filename), C.uint(line), C.uint(color))

	return int(ret)
}

func TracyZoneEnd(c int) {
	C.GoTracyZoneEnd(C.int(c))
}

func TracyZoneValue(c int, value int64) {
	C.GoTracyZoneValue(C.int(c), C.uint64_t(value))
}

func TracyZoneText(c int, text string) {
	C.GoTracyZoneText(C.int(c), allocString(text))
}

func TracyMessageL(msg string) {
	C.GoTracyMessageL(allocString(msg))
}

func TracyMessageLC(msg string, color uint32) {
	C.GoTracyMessageLC(allocString(msg), C.uint(color))
}

func TracyFrameMark() {
	C.GoTracyFrameMark()
}

func TracyFrameMarkName(name string) {
	C.GoTracyFrameMarkName(allocString(name))
}

func TracyFrameMarkStart(name string) {
	C.GoTracyFrameMarkStart(allocString(name))
}

func TracyFrameMarkEnd(name string) {
	C.GoTracyFrameMarkEnd(allocString(name))
}

func TracyPlotFloat(name string, val float32) {
	C.GoTracyPlotFloat(allocString(name), C.float(val))
}

func TracyPlotDouble(name string, val float64) {
	C.GoTracyPlotDoublet(allocString(name), C.double(val))
}

func TracyPlotInt(name string, val int) {
	C.GoTracyPlotInt(allocString(name), C.int(val))
}

func TracyMessageAppinfo(name string) {
	C.GoTracyMessageAppinfo(allocString(name))
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

		id := TracyZoneBegin("BLABLA", 0xFF00FF)

		time.Sleep(time.Nanosecond * 100)

		id2 := TracyZoneBegin("SUB BLABLA", 0xFF00FF)
		TracyZoneValue(id2, 100)
		TracyZoneText(id2, "To jest ważna informacja...")

		TracyFrameMarkStart("sin(x)")
		time.Sleep(time.Nanosecond * 100)
		TracyPlotDouble("sin", math.Sin(float64(i)/10))
		TracyFrameMarkEnd("sin(x)")
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

}
*/
