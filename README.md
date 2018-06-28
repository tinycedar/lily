# Lily [![GitHub release](https://img.shields.io/github/release/tinycedar/lily.svg)](https://github.com/tinycedar/lily/releases) [![license](https://img.shields.io/github/license/tinycedar/lily.svg)]()

Hosts manager takes effect immediately on switch

![](https://ws2.sinaimg.cn/large/006tNbRwgy1fkuanxlo29j30qa0etq2y.jpg)

## Motivation
We usually have to switch between a bunch of test environments which have same domains but different IPs.
Since it's troublesome to modify hosts directly, we use some tools to manage it.
I've tried HostAdmin, SwitchHosts!, HostManager and so forth but found that they all have the same problem: browser relaunch
is required in order to take effect after switch.

So I decide to solve this problem by means of "Reinventing the wheel", and as a Go enthusiast, I choose Go to develop it.
If you're interested in this project please fork it and pull request is preferred :)

## Features
* Takes effect immediately on switch
* Single process
* Realtime notification
* Hosts config saved automatically once you made modification
* Written purely in Go

## Supported browsers
* Chrome
* Firefox
* Opera
* Sogou Explorer (搜狗高速浏览器)
* QQ Browser (QQ浏览器)
* 360 Browser (360极速浏览器)
* 360 Secure Browser (360安全浏览器)
* Liebao (Cheetah) Browser (猎豹安全浏览器)
* Maxthon (遨游)

## Roadmap
Only Windows is supported recently, macOS and Linux will be supported in the near future.
Probably I will use [Electron](http://electron.atom.io/) framework to implement cross-platform feature.

## Install
Download binary in [Releases](https://github.com/tinycedar/lily/releases)

## Build
```
go get -u -v github.com/tinycedar/lily
go build -ldflags="-H windowsgui"
```
