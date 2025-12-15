package handlers

import (
	"log"
	"net/http"
	"sustainwear/internal/config"
	"sustainwear/pkg/fileupload"

	jsoniter "github.com/json-iterator/go"
)

type UploadHandler struct {
	config *config.Config
}

func NewUploadHandler(cfg *config.Config) *UploadHandler {
	return &UploadHandler{
		config: cfg,
	}
}

// UPLOAD IMAGES ENDPOINT
func (h *UploadHandler) UploadImages(w http.ResponseWriter, r *http.Request) {
	// PARSES MULTIPART FORM
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		log.Printf("UPLOAD: Failed to parse multipart form: %v", err)
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		log.Printf("UPLOAD: No files provided in request")
		http.Error(w, "No files provided. Use 'images' field for file uploads", http.StatusBadRequest)
		return
	}

	var uploadedPaths []string
	var errors []string

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			errMsg := "Failed to open file: " + fileHeader.Filename
			errors = append(errors, errMsg)
			continue
		}

		savedFilename, err := fileupload.ValidateAndSaveFile(file, fileHeader, h.config)
		file.Close()

		if err != nil {
			errMsg := fileHeader.Filename + ": " + err.Error()
			errors = append(errors, errMsg)
			continue
		}

		filePath := "uploads/" + savedFilename
		uploadedPaths = append(uploadedPaths, filePath)
	}

	if len(errors) > 0 && len(uploadedPaths) == 0 {
		log.Printf("UPLOAD: All %d file(s) failed to upload", len(files))
	} else if len(errors) > 0 {
		log.Printf("UPLOAD: Partial success - %d succeeded, %d failed", len(uploadedPaths), len(errors))
	} else {
		log.Printf("UPLOAD: Successfully uploaded %d file(s)", len(uploadedPaths))
	}

	response := map[string]interface{}{
		"success": len(uploadedPaths),
		"failed":  len(errors),
		"paths":   uploadedPaths,
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	statusCode := http.StatusOK
	if len(uploadedPaths) == 0 {
		statusCode = http.StatusBadRequest
	} else if len(errors) > 0 {
		statusCode = http.StatusMultiStatus
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsoniter.NewEncoder(w).Encode(response)
}
