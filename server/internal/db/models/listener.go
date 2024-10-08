package models

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"time"
)

// Listener
type Listener struct {
	ID        uuid.UUID `gorm:"primaryKey;->;<-:create;type:uuid;"`
	CreatedAt time.Time `gorm:"->;<-:create;"`
	Name      string    `gorm:"uniqueIndex"`
}

// BeforeCreate - GORM hook
func (l *Listener) BeforeCreate(tx *gorm.DB) (err error) {
	l.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}
	l.CreatedAt = time.Now()
	return nil
}
