package types

import "time"

// Achievement represents a badge/achievement that can be earned
type Achievement struct {
	AchievementID    int       `json:"achievement_id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	BadgeIcon        string    `json:"badge_icon"`
	BadgeColor       string    `json:"badge_color"`
	Category         string    `json:"category"`
	RequirementType  string    `json:"requirement_type"`
	RequirementValue int       `json:"requirement_value"`
	Points           int       `json:"points"`
	CreatedAt        time.Time `json:"created_at"`
}

// UserAchievement represents a user's progress or completion of an achievement
type UserAchievement struct {
	UserAchievementID int        `json:"user_achievement_id"`
	UserID            string     `json:"user_id"`
	AchievementID     int        `json:"achievement_id"`
	Progress          int        `json:"progress"`
	EarnedAt          *time.Time `json:"earned_at,omitempty"`
	// Embedded achievement details
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	BadgeIcon        string  `json:"badge_icon"`
	BadgeColor       string  `json:"badge_color"`
	Category         string  `json:"category"`
	RequirementValue int     `json:"requirement_value"`
	Points           int     `json:"points"`
	IsCompleted      bool    `json:"is_completed"`
	CompletionRate   float64 `json:"completion_rate"`
}

// AchievementStats represents overall achievement statistics for a user
type AchievementStats struct {
	TotalAchievements  int     `json:"total_achievements"`
	EarnedAchievements int     `json:"earned_achievements"`
	TotalPoints        int     `json:"total_points"`
	CompletionRate     float64 `json:"completion_rate"`
}
