package dto

type PresignedPutRequest struct {
	Key    string `json:"key" binding:"required"`
	Expire int    `json:"expire"`
}

type PresignedPutResponse struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

type HeadObjectResponse struct {
	Key  string `json:"key"`
	Size int64  `json:"size"`
}

type UploadResponse struct {
	Key      string `json:"key"`
	Path     string `json:"path"`
	FullPath string `json:"full_path"`
	Size     int64  `json:"size"`
}
