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

type UploadDriverDocumentRequest struct {
	DriverProfileID string `form:"driver_profile_id" binding:"required"`
	DocumentTypeID  string `form:"document_type_id" binding:"required"`
}

type UploadDriverDocumentItemResponse struct {
	DocumentID       string `json:"document_id"`
	Key              string `json:"key"`
	Size             int64  `json:"size"`
	OriginalFilename string `json:"original_filename"`
	DriverProfileID  string `json:"driver_profile_id"`
	DocumentTypeID   string `json:"document_type_id"`
	Bucket           string `json:"bucket"`
}

