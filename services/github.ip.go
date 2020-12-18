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

var address = []string{"github.com", "assets-cdn.github.com", "github.global.ssl.fastly.net"}

var searchPage = "https://www.ipaddress.com/search"
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
	for _, v := range address {
		url := v
		err = chromedp.Run(ctx, chromedp.Tasks{
			chromedp.Navigate(searchPage),
			chromedp.WaitVisible(`form`, chromedp.ByQuery),
			chromedp.SendKeys(`body > div > main > form > input[type=text]`, url, chromedp.BySearch),
			chromedp.Click("body > div > main > form > button:nth-child(2)", chromedp.ByQuery),
			chromedp.Sleep(5 * time.Second),
			chromedp.OuterHTML(`body`, &res, chromedp.ByQuery),
			chromedp.Sleep(5 * time.Second),
		})

		if err != nil {
			return nil, fmt.Errorf("搜索出错%s %w", url, err)
		}
		ip, err := getIP(res, v)
		if err != nil {
			return nil, err
		}
		if ip == "" {
			continue
		}
		domains = append(domains, &Domain{Domain: url, IP: ip})
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
