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
			allListCatalog = append(allListCatalog, models.ListCatalog{Catalog: file.Name()})
		} else {
			lCatalog.Files = append(lCatalog.Files, file.Name())
		}
	}

	finListCatalog = RecursionFolder(finListCatalog, allListCatalog, patch)
	finListCatalog = append(finListCatalog, lCatalog)
	return finListCatalog
}

func RecursionFolder(finListCatalog []models.ListCatalog, allListCatalog []models.ListCatalog, patch string) []models.ListCatalog {
	for _, v := range allListCatalog {
		patch2 := patch + "/" + v.Catalog
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
				allListCatalog = append(allListCatalog, models.ListCatalog{Catalog: file.Name()})
			} else {
				lCatalog.Files = append(lCatalog.Files, file.Name())
			}
		}
		finListCatalog = append(finListCatalog, lCatalog)
	}
	return finListCatalog
}
