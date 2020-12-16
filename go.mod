module github.com/micro-plat/ddns

go 1.14

replace github.com/micro-plat/hydra => ../../../github.com/micro-plat/hydra

replace github.com/micro-plat/lib4go => ../../../github.com/micro-plat/lib4go

require (
	github.com/PuerkitoBio/goquery v1.5.1
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
	github.com/chromedp/chromedp v0.5.3
	github.com/fsnotify/fsnotify v1.4.9
	github.com/go-sql-driver/mysql v1.5.0
	github.com/micro-plat/hydra v0.0.0-00010101000000-000000000000
	github.com/micro-plat/lib4go v1.0.2
	github.com/miekg/dns v1.1.34
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/sparrc/go-ping v0.0.0-20190613174326-4e5b6552494c
	github.com/zkfy/go-cache v2.1.0+incompatible
)
