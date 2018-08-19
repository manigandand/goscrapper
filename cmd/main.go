package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	content, err := ioutil.ReadFile("/home/manigandan/top_institutions.html")
	if err != nil {
		log.Println(err)
	}

	input := strings.NewReader(string(content))
	doc, err := goquery.NewDocumentFromReader(input)
	if err != nil {
		log.Println("could not scrap the page. err: ", err.Error())
		return
	}

	res := doc.Has("tbody")
	fmt.Printf("%+v\n\n", res)

	doc.Find("tr").Each(func(index int, item *goquery.Selection) {
		item.Children().Each(func(index int, item *goquery.Selection) {
			if index == 0 {
				fmt.Printf("%d ==> %+v\n\n", index, item.Text())
			}
		})
		// scriptText := item.Text()
		// fmt.Printf("%+v\n\n", scriptText)
	})

}
