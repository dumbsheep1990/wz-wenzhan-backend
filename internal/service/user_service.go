package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
	"wz-wenzhan-backend/internal/model"
	"wz-wenzhan-backend/internal/repository"
	"wz-wenzhan-backend/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Register(req *model.RegisterRequest) (*model.User, error)
	Login(req *model.LoginRequest) (string, *model.UserProfile, error)
	GetProfile(userID uint) (*model.UserProfile, error)
	UpdateProfile(userID uint, req *model.UpdateProfileRequest) error
	ChangePassword(userID uint, oldPassword, newPassword string) error
}

type userService struct {
	userRepo repository.UserRepository
	logger   *zap.Logger
}

func NewUserService(userRepo repository.UserRepository, logger *zap.Logger) UserService {
	return &userService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *userService) Register(req *model.RegisterRequest) (*model.User, error) {
	// 检查用户名是否已存在
	existingUser, err := s.userRepo.GetByUsername(req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	existingUser, err = s.userRepo.GetByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Status:   1,
	}

	if user.Nickname == "" {
		user.Nickname = user.Username
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	s.logger.Info("User registered", zap.String("username", user.Username))
	return user, nil
}

func (s *userService) Login(req *model.LoginRequest) (string, *model.UserProfile, error) {
	// 获取用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("用户名或密码错误")
		}
		return "", nil, err
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return "", nil, errors.New("用户已被禁用")
	}

	// 生成JWT token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", nil, err
	}

	// 更新最后登录时间
	err = s.userRepo.UpdateLastLogin(user.ID)
	if err != nil {
		s.logger.Error("Failed to update last login", zap.Error(err))
	}

	profile := &model.UserProfile{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Nickname: user.Nickname,
		Status:   user.Status,
	}

	s.logger.Info("User logged in", zap.String("username", user.Username))
	return token, profile, nil
}

func (s *userService) GetProfile(userID uint) (*model.UserProfile, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	profile := &model.UserProfile{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Nickname: user.Nickname,
		Status:   user.Status,
	}

	return profile, nil
}

func (s *userService) UpdateProfile(userID uint, req *model.UpdateProfileRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return err
	}

	s.logger.Info("User profile updated", zap.Uint("user_id", userID))
	return nil
}

func (s *userService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	err = s.userRepo.Update(user)
	if err != nil {
		return err
	}

	s.logger.Info("User password changed", zap.Uint("user_id", userID))
	return nil
}
