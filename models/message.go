package models

import "time"

type Message struct {
    ID         uint      `gorm:"primaryKey"`
    SenderID   uint      `gorm:"not null"`
    ReceiverID uint      `gorm:"not null"`
    Content    string    `gorm:"type:text;not null"`
    SentAt     time.Time `gorm:"autoCreateTime"`

    Sender   User `gorm:"foreignKey:SenderID"`
    Receiver User `gorm:"foreignKey:ReceiverID"`
}
