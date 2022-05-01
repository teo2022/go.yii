package controller

import (
	"fmt"
	"go.yii/models"
	"go.yii/utils"
	"io/ioutil"
	"log"
	"strings"
	"unicode"
)

func GetStruct(list []models.ListCatalog, up []string) []models.GroupLine {
	var groupLine []models.GroupLine
	var books []string
	allCount := 0
	for _, v := range list {
		if !utils.SearchArrString(up, strings.ToLower(v.Catalog)) {
			continue
		}
		var listFile []models.File
		newGroup := models.GroupLine{
			Catalog:  v.Catalog,
			Patch:    v.Patch,
		}

		for _, r := range v.Files {
			if !strings.Contains(r, ".php") {
				continue
			}
			var lineList []models.Line
			input, err := ioutil.ReadFile(v.Patch + "/" + r)
			if err != nil {
				log.Fatalln(err)
			}
			lines := strings.Split(string(input), "\n")
			for i, line := range lines {
				if !utils.IsEngByLoop(line) {
					if strings.Contains(line, "@property") {
						continue
					}
					if strings.Contains(line, "\\Yii::t('app'") {
						continue
					}
					word := GetLine(line)
					if !utils.SearchArrString(books, word){
						books = append(books, word)
					}
					allCount = allCount + 1
					lineList = append(lineList, models.Line{
						Line: i,
						Old: line,
					})
				}
			}
			listFile = append(listFile, models.File{
				File:    r,
				ListLine: lineList,
			})
		}
		newGroup.ListFile = listFile
		groupLine = append(groupLine, newGroup)
	}

	listFile := "return [\n"
	for _, v := range books {
		n := strings.Replace(v, "\"", "", -1)
		n = strings.Replace(n, "'", "", -1)
		listFile = listFile + fmt.Sprintf("\t'%v' => 'text',\n", n)
	}
	listFile = listFile + "];\n"
	err2 := ioutil.WriteFile("books.txt", []byte(listFile), 0644)
	if err2 != nil {
		fmt.Println(err2)
	}
	models.SetBooks(books)
	return groupLine
}

func GetLine(line string) string {
	idDone := 0
	finString := ""
	lenLine := len(line)
	for i, v := range line {
		if v > unicode.MaxASCII {
			idDone = i
			break
		}
	}
	idStart := idDone
	isOne := false
	isTwo := false
	if idDone > 0 {
		for idStart < 0 {
			if string(line[idStart]) == "\"" {
				isTwo = true
				break
			}
			if string(line[idStart]) == "'" {
				isOne = true
				break
			}
			idStart = idStart - 1
		}
	}
	idFin := idStart + 1
	if idStart > 0 {
		for idFin < lenLine {
			if string(line[idFin]) == "\"" {
				isTwo = true
				break
			}
			if string(line[idFin]) == "'" {
				isOne = true
				break
			}
			idFin = idFin + 1
		}
	}
	finString = line[idStart:idFin]
	if isTwo {
		finString = fmt.Sprintf("\"%v\"", finString)
	}
	if isOne {
		finString = fmt.Sprintf("'%v'", finString)
	}
	return finString
}
