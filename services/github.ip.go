package services

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

var address = map[string]string{
	"github.com":                   "https://github.com.ipaddress.com",
	"assets-cdn.github.com":        "https://github.com.ipaddress.com/assets-cdn.github.com",
	"github.global.ssl.fastly.net": "https://fastly.net.ipaddress.com/github.global.ssl.fastly.net",
}

var checkURL = "https://github.com/"

// getIP 获取ip
func getIP(text string, address string) (string, error) {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(text))
	if err != nil {
		return "", err
	}
	ipA := ""
	dom.Find("table.table-v tbody tr td ul.comma-separated,table.faq tbody tr td ul").Each(func(i int, contentSelection *goquery.Selection) {
		ip := contentSelection.Find("li").Text()
		re, err := regexp.Compile("((?:(?:25[0-5]|2[0-4]\\d|[01]?\\d?\\d)\\.){3}(?:25[0-5]|2[0-4]\\d|[01]?\\d?\\d))")
		if err != nil {
			log.Fatal(err)
		}
		ipA = re.FindString(ip)
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
	for k, v := range address {
		url := v
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
		fmt.Println(k,ip)
		if ip == "" {
			continue
		}
		domains = append(domains, &Domain{Domain: k, IP: ip})
	}
	return
}

//Check 检查当前IP是否可用`
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
