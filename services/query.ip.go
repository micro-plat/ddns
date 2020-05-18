package services

import (
	"context"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/micro-plat/lib4go/types"
)

// GetIP 获取ip
func GetIP(text string, address string) (string, error) {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(text))
	if err != nil {
		return "", err
	}
	ipA := ""
	ttlA := 0
	dom.Find("li.ReListCent").Each(func(i int, contentSelection *goquery.Selection) {
		ip := contentSelection.Find("div.w60-0").Text()
		ttl := contentSelection.Find("div.w14-0").Text()
		_, err := strconv.ParseFloat(ttl, 64)
		if err == nil && ip != "-" {
			// 匹配ip
			re, err := regexp.Compile("((?:(?:25[0-5]|2[0-4]\\d|[01]?\\d?\\d)\\.){3}(?:25[0-5]|2[0-4]\\d|[01]?\\d?\\d))")
			if err != nil {
				log.Fatal(err)
			}
			if ttlA == 0 {
				ttlA = types.GetInt(ttl)
				ipA = re.FindString(ip)
			}
			if types.GetInt(ttl) > ttlA {
				ttlA = types.GetInt(ttl)
				ipA = re.FindString(ip)
			}
		}
	})

	return ipA, nil
}

// RemoteChromedp chromedp请求
func RemoteChromedp() (domains []*Domain, err error) {

	address := []string{"github.com", "github.global.ssl.fastly.net", "assets-cdn.github.com"}
	pageurl := "http://tool.chinaz.com/dns?type=1&host="

	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithDebugf(log.Printf))
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	var res string
	for _, v := range address {
		err = chromedp.Run(ctx, chromedp.Tasks{chromedp.Navigate(pageurl + v), chromedp.Sleep(5 * time.Second),
			chromedp.OuterHTML(`body`, &res, chromedp.ByQuery), chromedp.Sleep(5 * time.Second),
		})

		if err != nil {
			return nil, err
		}
		ip, err := GetIP(res, v)
		if err != nil {
			return nil, err
		}
		if ip == "" {
			continue
		}
		domains = append(domains, &Domain{Domain: v, IP: ip})
	}
	return
}
