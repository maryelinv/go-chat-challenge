package db

import "gorm.io/gorm"

func SaveMessage(g *gorm.DB, m *Message) error { return g.Create(m).Error }

func LastNMessages(g *gorm.DB, room string, n int) ([]Message, error) {
	var mm []Message
	err := g.Where("room = ?", room).
		Order("created_at DESC").Limit(n).Find(&mm).Error

	return mm, err
}

func LastNMessagesByUser(g *gorm.DB, room, username string, n int) ([]Message, error) {
	var mm []Message
	err := g.Where("room = ? AND username = ?", room, username).
		Order("created_at DESC").Limit(n).Find(&mm).Error
	return mm, err
}

func RoomsForUser(g *gorm.DB, username string) ([]string, error) {
	var rows []struct{ Room string }
	err := g.Model(&Message{}).
		Select("DISTINCT room").
		Where("username = ?", username).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	rooms := make([]string, 0, len(rows))
	for _, r := range rows {
		rooms = append(rooms, r.Room)
	}
	return rooms, nil
}
