package models

type Line struct {
	Line int
	Old string
	New string
}
type Book struct {
	List []string
}
type File struct {
	File string
	ListLine []Line
}

type GroupLine struct {
	Catalog string
	Patch string
	ListFile []File
}

type ListModel struct {
	Name string
	Structure string
	Jsonschema string
}

type ListCatalog struct {
	Catalog string
	Patch string
	Files   []string
}

type GroupRoute struct {
	Name string
	Group []ListRoute
}
type ListRoute struct {
	Name string
	Link string
	Function string
	Folder string
	MetodFunc string
	Origin string
	RawBody string
}

var GlobalBooks Book

func SetBooks(v []string)  {
	GlobalBooks.List = v
}

func GetBooks () []string {
	return GlobalBooks.List
}