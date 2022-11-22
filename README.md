# GoTracy - About
This project is helper binding library. It allow run Tracy Client in  GoLang projects. Whole library is based on client/server code wersion 0.9 avaliable on:
https://github.com/wolfpld/tracy

Note: Library curently works only on Windows.

# Quick start - Windows
To install library just run commands:

>go get github.com/grzesl/gotracy

To commpile library you need installed gcc compiler. I prefer one placed below:

https://jmeubank.github.io/tdm-gcc/download/


# Example of use

    package main

    import (
        "math"
        "strconv"
        "time"

        "github.com/grzesl/gotracy"
    )

    func exampleFunction() {
        gotracy.TracySetThreadName("exampleFunction")
        for i := 0.0; i < math.Pi; i += 0.1 {

            zoneid := gotracy.TracyZoneBegin("Calculating Sin", 0xF0F0F0)

            time.Sleep(time.Millisecond * 150)
            gotracy.TracyZoneValue(zoneid, 150)
            time.Sleep(time.Millisecond * 50)
            gotracy.TracyZoneText(zoneid, "Sleep 50")

            gotracy.TracyFrameMarkStart("Calculating sin(x)")
            sin := math.Sin(i)
            gotracy.TracyFrameMarkEnd("Calculating sin(x)")

            gotracy.TracyMessageLC("Sin(x) = "+strconv.FormatFloat(sin, 'E', -1, 64), 0xFF0F0F)
            gotracy.TracyPlotDouble("sin(x)", sin)

            gotracy.TracyZoneEnd(zoneid)

            gotracy.TracyFrameMark()
        }
    }

    func main() {
        gotracy.TracyInit()
        exampleFunction()
    }



Example output is similar to:

![Tracy](/images/tracy_example.png)

# Currently suported functions
    func TracySetThreadName(name string) 
    func TracyZoneBegin(name string, color uint32) int 
    func TracyZoneEnd(c int) 
    func TracyZoneValue(c int, value int64) 
    func TracyZoneText(c int, text string)
    func TracyMessageL(msg string) 
    func TracyMessageLC(msg string, color uint32)
    func TracyFrameMark() 
    func TracyFrameMarkName(name string) 
    func TracyFrameMarkStart(name string)
    func TracyFrameMarkEnd(name string) 
    func TracyPlotFloat(name string, val float32) 
    func TracyPlotDouble(name string, val float64) 
    func TracyPlotInt(name string, val int) 
    func TracyMessageAppinfo(name string) 

# Other screenshots
![Tracy Sinus](/images/tracy_sin.png)
![Tracy Memory](/images/tracy_memory.png)

# License 
Library is MIT licensed