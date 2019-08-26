from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.action_chains import ActionChains
from selenium.webdriver.common.keys import Keys
import re
import requests

chrome_options = webdriver.ChromeOptions()
chrome_options.add_argument('--headless')  # 无界面选项

print("----------------查询最快的github.com DNS IP--------------------")

print("1. 从tool.chinaz.com获取github.com所有可用的IP")
browser = webdriver.Chrome(chrome_options=chrome_options)
browser.get("http://tool.chinaz.com/dns?type=1&host=github.com")
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
    "http://localhost:9090/ddns/request?domain=github.com&ip="+ip)
if response.status_code == 200:
    print("4. 成功返回", ip, ttl)
else:
    print("4. 处理失败:", response.status_code)
browser.quit()
