package files

type bindPath struct {
	Path string `json:"path" form:"path" binding:"required"`
}

type bindList struct {
	Path string `form:"path"`
}

type bindRename struct {
	OldPath string `json:"oldPath" binding:"required"`
	NewPath string `json:"newPath" binding:"required"`
}

type bindContent struct {
	Path    string `json:"path" binding:"required"`
	Content string `json:"content"`
}

type bindChmod struct {
	Path string `json:"path" binding:"required"`
	Mode string `json:"mode" binding:"required"`
}
