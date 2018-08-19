package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"code.aircto.in/root/vendor_old/github.com/jinzhu/gorm"
	_ "code.aircto.in/root/vendor_old/github.com/lib/pq"
	_ "code.aircto.in/root/vendor_old/github.com/mattn/go-sqlite3"
	"github.com/PuerkitoBio/goquery"
)

var (
	path       = flag.String("path", "", "user data csv file path")
	actionType = flag.String("type", "", "type of operation")
	ts         = strings.TrimSpace
	db         *gorm.DB
	dErr       error
)

func init() {
	flag.Parse()
	db, dErr = gorm.Open("postgres", "user=ubuntu password=ubuntu dbname=kaggle sslmode=disable")
	if dErr != nil {
		log.Fatal(dErr)
	}
}

func goscrapper(i int, ku *KaggleUser, wg *sync.WaitGroup, done chan bool, errChan chan *KaggleErr) {
	var (
		proxyClient *http.Client
		proxyURL    string
		isError     bool
	)
	time.Sleep(40 * time.Second)
	proxyClient, proxyURL = getProxyCLient()

	defer func() {
		wg.Done()
		done <- true
		errChan <- &KaggleErr{
			ProxyIP: proxyURL,
			IsError: isError,
		}
	}()

	kaggleURL := fmt.Sprintf("https://www.kaggle.com/%s", ku.UserName)
	fmt.Printf("Scarpping: %d >> %s >> using proxy ip: %s\n\n", i, kaggleURL, proxyURL)
	fmt.Println("--------BlockedProxyIPs:", len(BlockedProxyIPs), " ----------")

	//generating the HTTP GET request
	request, err := http.NewRequest("GET", kaggleURL, nil)
	if err != nil {
		// bolckProxyIP(proxyURL)
		log.Println(err)
		isError = true
		return
	}
	//calling the URL
	response, err := proxyClient.Do(request)
	if err != nil {
		// bolckProxyIP(proxyURL)
		log.Println(err)
		isError = true
		return
	}
	request.Close = true

	//getting the response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		// bolckProxyIP(proxyURL)
		log.Println(err)
		isError = true
		return
	}
	defer response.Body.Close()
	fmt.Println(string(body))
	ip := strings.NewReader(string(body))
	doc, err := goquery.NewDocumentFromReader(ip)
	if err != nil {
		// bolckProxyIP(proxyURL)
		log.Println("could not scrap the page. err: ", err.Error())
		isError = true
		return
	}

	res := scrapKaggleUserData(doc)
	if res != nil {
		fmt.Printf("%+v\n\n", res)
		if res.UserID == 0 {
			bolckProxyIP(proxyURL)
			log.Println("could not scrap the page. Might Be blocked us.")
			isError = true
			return
			// fmt.Println("wait for 5 mins")
			// // wait for 5 mins
			// time.Sleep(5 * time.Minute)
			// // recursive call
			// kaggleScrapper()
		}

		if ts(res.DisplayName) == "" {
			res.DisplayName = ku.DisplayName
		}

		res.Country = ts(res.Country)
		JSONStr, _ := json.Marshal(res)
		if strings.ToLower(res.Country) == "india" {
			iuFilePath := fmt.Sprintf("indian_users/%d_%s.json", ku.Id, ku.UserName)
			fErr := ioutil.WriteFile(iuFilePath, JSONStr, 0644)
			if fErr != nil {
				log.Println("can't able to create india json file.", fErr.Error())
				isError = true
				return
			}
			ku.IsVisited = true
			db.Save(ku)
			fmt.Printf("Indian User Data: %d-%s-%s\n\n\n", ku.Id, ku.UserName, ku.DisplayName)
			isError = false
			return
		}

		if res.Country == "" {
			niuFilePath := fmt.Sprintf("non_indian_users/%d_%s.json", ku.Id, ku.UserName)
			fErr := ioutil.WriteFile(niuFilePath, JSONStr, 0644)
			if fErr != nil {
				log.Println("can't able to create json file.", fErr.Error())
				isError = true
				return
			}
			ku.IsVisited = true
			db.Save(ku)
			fmt.Printf("Nil Country User Data: %d-%s-%s\n\n\n", ku.Id, ku.UserName, ku.DisplayName)
			isError = false
			return
		}
		niuFilePath := fmt.Sprintf("other_country_users/%d_%s.json", ku.Id, ku.UserName)
		fErr := ioutil.WriteFile(niuFilePath, JSONStr, 0644)
		if fErr != nil {
			log.Println("can't able to create json file.", fErr.Error())
			isError = true
			return
		}

		ku.IsVisited = true
		db.Save(ku)
		fmt.Printf("%s Country User Data: %d-%s-%s\n\n\n", res.Country, ku.Id, ku.UserName, ku.DisplayName)
		isError = false
		return
	}

	fmt.Println("deferring.......... ")
	isError = false
	return
}

func kaggleScrapper() {
	for interval := 0; interval < 4; interval++ {
		start := time.Now()
		var (
			wg          sync.WaitGroup
			kaggleUsers []*KaggleUser
		)

		maxNbConcurrentGoroutines := 20
		concurrentGoroutines := make(chan struct{}, maxNbConcurrentGoroutines)
		for i := 0; i < maxNbConcurrentGoroutines; i++ {
			concurrentGoroutines <- struct{}{}
		}
		done := make(chan bool)
		errChan := make(chan *KaggleErr)

		db.Limit(1000000).
			Find(&kaggleUsers, "is_visited=?", false)

		totalUsers := len(kaggleUsers)
		fmt.Println("totalUsers: ", totalUsers)
		// Collect all the jobs, and since the job is finished, we can
		// release another spot for a goroutine.
		go func() {
			totalDone := 0
			for i := 0; i < totalUsers*2; i++ {
				select {
				case isDone := <-done:
					if isDone {
						fmt.Println("received done signal....")
						fmt.Printf("INTERVAL: %d ==> Processed %d out of %d.\n", interval+1, totalDone+1, totalUsers)
						// Say that another goroutine can now start.
						concurrentGoroutines <- struct{}{}
						totalDone++
					}
				case isErr := <-errChan:
					if isErr.IsError {
						fmt.Println("==========================================")
						fmt.Println("==========================================")
						fmt.Println("::::::::::::::::::::::::::::::::::::::::::")
						fmt.Printf("ERROR HAPPENED: Proxy IP: %s\nIsError: %+v\n", isErr.ProxyIP, isErr.IsError)
						fmt.Printf("You many need to remove this ip from the list\n")
						fmt.Println("::::::::::::::::::::::::::::::::::::::::::")
						fmt.Println("==========================================")
						fmt.Println("==========================================")
						// fmt.Println("wait for 5 mins")
						// log.Fatal("could not scrap the page. Might Be blocked us.")
						// wait for 5 mins
						// time.Sleep(5 * time.Minute)
						// recursive call
						// kaggleScrapper()
					}
				}
			}
		}()

		wg.Add(totalUsers)
		for i, ku := range kaggleUsers {
			i++
			time.Sleep(30 * time.Millisecond)
			// intervals(i)
			<-concurrentGoroutines
			go goscrapper(i, ku, &wg, done, errChan)
			/*
				doc, err := goquery.NewDocument(kaggleURL)
				if err != nil {
					log.Fatal("could not scrap the page. err: ", err.Error())
				}
			*/
			// str, err := doc.Html()
			// fmt.Printf("%+v \n\n\n%+v \n", str, err)
			// fmt.Printf("\nTITLE: %s\n", doc.Find("title").Contents().Text())
		}

		// close the channel when all goroutine initiated
		wg.Wait()
		close(errChan)
		close(done)
		// fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<")
		// fmt.Printf("There are total %d indian users data found\n\n.", indianUsers)

		// fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<")
		// fmt.Printf("There are total %d users not updated there locations.\n\n.", nilCountryUsers)
		elapsed := time.Since(start)
		log.Println("Scapping Script took: ", elapsed)

		// give 3 mins gap for each interval
		time.Sleep(30 * time.Millisecond)
	}

	return
}

func main() {
	start := time.Now()
	if ts(*actionType) == "" {
		log.Fatalln("actionType not provided. Use -type dumb_db or scrap_data.")
	}

	switch *actionType {
	case "excel_report":
		excelReport()
	case "dumb_db":
		// if err := dumbCSV(); err != nil {
		// 	log.Fatalln(err.Error())
		// }
		// return
	case "scrap_data":
		fmt.Println(">>>>>>>>>> Scarpping Kaggle Data <<<<<<<<<<<<<<<<<<<<<<<<")
		getNewProxtList(1)
		loadProxyIPList()
		go goProxyJob()
		kaggleScrapper()
	}

	elapsed := time.Since(start)
	log.Println("Scapping Script took: ", elapsed)

	return
}

func scrapKaggleUserData(doc *goquery.Document) *Kaggle {
	res := new(Kaggle)
	doc.Find("script").Each(func(index int, item *goquery.Selection) {
		scriptText := item.Text()
		if isContainKaggleData(scriptText) && ts(scriptText) != "" {
			trimmed := ts(splitKaggleJSON(scriptText))
			fmt.Printf("JSON String >>>>: \n%+v\n ", trimmed)
			if ts(trimmed) != "" && len(trimmed) > 25 {
				err := json.Unmarshal([]byte(trimmed), res)
				if err != nil {
					// Fatal?
					fmt.Println(err)
				}
				// fmt.Printf("kaggle: \n%+v\n ", res)
				return
			}
			// fmt.Printf("Index: %d\n Data:\n%+v\n\n", index, scriptText)
		}
	})

	return res
}

func splitKaggleJSON(inputStr string) string {
	/*
		var res string
		str1 := strings.Split(inputStr, ");")
		if len(str1) > 0 {
			strs := strings.Split(str1[0], "];")
			if len(strs) > 0 {
				if len(strs[1]) > 30 {
					res = parseKaggleDataString(strs[1])
				}
			}
		}
	*/

	var substitution = ``
	var replace1 = `var Kaggle=window.Kaggle||{};Kaggle.State=Kaggle.State||[];Kaggle.State.push(`
	var replace2 = `);performance && performance.mark && performance.mark("ProfileContainerReact.componentCouldBootstrap");`
	JSONStr := strings.Replace(inputStr, replace1, substitution, -1)
	JSONStr = strings.Replace(JSONStr, replace2, substitution, -1)
	if len(JSONStr) > 30 {
		return JSONStr
	}

	return ""
}

func isContainKaggleData(inputStr string) bool {
	return (strings.Contains(inputStr, "Kaggle.State.push") &&
		strings.Contains(inputStr, `performance.mark("ProfileContainerReact.componentCouldBootstrap")`))
}

func parseKaggleDataString(inputStr string) string {
	if !isContainKaggleData(inputStr) {
		return ""
	}
	// fmt.Printf("Kaggle.State.push(string): ===> \n %s \n\n", inputStr)

	var re = regexp.MustCompile(`(?im)Kaggle.State.push[(]`)
	var substitution = ``

	return re.ReplaceAllString(inputStr, substitution)
}

func parseKaggleData(inputStr string) string {
	if !isContainKaggleData(inputStr) {
		return ""
	}

	var re = regexp.MustCompile(`(?m)Kaggle.State.push`)
	var str = []byte(inputStr)
	var substitution = []byte(``)
	var i = 0
	var count = 1

	str = re.ReplaceAllFunc(str, func(s []byte) []byte {
		if count < 0 {
			return substitution
		} else if i < count {
			i++
			return substitution
		}

		return s
	})

	return string(str)
}

func basicScrapper() {
	doc, err := goquery.NewDocument("https://www.kaggle.com/manigandanlaunchyard")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%+v\n\n", doc)
	str, err := doc.Html()
	fmt.Printf("%+v \n\n\n%+v \n", str, err)

	// doc.Find("#main article .entry-title").Each(func(index int, item *goquery.Selection) {
	// 	title := item.Text()
	// 	linkTag := item.Find("a")
	// 	link, _ := linkTag.Attr("href")
	// 	fmt.Printf("Post #%d: %s - %s\n", index, title, link)
	// })

	doc.Find("body a").Each(func(index int, item *goquery.Selection) {
		linkTag := item
		link, _ := linkTag.Attr("href")
		scriptText := linkTag.Text()
		fmt.Printf("Link #%d: '%s' - '%s'\n", index, scriptText, link)
	})

	fmt.Printf("\nTITLE: %s\n", doc.Find("title").Contents().Text())
}

func dumbCSV() error {
	if !db.HasTable(&KaggleUser{}) {
		if err := db.CreateTable(&KaggleUser{}).Error; err != nil {
			log.Fatalf("critical.KaggleUser.migrate.create_table.%s", err)
		}
	}

	if *path == "" {
		return errors.New("file not provided. Use -path filepath")
	}

	csvFile, csvErr := os.Open(*path)
	if csvErr != nil {
		log.Fatalln(csvErr.Error())
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		if line[1] == "UserName" {
			continue
		}

		userID, _ := strconv.Atoi(ts(line[0]))
		userName := ts(line[1])
		displayName := ts(line[2])
		registerDate := ts(line[3])
		performanceTier := ts(line[4])

		user := &KaggleUser{
			Id:              userID,
			UserName:        userName,
			DisplayName:     displayName,
			RegisterDate:    registerDate,
			PerformanceTier: performanceTier,
			IsVisited:       false,
		}

		if err := db.Save(user).Error; err != nil {
			return fmt.Errorf("can't able to save userdata: user:\n %+v", user)
		}
	}

	return nil
}

func intervals(i int) {
	// each 5 users give extra 5 sec sleep
	if (i % 5) == 0 {
		fmt.Printf("providing extra %d sec sleep for the %dth user\n", 5, i)
		time.Sleep(5 * time.Second)
	}
	// each 10 users give extra 10 sec sleep
	if (i % 10) == 0 {
		fmt.Printf("providing extra %d sec sleep for the %dth user\n", 10, i)
		time.Sleep(10 * time.Second)
	}
	// each 30 users give extra 15 sec sleep
	if (i % 30) == 0 {
		fmt.Printf("providing extra %d sec sleep for the %dth user\n", 15, i)
		time.Sleep(15 * time.Second)
	}
	// each 50 users give extra 25 sec sleep
	if (i % 50) == 0 {
		fmt.Printf("providing extra %d sec sleep for the %dth user\n", 25, i)
		time.Sleep(25 * time.Second)
	}
	// each 100 users give extra 1 min sleep
	if (i % 100) == 0 {
		fmt.Printf("providing extra %d min sleep for the %dth user\n", 1, i)
		time.Sleep(1 * time.Minute)
	}
}

/*
	ip := strings.NewReader(Input)
	doc, err := goquery.NewDocumentFromReader(ip)
	if err != nil {
		log.Fatal("could not scrap the page. err: ", err.Error())
	}
*/
