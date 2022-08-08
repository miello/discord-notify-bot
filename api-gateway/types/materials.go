package types

type File struct {
	Title string `json:"title"`
	Href  string `json:"href"`
}

type MaterialView struct {
	FolderName string `json:"folderName"`
	File       []File `json:"file"`
}
