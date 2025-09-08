package repo

import (
	"context"

	"github.com/hassiimykyta/life-rpg/apps/auth-svc/internal/models"
	"gorm.io/gorm"
)

type IdentityRepo struct {
	db *gorm.DB
}

func NewIdentityRepo(db *gorm.DB) *IdentityRepo { return &IdentityRepo{db: db} }

func (r *IdentityRepo) Create(ctx context.Context, id models.Identity) error {
	return r.db.WithContext(ctx).Create(&id).Error
}

func (r *IdentityRepo) FindByEmail(ctx context.Context, email string) (models.Identity, error) {
	var m models.Identity
	err := r.db.WithContext(ctx).First(&m, "email = ?", email).Error
	return m, err
}

func (r *IdentityRepo) FindByUsername(ctx context.Context, username string) (models.Identity, error) {
	var m models.Identity
	err := r.db.WithContext(ctx).First(&m, "username = ?", username).Error
	return m, err
}

func (r *IdentityRepo) FindByUserID(ctx context.Context, userID string) (models.Identity, error) {
	var m models.Identity
	err := r.db.WithContext(ctx).First(&m, "user_id = ?", userID).Error
	return m, err
}
