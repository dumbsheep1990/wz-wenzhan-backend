package service

import (
	"wz-wenzhan-backend/internal/model"
	"wz-wenzhan-backend/internal/repository"
	"strings"

	"go.uber.org/zap"
)

type SearchService interface {
	SearchDocuments(userID uint, req *model.SearchRequest) (*model.PaginationResponse, error)
	SearchAll(userID uint, req *model.SearchRequest) (map[string]interface{}, error)
}

type searchService struct {
	documentRepo repository.DocumentRepository
	folderRepo   repository.FolderRepository
	logger       *zap.Logger
}

func NewSearchService(
	documentRepo repository.DocumentRepository,
	folderRepo repository.FolderRepository,
	logger *zap.Logger) SearchService {
	return &searchService{
		documentRepo: documentRepo,
		folderRepo:   folderRepo,
		logger:       logger,
	}
}

func (s *searchService) SearchDocuments(userID uint, req *model.SearchRequest) (*model.PaginationResponse, error) {
	req.SetDefaults()

	// 使用基本的文档查询实现搜索
	documents, total, err := s.searchDocumentsBasic(userID, req)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
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

	return model.NewPaginationResponse(responses, total, req.Page, req.PageSize), nil
}

func (s *searchService) SearchAll(userID uint, req *model.SearchRequest) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// 搜索文档
	documents, _, err := s.searchDocumentsBasic(userID, req)
	if err != nil {
		return nil, err
	}

	// 搜索文件夹
	folders, err := s.searchFolders(userID, req.Keyword)
	if err != nil {
		return nil, err
	}

	result["documents"] = documents
	result["folders"] = folders

	return result, nil
}

func (s *searchService) searchDocumentsBasic(userID uint, req *model.SearchRequest) ([]model.Document, int64, error) {
	// 构建文档列表请求
	listReq := &model.DocumentListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
	}

	return s.documentRepo.GetUserDocuments(userID, listReq)
}

func (s *searchService) searchFolders(userID uint, keyword string) ([]model.Folder, error) {
	if keyword == "" {
		return []model.Folder{}, nil
	}

	// 获取用户所有文件夹，然后在应用层过滤
	folders, err := s.folderRepo.GetUserFolders(userID)
	if err != nil {
		return nil, err
	}

	var matchedFolders []model.Folder
	for _, folder := range folders {
		// 简单的名称匹配（忽略大小写）
		if strings.Contains(strings.ToLower(folder.Name), strings.ToLower(keyword)) {
			matchedFolders = append(matchedFolders, folder)
		}
	}

	return matchedFolders, nil
}
