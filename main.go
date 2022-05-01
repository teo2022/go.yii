package main

import (
	"fmt"
	"go.yii/controller"
	"go.yii/models"
	"io/ioutil"
	"strings"
)

func main() {
	listUp := []string{"controllers", "views", "models", "data"}
	TeoStartGenerate("/Users/alex/Downloads/zh/bondus.mir-crm.loc",listUp)
}

func TeoStartGenerate(patch string, listUp []string) {
	AllFolder := controller.GetFolders(patch)
	listStructure := controller.GetStruct(AllFolder, listUp)
	gen := controller.GenerateTag(listStructure)
	fmt.Println(gen)
	//FinChangeFile(gen)
}

func FinChangeFile(list []models.GroupLine) {
	for _, group := range list {
		for _, line := range group.ListFile {
			var newFile string
			input, _ := ioutil.ReadFile(group.Patch + "/" + line.File)
			newFile = string(input)
			for _, v := range line.ListLine {
				newFile = strings.Replace(newFile, v.Old, v.New, -1)
			}
			err2 := ioutil.WriteFile(group.Patch + "/" + line.File, []byte(newFile), 0644)
			if err2 != nil {
				fmt.Println(err2)
			}
		}
	}
}
