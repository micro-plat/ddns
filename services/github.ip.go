package services

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/micro-plat/hydra"
)

var address = map[string]string{
	"github.com":                   "https://github.com.ipaddress.com",
	"assets-cdn.github.com":        "https://github.com.ipaddress.com/assets-cdn.github.com",
	"github.global.ssl.fastly.net": "https://fastly.net.ipaddress.com/github.global.ssl.fastly.net",
}

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
	for domain, url := range address {
		res, status, err := hydra.C.HTTP().GetRegularClient().Get(url)
		if err != nil || status != 200 {
			return nil, fmt.Errorf("请求出错%s %d %w", url, status, err)
		}
		ip, err := getIP(res, url)
		if err != nil {
			return nil, err
		}
		if ip == "" {
			return nil, fmt.Errorf("地址解析失败，无法正确获得IP地址信息")
		}
		domains = append(domains, &Domain{Domain: domain, IP: ip})
	}
	return
}
