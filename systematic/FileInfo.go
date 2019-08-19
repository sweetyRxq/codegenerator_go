package systematic

type FileInfo struct{
	FileId string `json:"fileId"`
	FileHash string `json:"fileHash"`
	FileName  string `json:"fileName"`
	FileSize string `json:"fileSize"`
	FileType string `json:"fileType"`
	UploadTime string `json:"uploadTime"`
	DataType string `json:"dataType"`
}