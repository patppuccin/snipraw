package web

type Project struct {
	Name      string
	FileCount int
}

type FileNode struct {
	Name     string
	Path     string
	IsDir    bool
	Children []FileNode
}
