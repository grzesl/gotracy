# GoTracy - About
This project is helper binding library. It allow run Tracy Client in  GoLang projects. Whole library is based on client/server code wersion 0.9 avaliable on:
https://github.com/wolfpld/tracy

Note: Library curently works only on Windows.

# Quick start
To install library just run commands:

>go get github.com/grzesl/gotracy

To commpile library you need installed gcc, compiler wich i prefer is placed below:

https://jmeubank.github.io/tdm-gcc/download/


# Example of use

    package main

    import (
        "time"

        "github.com/grzesl/gotracy"
    )

    func exampleFunction() {
        gotracy.TracySetThreadName("exampleFunction")
        for {
            ido := gotracy.TracyZoneBegin("TEST", 0xF0F0FA)
            time.Sleep(time.Millisecond * 1500)
            gotracy.TracyZoneValue(ido, 1500)
            time.Sleep(time.Millisecond * 500)
            gotracy.TracyZoneValue(ido, 500)
            gotracy.TracyMessageLC("MESSAGE FROM TEST ZONE", 0xFF0F0F)
            gotracy.TracyZoneEnd(ido)
            time.Sleep(time.Second * 3)
            gotracy.TracyFrameMarkName("thread")
        }
    }

    func main() {
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
![Tracy](/images/tracy_sin.png)

# License 
Library is MIT licensed