package db

import (
	"regexp"
	"strings"

	"gorm.io/gorm"
)

var validRoomRgx = regexp.MustCompile(`^[a-z0-9]([a-z0-9-_]{0,30}[a-z0-9])?$`)

func NormalizeRoomName(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

func IsValidRoomName(s string) bool {
	return validRoomRgx.MatchString(s)
}

func CreateRoom(g *gorm.DB, name, creator string) error {
	r := Room{Name: name, Creator: creator}
	return g.Create(&r).Error
}

func ListRoomsForUser(g *gorm.DB, username string) ([]string, error) {
	var created []Room
	if err := g.Where("creator = ?", username).Find(&created).Error; err != nil {
		return nil, err
	}

	var posted []struct{ Room string }
	if err := g.Model(&Message{}).
		Select("DISTINCT room").
		Where("username = ?", username).
		Find(&posted).Error; err != nil {
		return nil, err
	}

	m := map[string]struct{}{}
	for _, r := range created {
		m[r.Name] = struct{}{}
	}
	for _, r := range posted {
		m[r.Room] = struct{}{}
	}

	m["general"] = struct{}{}

	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out, nil
}
