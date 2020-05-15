from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.action_chains import ActionChains
from selenium.webdriver.common.keys import Keys
import re
import requests


print("----------------查询最快的github.com DNS IP--------------------")

addrs = ["github.com", "assets-cdn.github.com", "github.global.ssl.fastly.net"]

for addr in addrs:

    chrome_options = webdriver.ChromeOptions()
    chrome_options.add_argument('--headless')  # 无界面选项
    browser = webdriver.Chrome(chrome_options=chrome_options)
    print("1. 从tool.chinaz.com获取"+addr+"所有可用的IP")

    url = "http://tool.chinaz.com/dns?type=1&host="+addr
    print(url)
    browser.get(url)
    content = browser.find_element_by_class_name("DnsResuListWrap")
    ttl = 0
    ip = ""
    print("2. 查找速度最快的IP")
    for li in browser.find_elements_by_xpath("//ul[@class='DnsResuListWrap fl DnsWL']/li"):
        nip = li.find_element_by_class_name("w60-0").text
        cttl = li.find_element_by_class_name("w14-0").text
        if nip == "" or not cttl.isdigit():
            continue
        nttl = float(cttl)
        if ttl == 0 or nttl < ttl:
            ttl = nttl
            ip = re.findall(r"\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b", nip)[0]

        if ip == "" or ttl == 0:
            print("获取ip失败")
            exit

        print("3. 发送请求更新IP地址")
        response = requests.get(
            "http://192.168.106.189:9090/ddns/request?domain="+addr+"&ip="+ip)
        if response.status_code == 200:
            print("4. 成功返回", ip, ttl)
        else:
            print("4. 处理失败:", response.status_code)
        browser.quit()
