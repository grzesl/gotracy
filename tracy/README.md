# GoTracy - About
This project is helper binding library. It allow run Tracy Client in  GoLang projects. Whole library is based on client/server code wersion 0.9 avaliable on:
https://github.com/wolfpld/tracy

# Quick start
To install library just run commands:

>go get github.com/grzesl/gotracy

# Example of use

    func exampleFunction() {
        TracySetThreadName("exampleFunction")
        for {
            ido := TracyZoneBegin("TEST", 0xF0F0FA)
            time.Sleep(time.Millisecond * 1500)
            TracyZoneValue(ido, 1500)
            time.Sleep(time.Millisecond * 500)
            TracyZoneValue(ido, 500)
            TracyMessageLC("MESSAGE FROM TEST ZONE", 0xFF0F0F)
            TracyZoneEnd(ido)
            time.Sleep(time.Second * 3)
            TracyFrameMarkName("thread")
            }
        }

Example output is similar to:

![Tux, the Linux mascot](/images/tracy_example.png)


# License 
Library is MIT licensed