package domains

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	UserID     *uuid.UUID
	UserName   string
	Content    string
	ConnStatus string
	Time       time.Time
}
