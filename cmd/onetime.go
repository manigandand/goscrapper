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
	"time"
)

var proxyIPs = []string{}

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func main() {
	// for i := 0; i < 30; i++ {
	// 	time.Sleep(10 * time.Millisecond)
	// 	go func() {
	// 		time.Sleep(30 * time.Millisecond)
	// 		fmt.Println(random(0, 400))
	// 	}()
	// }

	// time.Sleep(10 * time.Second)
	//

	content, err := ioutil.ReadFile("cmd/exp.txt")
	if err != nil {
		log.Println(err)
	}
	proxyIPLists := strings.Split(string(content), "\n")
	re := regexp.MustCompile("[#a-zA-Z]*")
	var (
		newProxyIPSet []string
		fileContent   []byte
	)

	for _, proxyIP := range proxyIPLists {
		httpProxyIP := fmt.Sprintf("http://%s", re.ReplaceAllString(proxyIP, ""))
		newProxyIPSet = append(newProxyIPSet, httpProxyIP)
		fileContent = append(fileContent, []byte(httpProxyIP)...)
		fileContent = append(fileContent, []byte("\n")...)
	}

	fErr := ioutil.WriteFile("cmd/test.txt", fileContent, 0644)
	if fErr != nil {
		log.Println(fErr)
		return
	}

	return
	// get free proxy address from here: https://free-proxy-list.net/
	//creating the proxyURL
	// proxyStr := "http://138.118.85.25:53281"
	for _, proxyStr := range proxyIPs {
		proxyURL, err := url.Parse(proxyStr)
		if err != nil {
			fmt.Printf("%+v\n\n", proxyStr)
			log.Println(err)
		}
		// creating proxy client
		proxyClient := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
			Timeout: 15 * time.Second,
		}
		//creating the URL to be loaded through the proxy
		// urlStr := "http://httpbin.org/get"
		urlStr := "http://api.openweathermap.org/data/2.5/weather?q=London"
		kaggleURL, err := url.Parse(urlStr)
		if err != nil {
			fmt.Println("=======================================")
			fmt.Printf("%+v\n\n", proxyStr)
			log.Println(err)
			fmt.Println("=======================================")
			fmt.Println()
			fmt.Println()
			// return
			continue
		}

		//generating the HTTP GET request
		request, err := http.NewRequest("GET", kaggleURL.String(), nil)
		if err != nil {
			fmt.Println("=======================================")
			fmt.Printf("%+v\n\n", proxyStr)
			log.Println(err)
			fmt.Println("=======================================")
			fmt.Println()
			fmt.Println()
			// return
			continue
		}
		// fmt.Printf("%+v\n\n", request)
		//calling the URL
		response, err := proxyClient.Do(request)
		if err != nil {
			fmt.Println("=======================================")
			fmt.Printf("%+v\n\n", proxyStr)
			log.Println(err)
			fmt.Println("=======================================")
			fmt.Println()
			fmt.Println()
			// return
			continue
		}
		//getting the response
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("=======================================")
			fmt.Printf("%+v\n\n", proxyStr)
			log.Println(err)
			fmt.Println("=======================================")
			fmt.Println()
			fmt.Println()
			// return
			continue
		}
		//printing the response
		log.Println(string(data))
		fmt.Printf("\n\n")

		// fmt.Printf("SUCCESS >>>>>: %+v\n\n", proxyStr)
	}
}
