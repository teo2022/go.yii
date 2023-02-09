package controller

import (
	"fmt"
	"go.yii/models"
	"go.yii/utils"
	"io/ioutil"
	"regexp"
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
			Catalog: v.Catalog,
			Patch:   v.Patch,
		}

		for _, r := range v.Files {
			if !strings.Contains(r, ".php") {
				continue
			}
			var lineList []models.Line
			input, err := ioutil.ReadFile(r)
			if err != nil {
				fmt.Println(err)
				continue
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
					if strings.Contains(line, "Yii::t('app'") {
						continue
					}
					if strings.Contains(line, "const") {
						continue
					}
					_, _, word := GetLine(line)
					if !utils.SearchArrString(books, word) {
						books = append(books, word)
					}
					allCount = allCount + 1
					lineList = append(lineList, models.Line{
						Line: i,
						Old:  line,
					})
				}
			}
			listFile = append(listFile, models.File{
				File:     r,
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
		listFile = listFile + fmt.Sprintf("\t'%v' => '',\n", n)
	}
	listFile = listFile + "];\n"
	err2 := ioutil.WriteFile("books.txt", []byte(listFile), 0644)
	if err2 != nil {
		fmt.Println(err2)
	}
	models.SetBooks(books)
	return groupLine
}

func GetLine(line string) (bool, bool, string) {
	idDone := 0
	finString := ""
	lenLine := len(line) - 1
	isStartTwo, isStartOne := SetCountBaet(line)
	for i, v := range line {
		if v > unicode.MaxASCII {
			idDone = i
			break
		}
	}
	idStart := idDone
	idFin := idStart + 1
	if idStart > 0 {
		for idFin <= lenLine {
			if idFin+1 <= lenLine && idFin+2 <= lenLine {
				searPhp := fmt.Sprintf("%v%v%v", string(line[idFin]), string(line[idFin+1]), string(line[idFin+2]))
				if searPhp == "<?=" {
					break
				}
			}
			n := string(line[idFin])
			if string(n) == "\"" && isStartTwo {
				break
			}
			if string(line[idFin]) == "'" && isStartOne {
				break
			}
			idFin = idFin + 1
		}
	}

	finString = string([]rune(line[idStart:idFin]))
	finString = strings.Replace(finString, "\n", "", -1)
	finString = strings.Replace(finString, "\r", "", -1)
	finString = strings.Replace(finString, "</span>", "", -1)
	finString = strings.Replace(finString, "</h4>", "", -1)
	finString = strings.Replace(finString, "</th>", "", -1)
	finString = strings.Replace(finString, "</strong>", "", -1)
	finString = strings.Replace(finString, "</div>", "", -1)
	finString = strings.Replace(finString, "</p>", "", -1)
	finString = strings.Replace(finString, "</a>", "", -1)
	finString = strings.Replace(finString, "</b>", "", -1)
	finString = strings.Replace(finString, "</label>", "", -1)
	finString = strings.Replace(finString, "</td>", "", -1)
	finString = strings.Replace(finString, "</title>", "", -1)
	finString = strings.Replace(finString, "</h2>", "", -1)
	finString = strings.Replace(finString, "<i class=fa fa-download></i>", "", -1)
	finString = strings.Replace(finString, "</i>", "", -1)
	finString = strings.Replace(finString, "<i class=fa fa-download>", "", -1)
	finString = strings.Replace(finString, "</h5>", "", -1)
	finString = strings.Replace(finString, "<br>", "", -1)
	finString = strings.Replace(finString, "<div class=col-md-6 style=text-align: right;>", "", -1)

	return isStartOne, isStartTwo, finString
}

func SetCountBaet(line string) (bool, bool) {
	countOne := 0
	countTwo := 0
	isOne := false
	isTwo := false
	isStartTwo := false
	isStartOne := false
	re := regexp.MustCompile("[А-Яа-я]+?") //проверяем на киррилические символы
	for _, v := range line {
		//b := string(v)
		//fmt.Println(b)
		if string(v) == "\"" && isTwo == false {
			isTwo = true
			countTwo = countTwo + 1
			continue
		}
		if string(v) == "\"" && isTwo == true {
			isTwo = false
			countTwo = countTwo - 1
			continue
		}
		if string(v) == "'" && isOne == false {
			isOne = true
			countOne = countOne + 1
			continue
		}
		if string(v) == "'" && isOne == true {
			isOne = false
			countOne = countOne - 1
			continue
		}
		isRussian := re.MatchString(string(v))
		if isRussian {
			if isOne == false && isTwo == true {
				isStartTwo = true
			}
			if isOne == true && isTwo == false {
				isStartOne = true
			}
			break
		}
	}
	return isStartTwo, isStartOne
}
