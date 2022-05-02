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
				for i, l := range linesPhp {
					if strings.Contains(l, "<?php") {
						isOpenPhp = true
					}
					if strings.Contains(l, "<?=") {
						isOpenPhp = true
					}
					if isOpenPhp == true && strings.Contains(l, "?>") {
						isOpenPhp = false
					}
					if i == v.Line {
						break
					}
				}

				n1 = n1 + 1
				newL := ReplaceTag(v.Old, isOpenPhp)
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

func ReplaceTag(value string, php bool) string {
	var newL string
	//if strings.Contains(value," Html::a('Добавить"){
	//	fmt.Println("das")
	//}
	isStartOne, isStartTwo, searchBlock := GetLine(value)
	if !php {
		newL = strings.Replace(value, searchBlock, fmt.Sprintf("<?=\\Yii::t('app', '%v')?>\"",searchBlock), -1)
		return newL
	}
	if isStartTwo {
		newL = strings.Replace(value, searchBlock, fmt.Sprintf("\".\\Yii::t('app', '%v').\"",searchBlock), -1)
	}
	if isStartOne {
		newL = strings.Replace(value, searchBlock, fmt.Sprintf("'.\\Yii::t('app', '%v').'",searchBlock), -1)
	}
	if !isStartTwo && !isStartOne {
		newL = strings.Replace(value, searchBlock, fmt.Sprintf("\\Yii::t('app', %v)",searchBlock), -1)
	}
	return newL
}
