package services

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/micro-plat/lib4go/types"
)

var address = []string{"github.com", "developer.github.com", "github.global.ssl.fastly.net", "assets-cdn.github.com"}

var pageURL = "http://tool.chinaz.com/dns?type=1&host="
var checkURL = "https://github.com/"

// getIP 获取ip
func getIP(text string, address string) (string, error) {
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

// GetGithubDomains chromedp请求
func GetGithubDomains() (domains []*Domain, err error) {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithErrorf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	var res string
	for _, v := range address {
		url := pageURL + v
		err = chromedp.Run(ctx, chromedp.Tasks{
			chromedp.Navigate(url),
			chromedp.Sleep(5 * time.Second),
			chromedp.OuterHTML(`body`, &res, chromedp.ByQuery),
			chromedp.Sleep(5 * time.Second),
		})

		if err != nil {
			return nil, fmt.Errorf("请求出错%s %w", url, err)
		}
		ip, err := getIP(res, v)
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

//Check 检查当前IP是否可用
func Check() error {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithErrorf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(checkURL),
		chromedp.Sleep(5 * time.Second),
	})

	if err != nil {
		return fmt.Errorf("请求出错%s %w", checkURL, err)
	}
	return nil

}
