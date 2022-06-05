package controller

import (
	"fmt"
	"go.yii/models"
	"io/ioutil"
	"log"
	"strings"
)

func GenerateTag(list []models.GroupLine) []models.GroupLine {
	listFile := ""
	noRep := 0
	newRep := 0
	books := models.GlobalBooks
	n1 := 0
	for i1, group := range list {
		for i2, line := range group.ListFile {
			for i3, v := range line.ListLine {
				input, err := ioutil.ReadFile(line.File)
				if err != nil {
					log.Fatalln(err)
				}
				linesPhp := strings.Split(string(input), "\n")
				isOpenPhp := false
				countOpenPhp := 0
				dopSearch := false
				for i, l := range linesPhp {
					if strings.Contains(l, "<?php") {
						isOpenPhp = true
						countOpenPhp += 1
					}
					if strings.Contains(l, "<?=") {
						isOpenPhp = true
						countOpenPhp += 1
					}
					if isOpenPhp == true && strings.Contains(l, "?>") {
						countOpenPhp -= 1
						isOpenPhp = false
					}
					if i == v.Line {
						if countOpenPhp > 0 {
							isOpenPhp = true
							dopSearch = true
						} else {
							if strings.Contains(l, "<?=") || strings.Contains(l, "<?php") {
								dopSearch = true
							}
						}
						break
					}
				}

				n1 = n1 + 1
				newL := ReplaceTag(v.Old, isOpenPhp, dopSearch)
				if newL == v.Old {
					noRep = noRep + 1
				} else {
					newRep = newRep + 1
				}
				list[i1].ListFile[i2].ListLine[i3].New = newL
				listFile = listFile + fmt.Sprintf("%v  ->  %v \n\n", strings.Trim(v.Old, " "), strings.Trim(newL, " "))
			}
		}
	}
	for _, book := range books.List {
		listFile = listFile + fmt.Sprintf("%v \n", book)
	}
	err2 := ioutil.WriteFile("finModel.txt", []byte(listFile), 0644)
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Printf("Вставлено строк будет:%v   Остались без изменений:%v \n", newRep, noRep)
	return list
}

func ReplaceTag(value string, php bool, dopSearch bool) string {
	var newL string
	//if strings.Contains(value," Html::a('Добавить"){
	//	fmt.Println("das")
	//}
	isStartOne, isStartTwo, searchBlock := GetLine(value)
	if dopSearch {
		if SearchClosePhp(value, searchBlock) {
			php = false
		} else {
			php = true
		}
	}
	if !php {
		newL = strings.Replace(value, searchBlock, fmt.Sprintf("<?=\\Yii::t('app', '%v')?>", searchBlock), -1)
		return newL
	}
	if isStartTwo {
		newL = strings.Replace(value, searchBlock, fmt.Sprintf("\".\\Yii::t('app', '%v').\"", searchBlock), -1)
		newL = SearchTwoCharts(newL)
	}
	if isStartOne {
		newL = strings.Replace(value, searchBlock, fmt.Sprintf("'.\\Yii::t('app', '%v').'", searchBlock), -1)
		newL = SearchOneCharts(newL)

	}
	if !isStartTwo && !isStartOne {
		newL = strings.Replace(value, searchBlock, fmt.Sprintf("\\Yii::t('app', %v)", searchBlock), -1)
	}
	return newL
}

func SearchOneCharts(newL string) string {
	result := ""
	isLeft := false
	countLeft := 0
	isRight := false
	countRight := 0
	for i, v := range newL {
		if len(newL) > i+2 {
			if !isLeft && string(v) == "'" && string(newL[i+1]) == "'" && string(newL[i+2]) == "." {
				isLeft = true
				countLeft = 3
			}
		}
		if isLeft && countLeft > 0 {
			countLeft -= 1
			continue
		}
		if countLeft == 0 {
			isLeft = false
		}
		if len(newL) > i+2 {
			if string(v) == "." && string(newL[i+1]) == "'" && string(newL[i+2]) == "'" {
				isRight = true
				countRight = 3
			}
		}
		if isRight && countRight > 0 {
			countRight -= 1
			continue
		}
		if countRight == 0 {
			isRight = false
		}
		result += string(v)
	}
	return result
}

func SearchTwoCharts(newL string) string {
	result := ""
	isLeft := false
	countLeft := 0
	isRight := false
	countRight := 0
	for i, v := range newL {
		if len(newL) > i+2 {
			if string(v) == "\"" && string(newL[i+1]) == "\"" && string(newL[i+2]) == "." {
				isLeft = true
				countLeft = 3
			}
		}
		if isLeft && countLeft > 0 {
			countLeft -= 1
			continue
		}
		if countLeft == 0 {
			isLeft = false
		}
		if len(newL) > i+2 {
			if string(v) == "." && string(newL[i+1]) == "\"" && string(newL[i+2]) == "\"" {
				isRight = true
				countRight = 3
			}
		}
		if isRight && countRight > 0 {
			countRight -= 1
			continue
		}
		if countRight == 0 {
			isRight = false
		}
		result += string(v)
	}
	return result
}

func SearchClosePhp(value string, searchBlock string) bool {
	newL := strings.Replace(value, searchBlock, "{{}}", -1)
	idStart := 0
	isClouse := false
	for i, v := range newL {
		if string(v) == "{" && string(newL[i+1]) == "{" && string(newL[i+2]) == "}" && string(newL[i+3]) == "}" {
			idStart = i
			break
		}
	}
	idFin := idStart
	for idFin != 0 {
		if string(newL[idFin]) == ">" && string(newL[idFin-1]) == "?" {
			return true
		}
		idFin -= 1
	}
	return isClouse
}
