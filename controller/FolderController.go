package controller

import (
	"fmt"
	"go.yii/models"
	"io/ioutil"
	"strings"
)

func GetFolders(patch string) []models.ListCatalog {
	var finListCatalog []models.ListCatalog
	files, err := ioutil.ReadDir(patch)
	if err != nil {
		fmt.Println(err)
		return finListCatalog
	}
	var allListCatalog []models.ListCatalog
	lCatalog := models.ListCatalog{}
	lCatalog.Catalog = "/"
	lCatalog.Patch = patch + "/"
	for _, file := range files {
		if file.IsDir() {
			allListCatalog = append(allListCatalog, models.ListCatalog{Catalog: file.Name(), Patch: patch + "/"+ file.Name()})
		} else {
			lCatalog.Files = append(lCatalog.Files, file.Name())
		}
	}

	finListCatalog = RecursionFolder(finListCatalog, allListCatalog, patch)
	finListCatalog = append(finListCatalog, lCatalog)
	return finListCatalog
}

func RecursivFile(patch string, list []string) []string {
	files, err := ioutil.ReadDir(patch)
	if err != nil {
		fmt.Println(err)
		return list
	}
	for _, file := range files {
		if file.IsDir() {
			list = RecursivFile(patch  +"/"+ file.Name(),list)
		} else {
			list = append(list, patch + "/" + file.Name())
		}
	}
	return list
}
func RecursionFolder(finListCatalog []models.ListCatalog, allListCatalog []models.ListCatalog, patch string) []models.ListCatalog {
	count := 0
	//var newL []models.ListCatalog
	for _, v := range allListCatalog {
		count = count + 1
		patch2 := v.Patch
		if strings.Contains(v.Catalog, ".") {
			continue
		}
		files, err := ioutil.ReadDir(patch2)
		if err != nil {
			fmt.Println(err)
			break
		}
		lCatalog := models.ListCatalog{}
		lCatalog.Catalog = v.Catalog
		lCatalog.Patch = patch2
		for _, file := range files {
			if file.IsDir() {
				lCatalog.Files = RecursivFile(patch + "/" + v.Catalog +"/"+ file.Name(), lCatalog.Files)
				//newL = append(newL, models.ListCatalog{Catalog: file.Name(),Patch: patch + "/" + v.Catalog +"/"+ file.Name()})
			} else {
				lCatalog.Files = append(lCatalog.Files, patch + "/" + v.Catalog +"/"+ file.Name())
			}
		}
		finListCatalog = append(finListCatalog, lCatalog)
	}
	//if len(newL) != len(allListCatalog) {
	//	finListCatalog = RecursionFolder(finListCatalog, newL, patch)
	//}
	//fmt.Println(count)
	return finListCatalog
}

