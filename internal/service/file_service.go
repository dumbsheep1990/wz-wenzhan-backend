package service

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
	"wz-wenzhan-backend/internal/model"

	"go.uber.org/zap"
)

type FileService interface {
	UploadFile(userID uint, file *multipart.FileHeader) (*model.FileUploadResponse, error)
	DeleteFile(userID uint, fileID string) error
	GetFileURL(fileID string) string
}

type fileService struct {
	uploadPath string
	baseURL    string
	logger     *zap.Logger
}

func NewFileService(uploadPath, baseURL string, logger *zap.Logger) FileService {
	// 确保上传目录存在
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		logger.Error("Failed to create upload directory", zap.Error(err))
	}

	return &fileService{
		uploadPath: uploadPath,
		baseURL:    baseURL,
		logger:     logger,
	}
}

func (s *fileService) UploadFile(userID uint, fileHeader *multipart.FileHeader) (*model.FileUploadResponse, error) {
	// 检查文件大小（限制为10MB）
	maxSize := int64(10 << 20) // 10MB
	if fileHeader.Size > maxSize {
		return nil, fmt.Errorf("文件大小超过限制，最大支持10MB")
	}

	// 检查文件类型
	allowedTypes := map[string]bool{
		"image/jpeg":                               true,
		"image/png":                                true,
		"image/gif":                                true,
		"application/pdf":                          true,
		"application/msword":                       true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   true,
		"application/vnd.ms-excel":                 true,
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         true,
		"application/vnd.ms-powerpoint":            true,
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": true,
		"text/plain":                               true,
		"text/markdown":                            true,
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return nil, fmt.Errorf("不支持的文件类型: %s", contentType)
	}

	// 打开上传的文件
	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// 生成文件ID和路径
	fileID := s.generateFileID(userID, fileHeader.Filename)
	userDir := filepath.Join(s.uploadPath, fmt.Sprintf("user_%d", userID))
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return nil, err
	}

	// 生成文件扩展名
	ext := filepath.Ext(fileHeader.Filename)
	fileName := fileID + ext
	filePath := filepath.Join(userDir, fileName)

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// 复制文件内容
	_, err = io.Copy(dst, src)
	if err != nil {
		return nil, err
	}

	// 生成文件URL
	fileURL := fmt.Sprintf("%s/uploads/user_%d/%s", s.baseURL, userID, fileName)

	response := &model.FileUploadResponse{
		ID:       fileID,
		Filename: fileHeader.Filename,
		URL:      fileURL,
		Size:     fileHeader.Size,
		MimeType: contentType,
	}

	s.logger.Info("File uploaded",
		zap.Uint("user_id", userID),
		zap.String("file_id", fileID),
		zap.String("filename", fileHeader.Filename),
		zap.Int64("size", fileHeader.Size))

	return response, nil
}

func (s *fileService) DeleteFile(userID uint, fileID string) error {
	// 查找文件
	userDir := filepath.Join(s.uploadPath, fmt.Sprintf("user_%d", userID))
	
	// 遍历用户目录找到对应的文件
	files, err := os.ReadDir(userDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), fileID) {
			filePath := filepath.Join(userDir, file.Name())
			err := os.Remove(filePath)
			if err != nil {
				return err
			}

			s.logger.Info("File deleted",
				zap.Uint("user_id", userID),
				zap.String("file_id", fileID))

			return nil
		}
	}

	return fmt.Errorf("文件不存在")
}

func (s *fileService) GetFileURL(fileID string) string {
	// 这里可以根据文件ID生成访问URL
	return fmt.Sprintf("%s/uploads/%s", s.baseURL, fileID)
}

func (s *fileService) generateFileID(userID uint, filename string) string {
	// 使用用户ID、文件名和时间戳生成唯一ID
	data := fmt.Sprintf("%d_%s_%d", userID, filename, time.Now().UnixNano())
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}
