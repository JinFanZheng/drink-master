package services

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/models"
)

// MemberService handles member operations
type MemberService struct {
	db *gorm.DB
}

// NewMemberService creates a new member service
func NewMemberService(db *gorm.DB) *MemberService {
	return &MemberService{
		db: db,
	}
}

// FindByOpenID finds a member by WeChat OpenID
func (s *MemberService) FindByOpenID(openID string) (*models.Member, error) {
	var member models.Member
	err := s.db.Where("we_chat_open_id = ?", openID).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// FindOrCreateByOpenID finds a member by OpenID or creates a new one
func (s *MemberService) FindOrCreateByOpenID(openID, nickname, avatarUrl string) (*models.Member, error) {
	var member models.Member

	// Try to find existing member
	err := s.db.Where("we_chat_open_id = ?", openID).First(&member).Error
	if err == nil {
		// Member exists, update avatar and nickname if provided
		if nickname != "" {
			member.Nickname = nickname
		}
		if avatarUrl != "" {
			member.Avatar = avatarUrl
		}
		if err := s.db.Save(&member).Error; err != nil {
			return nil, err
		}
		return &member, nil
	}

	// Member not found, create new one
	if err == gorm.ErrRecordNotFound {
		member = models.Member{
			ID:           uuid.New().String(),
			Nickname:     nickname,
			Avatar:       avatarUrl,
			WeChatOpenId: openID,
			Role:         "Member", // Default role
		}

		if err := s.db.Create(&member).Error; err != nil {
			return nil, err
		}
		return &member, nil
	}

	return nil, err
}

// FindByID finds a member by ID
func (s *MemberService) FindByID(id string) (*models.Member, error) {
	var member models.Member
	err := s.db.Where("id = ?", id).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}
