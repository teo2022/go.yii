package controller

import (
	"fmt"
	"go.yii/models"
	"io/ioutil"
	"strings"
)

func GenerateTag(list []models.GroupLine) []models.GroupLine {
	listFile := ""
	noRep := 0
	books := models.GlobalBooks
	for i1, group := range list {
		for i2, line := range group.ListFile {
			for i3, v := range line.ListLine {
				newL := ReplaceTag(v.Old)
				if newL == v.Old {
					noRep = noRep + 1
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
	return list
}

func ReplaceTag(value string) string {
	books := models.GlobalBooks
	for _, book := range books.List {
		//book = fmt.Sprintf("\"%v\"",book)
		if strings.Contains(value, book) {
			newL := strings.Replace(value, book, fmt.Sprintf("\\Yii::t('app', %v)",book), -1)
			return newL
		}
	}

	return value
}