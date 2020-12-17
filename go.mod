module github.com/micro-plat/ddns

go 1.14

replace github.com/micro-plat/hydra => ../../../github.com/micro-plat/hydra

replace github.com/micro-plat/lib4go => ../../../github.com/micro-plat/lib4go

require (
	github.com/PuerkitoBio/goquery v1.5.1
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
	github.com/chromedp/chromedp v0.5.3
	github.com/fsnotify/fsnotify v1.4.9
	github.com/micro-plat/hydra v0.0.0-00010101000000-000000000000
	github.com/micro-plat/lib4go v1.0.2
	github.com/miekg/dns v1.1.34
	github.com/zkfy/go-cache v2.1.0+incompatible
	golang.org/x/sys v0.0.0-20201015000850-e3ed0017c211
	gopkg.in/fsnotify.v1 v1.4.7
)
