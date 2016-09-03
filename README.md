# lily: hosts manager takes effect immediately on switch

![capture](https://cloud.githubusercontent.com/assets/8019222/18225305/4a399b78-7222-11e6-8e1e-9e0037c63d2b.PNG)

## Motivation
We usually have to switch between a bunch of test environments which have same domain but with different ip.
Since it's troublesome to modify hosts file each time we make switch, so we use some tools to manage it.
I've tried HostAdmin, SwitchHosts! and some other tools but they all have the some problem: we have to restart
browser in order to take effect after hosts switched.

So I decide to solve this problem by "Reinventing the wheel", and as a Go enthusiast, I choose Go to develop it.
If you're interested in this project please fork it and pull request is prefered :)

## Feature
* Takes effect immediately on swtich
* Easy to use
* Clean code
* Written purely in Go

## Build
```
get get -v github.com/tinycedar/lily
go build -ldflags="-H windowsgui"
```

## Run
Double click Lily.exe
