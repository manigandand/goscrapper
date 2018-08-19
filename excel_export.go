package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/tealeg/xlsx"
)

func excelReport() {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	files, err := dirWalk("/home/manigandan/Projects/aircto/go/src/goscrapper/indian_users_25")
	if err != nil {
		panic(err)
	}

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "Name"
	cell = row.AddCell()
	cell.Value = "Kaggle URL" // userUrl
	cell = row.AddCell()
	cell.Value = "Country"
	cell = row.AddCell()
	cell.Value = "City"
	cell = row.AddCell()
	cell.Value = "GitHub UserName"
	cell = row.AddCell()
	cell.Value = "Twitter UserName"
	cell = row.AddCell()
	cell.Value = "linkedIn Url"
	cell = row.AddCell()
	cell.Value = "website Url"

	for _, file := range files {
		res := new(Kaggle)
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Println(err)
		}
		err = json.Unmarshal(content, res)
		if err != nil {
			fmt.Println(err)
		}

		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = res.DisplayName
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("https://kaggle.com%s", res.UserURL)
		cell = row.AddCell()
		cell.Value = res.Country
		cell = row.AddCell()
		cell.Value = res.City
		cell = row.AddCell()
		cell.Value = res.GitHubUserName
		cell = row.AddCell()
		cell.Value = res.GitHubUserName
		cell = row.AddCell()
		cell.Value = res.LinkedInURL
		cell = row.AddCell()
		cell.Value = res.WebsiteURL

		continue
	}

	err = file.Save("kaggle_data.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}

}
