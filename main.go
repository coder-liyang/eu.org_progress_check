package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Domain struct {
	Domain  string `json:"domain"`
	Created string `json:"created"`
	Updated string `json:"updated"`
	DNSSEC  string `json:"dnssec"`
	Flags   string `json:"flags"`
}

func euOrg() {
	var allDomain []Domain
	url := "https://nic.eu.org/arf/en/"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	cookie, err := readCookieFromFile("cookie.txt")
	if err != nil {
		log.Fatalf("Error reading cookie from file: %v", err)
	}

	if len(cookie) < 50 {
		log.Fatalf("正确的cookie大约86个字符")
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	//req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Cookie", cookie)
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "none")
	req.Header.Add("Sec-Fetch-User", "?1")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	//{
	//	body, err := ioutil.ReadAll(res.Body)
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	fmt.Println(string(body))
	//}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	if doc.Find("#id_loginform").Length() > 0 {
		log.Fatal("cookie失效")
	}

	// Find the review items
	domains := doc.Find(".domainlist tr").Not(":first-child")
	domains.Each(func(i int, s *goquery.Selection) {
		tds := s.Children()
		domain := Domain{}
		tds.Each(func(i int, selection *goquery.Selection) {
			text := strings.Trim(selection.Text(), "\n ")
			switch i {
			case 0:
				domain.Domain = text
			case 1:
				domain.Created = text
			case 2:
				domain.Updated = text
			case 3:
				domain.DNSSEC = text
			case 4:
				domain.Flags = text
			}
		})
		allDomain = append(allDomain, domain)
	})
	fmt.Println(allDomain)
	fmt.Println("successfully")
}

// 从本地文件读取cookie
func readCookieFromFile(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 使用bufio读取文本行
	var cookieValue string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cookieValue = scanner.Text()
		break // 假设cookie在文件的第一行
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return cookieValue, nil
}
func main() {
	euOrg()
}
