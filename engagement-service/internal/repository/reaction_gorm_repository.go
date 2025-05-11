package repository

import (
	"engagement-service/internal/db"
	"engagement-service/internal/models"

	"gorm.io/gorm"
)

type GormReactionRepository struct {
	db *gorm.DB
}

func NewGormReactionRepository() *GormReactionRepository {
	return &GormReactionRepository{
		db: db.DB,
	}
}

func (r *GormReactionRepository) ToggleReaction(videoID, userID, reactionType string) error {
	var reaction models.Reaction

	err := r.db.Where("video_id = ? AND user_id = ?", videoID, userID).First(&reaction).Error

	if err == gorm.ErrRecordNotFound {
		// ما فيش تفاعل سابق → نضيف جديد
		return r.db.Create(&models.Reaction{
			VideoID: videoID,
			UserID:  userID,
			Type:    reactionType,
		}).Error
	} else if err != nil {
		return err
	}

	// التفاعل موجود فعلاً:
	if reaction.Type == reactionType {
		// نفس التفاعل → نعمل "إلغاء"
		return r.db.Delete(&reaction).Error
	}

	// تغيير نوع التفاعل
	return r.db.Model(&reaction).Update("type", reactionType).Error
}

func (r *GormReactionRepository) CountReactions(videoID string) (int64, int64, error) {
	var likes int64
	var dislikes int64

	err := r.db.Model(&models.Reaction{}).Where("video_id = ? AND type = ?", videoID, "like").Count(&likes).Error
	if err != nil {
		return 0, 0, err
	}

	err = r.db.Model(&models.Reaction{}).Where("video_id = ? AND type = ?", videoID, "dislike").Count(&dislikes).Error
	if err != nil {
		return 0, 0, err
	}

	return likes, dislikes, nil
}
