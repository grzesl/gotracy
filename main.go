package gotracy

/*
#cgo CPPFLAGS: -DTRACY_ENABLE -D_WIN32_WINNT=0x0602 -DWINVER=0x0602
#cgo LDFLAGS:  -L./ -luser32 -lkernel32 -lws2_32 -lImageHlp -lDbghelp -Wl,-rpath=./
#include "gotracy.h"
#include <stdlib.h>
#include <stdio.h>
*/
import "C"

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"strconv"
	"sync"
	"time"
	"unsafe"

	"github.com/enriquebris/goconcurrentqueue"
	"github.com/undefinedlabs/go-mpatch"
)

var tracyStringsMap map[string]*C.char = make(map[string]*C.char)
var allocStringMutex sync.Mutex
var tracyMutex sync.Mutex

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
	tracyMutex.Lock()
	C.GoTracySetThreadName(allocString(name))
	tracyMutex.Unlock()
}

func TracyZoneBegin(name string, color uint32) int {

	tracyMutex.Lock()

	pc, filename, line, _ := runtime.Caller(1)
	funcname := runtime.FuncForPC(pc).Name()

	ret := C.GoTracyZoneBegin(allocString(name), allocString(funcname),
		allocString(filename), C.uint(line), C.uint(color))

	tracyMutex.Unlock()
	return int(ret)
}

func TracyZoneEnd(c int) {
	tracyMutex.Lock()
	C.GoTracyZoneEnd(C.int(c))
	tracyMutex.Unlock()
}

func TracyZoneValue(c int, value int64) {
	tracyMutex.Lock()
	C.GoTracyZoneValue(C.int(c), C.uint64_t(value))
	tracyMutex.Unlock()
}

func TracyZoneText(c int, text string) {
	tracyMutex.Lock()
	C.GoTracyZoneText(C.int(c), allocString(text))
	tracyMutex.Unlock()
}

func TracyMessageL(msg string) {
	tracyMutex.Lock()
	C.GoTracyMessageL(allocString(msg))
	tracyMutex.Unlock()
}

func TracyMessageLC(msg string, color uint32) {
	tracyMutex.Lock()
	C.GoTracyMessageLC(allocString(msg), C.uint(color))
	tracyMutex.Unlock()
}

func TracyFrameMark() {
	tracyMutex.Lock()
	C.GoTracyFrameMark()
	tracyMutex.Unlock()
}

func TracyFrameMarkName(name string) {
	tracyMutex.Lock()
	C.GoTracyFrameMarkName(allocString(name))
	tracyMutex.Unlock()
}

func TracyFrameMarkStart(name string) {
	tracyMutex.Lock()
	C.GoTracyFrameMarkStart(allocString(name))
	tracyMutex.Unlock()
}

func TracyFrameMarkEnd(name string) {
	tracyMutex.Lock()
	C.GoTracyFrameMarkEnd(allocString(name))
	tracyMutex.Unlock()
}

func TracyPlotFloat(name string, val float32) {
	tracyMutex.Lock()
	C.GoTracyPlotFloat(allocString(name), C.float(val))
	tracyMutex.Unlock()
}

func TracyPlotDouble(name string, val float64) {
	tracyMutex.Lock()
	C.GoTracyPlotDoublet(allocString(name), C.double(val))
	tracyMutex.Unlock()
}

func TracyPlotInt(name string, val int) {
	tracyMutex.Lock()
	C.GoTracyPlotInt(allocString(name), C.int(val))
	tracyMutex.Unlock()
}

func TracyMessageAppinfo(name string) {
	tracyMutex.Lock()
	C.GoTracyMessageAppinfo(allocString(name))
	tracyMutex.Unlock()
}

func TracyMemoryAlloc(ptr uint64, size int64, secure int) {
	tracyMutex.Lock()
	C.GoTracyMemoryAlloc(C.ulonglong(ptr), C.ulonglong(size), C.int(secure))
	tracyMutex.Unlock()
}

func TracyMemoryAllocNamed(ptr uint64, size int64, secure int, name string) {
	tracyMutex.Lock()
	C.GoTracyMemoryAllocNamed(C.ulonglong(ptr), C.ulonglong(size), C.int(secure), allocString(name))
	tracyMutex.Unlock()
}

func TracyMemoryFree(ptr uint64, secure int) {
	tracyMutex.Lock()
	C.GoTracyMemoryFree(C.ulonglong(ptr), C.int(secure))
	tracyMutex.Unlock()
}

var mutex sync.Mutex

func otherThread() {

	TracySetThreadName("otherThread")

	for {

		ido := TracyZoneBegin("TEST", 0xF0F0FA)
		time.Sleep(time.Millisecond * 15)
		TracyZoneValue(ido, 1000)
		time.Sleep(time.Millisecond * 50)
		TracyZoneValue(ido, 500)
		TracyZoneEnd(ido)
		TracyMessageL("Id from TEST: " + strconv.Itoa(ido))
		TracyMessageLC("MESSAGE FROM TEST ZONE", 0xFF0F0F)
		time.Sleep(time.Second * 3)
		TracyFrameMarkName("oThread")
		//TracyFrameMarkEnd("oThread")
	}
}

func mainTestProgram() {
	TracySetThreadName("mainThread")
	fmt.Println("TEST")
	i := 0
	//go otherThread()

	for {
		//ZoneScoped()

		makslice := make([]int, 10000)
		makslice[0] = i
		id := TracyZoneBegin("Main Zone", 0xFF00FF)

		time.Sleep(time.Nanosecond * 10000)

		id2 := TracyZoneBegin("Sub main zone", 0xFF00FF)
		TracyZoneValue(id2, 100)
		TracyZoneText(id2, "This is important message...")

		TracyFrameMarkStart("sin(x)")
		time.Sleep(time.Nanosecond * 10000)
		TracyPlotDouble("sin", math.Sin(float64(i)/10))
		TracyFrameMarkEnd("sin(x)")

		time.Sleep(time.Millisecond * 100)
		TracyZoneValue(id2, 100)
		TracyZoneEnd(id2)
		//FrameMark()
		i++
		TracyZoneEnd(id)

		TracyMessageL("Id from BlaBla: " + strconv.Itoa(id))

		if i > 1000 {
			break
		}

		//time.Sleep(time.Second * 1)

		gostr := fmt.Sprintf("Loop %d", i)

		log.Println(gostr)
		//TracyMessageL("TEST" + gostr)
		//TracyFrameMarkEnd("mThread")
		TracyFrameMark()
		//go otherThread()

		//	pprof.WriteHeapProfile(os.Stdout)
	}
}

type tflag uint8
type nameOff int32
type typeOff int32

type _type struct {
	size uintptr
}

//go:linkname tracy_resolveNameOff runtime.resolveNameOff
func tracy_resolveNameOff(ptrInModule unsafe.Pointer, off nameOff) string

//go:noinline
//go:linkname tracy_tracealloc runtime.tracealloc
func tracy_tracealloc(p unsafe.Pointer, size uintptr, typ *_type)

//go:noinline
//go:linkname tracy_tracefree runtime.tracefree
func tracy_tracefree(p unsafe.Pointer, size uintptr)

//go:noinline
//go:linkname tracy_tracegc runtime.tracegc
func tracy_tracegc()

//go:linkname tracy_systemstack runtime.systemstack
func tracy_systemstack(fn func())

//go:linkname tracy_lock runtime.lock
func tracy_lock(l *mymutex)

//go:linkname tracy_unlock runtime.unlock
func tracy_unlock(l *mymutex)

type lockRank int

type lockRankStruct struct {
	// static lock ranking of the lock
	rank lockRank
	// pad field to make sure lockRankStruct is a multiple of 8 bytes, even on
	// 32-bit systems.
	pad int
}

type mymutex struct {
	// Empty struct if lock ranking is disabled, otherwise includes the lock rank
	lockRankStruct
	// Futex-based impl treats it as uint32 key,
	// while sema-based impl as M* waitm.
	// Used to be a union, but unions break precise GC.
	key uintptr
}

type tracy_debug_struct struct {
	cgocheck           int32
	clobberfree        int32
	efence             int32
	gccheckmark        int32
	gcpacertrace       int32
	gcshrinkstackoff   int32
	gcstoptheworld     int32
	gctrace            int32
	invalidptr         int32
	madvdontneed       int32 // for Linux; issue 28466
	scavtrace          int32
	scheddetail        int32
	schedtrace         int32
	tracebackancestors int32
	asyncpreemptoff    int32
	harddecommit       int32
	adaptivestackstart int32

	// debug.malloc is used as a combined debug check
	// in the malloc function and should be set
	// if any of the below debug options is != 0.
	malloc         bool
	allocfreetrace int32
	inittrace      int32
	sbrk           int32
}

//go:linkname tracy_debug runtime.debug
var tracy_debug any

var allocMyMutex mymutex

type MemTraceItem struct {
	P    uint64
	Size int64
}

var memTraceQueue goconcurrentqueue.Queue

var memTraceSlice []MemTraceItem

var memTraceAllocSlice []MemTraceItem
var memTraceAllocSliceHead int = 0
var memTraceFreeSlice []MemTraceItem
var memTraceFreeSliceHead int = 0

var memTraceAlloc []uint64
var memTraceAllocCount int
var memTraceFreeCount int

const memTraceSliceCount int = 100000

func checkMallocFreeCount() {
	if memTraceAllocCount > len(memTraceAlloc)/2 {
		memTraceAllocCount = 0
		// refresh table
		for _, value := range memTraceAlloc {
			if value == 0 {
				continue
			}
			memTraceAlloc[memTraceAllocCount] = value
			memTraceAllocCount++
		}
	}
}

var mybreak bool = false

func newTraceAllocMemory(duration int) {

	for {
		time.Sleep(time.Duration(duration) * time.Millisecond)
		mybreak = false
		for i, val := range memTraceAllocSlice {

			tracy_lock(&allocMyMutex)

			if i >= memTraceAllocSliceHead {
				mybreak = true
				memTraceAllocSliceHead = 0
			} else if val.P == 0 {
				mybreak = true
				memTraceAllocSliceHead = 0
			}

			tracy_unlock(&allocMyMutex)

			if mybreak {
				break
			}
			var found int = 0
			if val.Size != 0 {

				for _, v := range memTraceAlloc {
					if v == val.P {
						found = 1
						break
					}
				}

				if found == 1 {
					continue //wait maybe will be free layter
				}

				TracyMemoryAlloc(val.P, val.Size, 1)
				memTraceAlloc[memTraceAllocCount] = val.P
				memTraceAllocCount++
			} else {

				for i, v := range memTraceAlloc {
					if v == val.P {
						memTraceAlloc[i] = 0
						TracyMemoryFree(val.P, 1)
						break
					}
				}
				checkMallocFreeCount()
			}

		}

	}
}

var traceAllocItem *MemTraceItem

func my_tracealloc(p unsafe.Pointer, size uintptr, typ *_type) {
	tracy_lock(&allocMyMutex)
	my_debug.allocfreetrace = 0
	traceFreeItem = &memTraceAllocSlice[memTraceAllocSliceHead]
	traceFreeItem.P = uint64(uintptr(p))
	traceFreeItem.Size = int64(size)
	memTraceAllocSliceHead++
	memTraceAllocSlice[memTraceAllocSliceHead].P = 0
	my_debug.allocfreetrace = 1
	tracy_unlock(&allocMyMutex)

}

var traceFreeItem *MemTraceItem

func my_tracefree(p unsafe.Pointer, size uintptr) {
	tracy_lock(&allocMyMutex)
	my_debug.allocfreetrace = 0
	traceFreeItem = &memTraceAllocSlice[memTraceAllocSliceHead]
	traceFreeItem.P = uint64(uintptr(p))
	traceFreeItem.Size = 0
	memTraceAllocSliceHead++
	memTraceAllocSlice[memTraceAllocSliceHead].P = 0
	my_debug.allocfreetrace = 1
	tracy_unlock(&allocMyMutex)
}

func my_tracegc() {

}

var my_debug *tracy_debug_struct

func TracyInit() {

	log.Print("TracyInit \n")
	runtime.LockOSThread()

	memTraceAllocSlice = make([]MemTraceItem, memTraceSliceCount)
	memTraceAllocSliceHead = 0

	memTraceAlloc = make([]uint64, memTraceSliceCount)
	memTraceAllocCount = 0

	/* //SET GODEBUG=allocfreetrace=1
	err := os.Setenv("GODEBUG", "allocfreetrace=1")
	if err != nil {
		log.Print(err)
	}
	*/

	_, err := mpatch.PatchMethod(tracy_tracealloc, my_tracealloc)
	if err != nil {
		log.Panic(err)
	}
	_, err = mpatch.PatchMethod(tracy_tracefree, my_tracefree)
	if err != nil {
		log.Panic(err)
	}
	_, err = mpatch.PatchMethod(tracy_tracegc, my_tracegc)
	if err != nil {
		log.Panic(err)
	}

	go newTraceAllocMemory(10)

	my_debug = (*tracy_debug_struct)(unsafe.Pointer(&tracy_debug))
	my_debug.allocfreetrace = 1
	my_debug.malloc = true

}

/*
func main() {
	TracyInit()
	mainTestProgram()
}
*/
