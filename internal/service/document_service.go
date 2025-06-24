package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
	"wz-wenzhan-backend/internal/model"
	"wz-wenzhan-backend/internal/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DocumentService interface {
	Create(userID uint, req *model.CreateDocumentRequest) (*model.Document, error)
	GetByID(id, userID uint) (*model.DocumentDetailResponse, error)
	Update(id, userID uint, req *model.UpdateDocumentRequest) error
	Delete(id, userID uint) error
	List(userID uint, req *model.DocumentListRequest) ([]model.DocumentResponse, int64, error)
	Share(id, userID uint, req *model.ShareDocumentRequest) (string, error)
	Copy(id, userID uint) (*model.Document, error)
	GetByShareToken(token string) (*model.DocumentDetailResponse, error)
}

type documentService struct {
	documentRepo repository.DocumentRepository
	logger       *zap.Logger
}

func NewDocumentService(documentRepo repository.DocumentRepository, logger *zap.Logger) DocumentService {
	return &documentService{
		documentRepo: documentRepo,
		logger:       logger,
	}
}

func (s *documentService) Create(userID uint, req *model.CreateDocumentRequest) (*model.Document, error) {
	document := &model.Document{
		Title:    req.Title,
		Content:  req.Content,
		Type:     req.Type,
		Status:   model.DocumentStatusDraft,
		UserID:   userID,
		FolderID: req.FolderID,
		Tags:     req.Tags,
		Size:     int64(len(req.Content)),
	}

	err := s.documentRepo.Create(document)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Document created", 
		zap.Uint("user_id", userID), 
		zap.Uint("document_id", document.ID),
		zap.String("title", document.Title))

	return document, nil
}

func (s *documentService) GetByID(id, userID uint) (*model.DocumentDetailResponse, error) {
	document, err := s.documentRepo.GetByIDAndUserID(id, userID)
	if err != nil {
		return nil, err
	}

	// 更新查看次数
	go func() {
		s.documentRepo.UpdateViewCount(id)
	}()

	response := &model.DocumentDetailResponse{
		DocumentResponse: model.DocumentResponse{
			ID:        document.ID,
			Title:     document.Title,
			Type:      document.Type,
			Status:    document.Status,
			FolderID:  document.FolderID,
			Tags:      document.Tags,
			Size:      document.Size,
			ViewCount: document.ViewCount,
			IsShared:  document.IsShared,
			CreatedAt: document.CreatedAt,
			UpdatedAt: document.UpdatedAt,
		},
		Content: document.Content,
	}

	return response, nil
}

func (s *documentService) Update(id, userID uint, req *model.UpdateDocumentRequest) error {
	document, err := s.documentRepo.GetByIDAndUserID(id, userID)
	if err != nil {
		return err
	}

	if req.Title != "" {
		document.Title = req.Title
	}
	if req.Content != "" {
		document.Content = req.Content
		document.Size = int64(len(req.Content))
	}
	if req.Status != nil {
		document.Status = *req.Status
	}
	if req.FolderID != nil {
		document.FolderID = req.FolderID
	}
	if req.Tags != "" {
		document.Tags = req.Tags
	}

	err = s.documentRepo.Update(document)
	if err != nil {
		return err
	}

	s.logger.Info("Document updated", 
		zap.Uint("user_id", userID), 
		zap.Uint("document_id", id))

	return nil
}

func (s *documentService) Delete(id, userID uint) error {
	document, err := s.documentRepo.GetByIDAndUserID(id, userID)
	if err != nil {
		return err
	}

	// 直接删除文档（简化版本，不移动到回收站）
	err = s.documentRepo.Delete(id, userID)
	if err != nil {
		return err
	}

	s.logger.Info("Document deleted", 
		zap.Uint("user_id", userID), 
		zap.Uint("document_id", id))

	return nil
}

func (s *documentService) List(userID uint, req *model.DocumentListRequest) ([]model.DocumentResponse, int64, error) {
	documents, total, err := s.documentRepo.List(userID, req)
	if err != nil {
		return nil, 0, err
	}

	var responses []model.DocumentResponse
	for _, doc := range documents {
		responses = append(responses, model.DocumentResponse{
			ID:        doc.ID,
			Title:     doc.Title,
			Type:      doc.Type,
			Status:    doc.Status,
			FolderID:  doc.FolderID,
			Tags:      doc.Tags,
			Size:      doc.Size,
			ViewCount: doc.ViewCount,
			IsShared:  doc.IsShared,
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
		})
	}

	return responses, total, nil
}

func (s *documentService) Share(id, userID uint, req *model.ShareDocumentRequest) (string, error) {
	document, err := s.documentRepo.GetByIDAndUserID(id, userID)
	if err != nil {
		return "", err
	}

	// 生成分享token
	token, err := generateRandomString(32)
	if err != nil {
		return "", err
	}

	// 设置分享信息
	document.IsShared = true
	document.ShareToken = token
	expiry := time.Now().Add(time.Duration(req.ExpiryHours) * time.Hour)
	document.ShareExpiry = &expiry

	err = s.documentRepo.Update(document)
	if err != nil {
		return "", err
	}

	s.logger.Info("Document shared", 
		zap.Uint("user_id", userID), 
		zap.Uint("document_id", id),
		zap.String("token", token))

	return token, nil
}

func (s *documentService) Copy(id, userID uint) (*model.Document, error) {
	original, err := s.documentRepo.GetByIDAndUserID(id, userID)
	if err != nil {
		return nil, err
	}

	// 创建副本
	copy := &model.Document{
		Title:    original.Title + " - 副本",
		Content:  original.Content,
		Type:     original.Type,
		Status:   model.DocumentStatusDraft,
		UserID:   userID,
		FolderID: original.FolderID,
		Tags:     original.Tags,
		Size:     original.Size,
	}

	err = s.documentRepo.Create(copy)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Document copied", 
		zap.Uint("user_id", userID), 
		zap.Uint("original_id", id),
		zap.Uint("copy_id", copy.ID))

	return copy, nil
}

func (s *documentService) GetByShareToken(token string) (*model.DocumentDetailResponse, error) {
	document, err := s.documentRepo.GetByShareToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分享链接无效或已过期")
		}
		return nil, err
	}

	response := &model.DocumentDetailResponse{
		DocumentResponse: model.DocumentResponse{
			ID:        document.ID,
			Title:     document.Title,
			Type:      document.Type,
			Status:    document.Status,
			FolderID:  document.FolderID,
			Tags:      document.Tags,
			Size:      document.Size,
			ViewCount: document.ViewCount,
			IsShared:  document.IsShared,
			CreatedAt: document.CreatedAt,
			UpdatedAt: document.UpdatedAt,
		},
		Content: document.Content,
	}

	// 更新查看次数
	go func() {
		s.documentRepo.UpdateViewCount(document.ID)
	}()

	return response, nil
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
