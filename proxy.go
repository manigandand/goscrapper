package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

type BlockedProxyIP struct {
	// Index     int       `json:"index"`
	ProxyIP   string    `json:"proxy_ip"`
	BlockedAt time.Time `json:"blocked_at"`
}

var (
	BlockedProxyIPs []*BlockedProxyIP
	ProxyIPLists    []string
	ProxyLock       sync.RWMutex
)

func bolckProxyIP(proxyStr string) {
	BlockedProxyIPs = append(BlockedProxyIPs, &BlockedProxyIP{
		ProxyIP:   proxyStr,
		BlockedAt: time.Now(),
	})
}

func isProxyIPBlocked(proxyIP string) bool {
	if len(BlockedProxyIPs) > 0 {
		// fmt.Println("--------BlockedProxyIPs:", len(BlockedProxyIPs), " ----------")
	}
	for _, pip := range BlockedProxyIPs {
		if pip.ProxyIP == proxyIP {
			elapsed := time.Since(pip.BlockedAt)
			if elapsed.Minutes() > 15 {
				return false
			}
			return true
		}
	}

	return false
}

func goProxyJob() {
	i := 2
	for {
		time.Sleep(1 * time.Minute)
		getNewProxtList(i)
		i++
	}
}

func getNewProxtList(nthTime int) {
	var (
		newProxyIPSet []string
		fileContent   []byte
	)
	fmt.Println("==================================================")
	fmt.Println("Calling for new proxy list ", nthTime, " time")
	re := regexp.MustCompile("[#a-zA-Z]*")
	proxyListAPI := "http://list.didsoft.com/get?email=mail@manigandan.com&pass=mpcnhr&pid=httppremium&https=yes"
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	request, err := http.NewRequest("GET", proxyListAPI, nil)
	if err != nil {
		log.Println(err)
		return
	}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return
	}

	proxyIPLists := strings.Split(string(data), "\n")
	total := len(proxyIPLists)
	for i, proxyIP := range proxyIPLists {
		if ts(proxyIP) != "" {
			httpProxyIP := fmt.Sprintf("http://%s", re.ReplaceAllString(proxyIP, ""))
			newProxyIPSet = append(newProxyIPSet, httpProxyIP)
			fileContent = append(fileContent, []byte(httpProxyIP)...)
			if i+1 != total {
				fileContent = append(fileContent, []byte("\n")...)
			}
		}
	}

	filePath := fmt.Sprintf("proxy_ips_list/proxt_list_%d_%d.txt", nthTime, time.Now().Unix())
	fErr := ioutil.WriteFile(filePath, fileContent, 0644)
	if fErr != nil {
		log.Println(fErr)
		return
	}
	fmt.Println("we got ", len(newProxyIPSet), " new proxy ips")

	ProxyLock.Lock()
	proxyIPs = append(proxyIPs, newProxyIPSet...)
	ProxyLock.Unlock()

	fmt.Println("total proxy we have ", len(proxyIPs))
	fmt.Println("==================================================")

	return
}

func getProxyCLient() (*http.Client, string) {
	var (
		proxyURL *url.URL
		err      error
	)

	time.Sleep(30 * time.Millisecond)
	for {
		index := random(0, len(proxyIPs))
		ProxyLock.RLock()
		proxyStr := proxyIPs[index]
		ProxyLock.RUnlock()

		if !isProxyIPBlocked(proxyStr) {
			proxyURL, err = url.Parse(proxyStr)
			if err != nil {
				bolckProxyIP(proxyStr)
				log.Println("Index: ", index, " Proxy ip error: >>>>", proxyURL, " >>>>>>", err.Error())
				continue
			}

			// break here, we got new proxy ip
			break
		}
	}

	// creating proxy client
	proxyClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 20 * time.Second,
	}

	return proxyClient, proxyURL.String()
}

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
