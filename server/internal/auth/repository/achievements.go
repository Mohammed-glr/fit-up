package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/tdmdh/fit-up-server/internal/auth/types"
)

func (s *Store) GetUserAchievements(ctx context.Context, userID string) ([]types.UserAchievement, error) {
	query := `
		SELECT 
			a.achievement_id,
			a.name,
			a.description,
			a.badge_icon,
			a.badge_color,
			a.category,
			a.requirement_value,
			a.points,
			COALESCE(ua.progress, 0) as progress,
			ua.earned_at,
			CASE WHEN ua.earned_at IS NOT NULL THEN true ELSE false END as is_completed
		FROM achievements a
		LEFT JOIN user_achievements ua ON a.achievement_id = ua.achievement_id AND ua.user_id = $1
		ORDER BY a.category, a.requirement_value
	`

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []types.UserAchievement
	for rows.Next() {
		var achievement types.UserAchievement
		var earnedAt *string

		err := rows.Scan(
			&achievement.AchievementID,
			&achievement.Name,
			&achievement.Description,
			&achievement.BadgeIcon,
			&achievement.BadgeColor,
			&achievement.Category,
			&achievement.RequirementValue,
			&achievement.Points,
			&achievement.Progress,
			&earnedAt,
			&achievement.IsCompleted,
		)
		if err != nil {
			log.Printf("Error scanning achievement: %v", err)
			continue
		}

		if achievement.RequirementValue > 0 {
			achievement.CompletionRate = (float64(achievement.Progress) / float64(achievement.RequirementValue)) * 100
			if achievement.CompletionRate > 100 {
				achievement.CompletionRate = 100
			}
		}

		achievements = append(achievements, achievement)
	}

	return achievements, nil
}

func (s *Store) GetAchievementStats(ctx context.Context, userID string) (*types.AchievementStats, error) {
	var stats types.AchievementStats

	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(ua.earned_at) as earned,
			COALESCE(SUM(CASE WHEN ua.earned_at IS NOT NULL THEN a.points ELSE 0 END), 0) as points
		FROM achievements a
		LEFT JOIN user_achievements ua ON a.achievement_id = ua.achievement_id AND ua.user_id = $1
	`

	err := s.db.QueryRow(ctx, query, userID).Scan(
		&stats.TotalAchievements,
		&stats.EarnedAchievements,
		&stats.TotalPoints,
	)
	if err != nil {
		return nil, err
	}

	if stats.TotalAchievements > 0 {
		stats.CompletionRate = (float64(stats.EarnedAchievements) / float64(stats.TotalAchievements)) * 100
	}

	return &stats, nil
}

func (s *Store) UpdateAchievementProgress(ctx context.Context, userID string, achievementID int, progress int) error {
	query := `
		INSERT INTO user_achievements (user_id, achievement_id, progress, earned_at)
		VALUES ($1, $2, $3, CASE WHEN $3 >= (SELECT requirement_value FROM achievements WHERE achievement_id = $2) THEN CURRENT_TIMESTAMP ELSE NULL END)
		ON CONFLICT (user_id, achievement_id)
		DO UPDATE SET 
			progress = GREATEST(user_achievements.progress, $3),
			earned_at = CASE 
				WHEN user_achievements.earned_at IS NULL AND $3 >= (SELECT requirement_value FROM achievements WHERE achievement_id = $2)
				THEN CURRENT_TIMESTAMP 
				ELSE user_achievements.earned_at 
			END
	`

	_, err := s.db.Exec(ctx, query, userID, achievementID, progress)
	return err
}

func (s *Store) CheckAndAwardAchievements(ctx context.Context, userID string) ([]types.UserAchievement, error) {
	stats, err := s.GetUserStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	var totalVolume float64
	var prCount int

	err = s.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(reps_completed * weight_used), 0) as total_volume
		FROM progress_logs
		WHERE user_id = $1
	`, userID).Scan(&totalVolume)
	if err != nil {
		log.Printf("Error getting total volume: %v", err)
	}

	err = s.db.QueryRow(ctx, `
		SELECT COUNT(DISTINCT exercise_id) as pr_count
		FROM progress_logs
		WHERE user_id = $1 AND weight_used > 0
	`, userID).Scan(&prCount)
	if err != nil {
		log.Printf("Error getting PR count: %v", err)
	}

	requirementValues := map[string]int{
		"total_workouts":   stats.TotalWorkouts,
		"streak_days":      stats.CurrentStreak,
		"total_volume_lbs": int(totalVolume),
		"pr_count":         prCount,
	}

	achievements, err := s.GetAllAchievements(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get achievements: %w", err)
	}

	var newlyEarned []types.UserAchievement

	for _, achievement := range achievements {
		progress, exists := requirementValues[achievement.RequirementType]
		if !exists {
			continue
		}

		var existingProgress int
		var earnedAt *string
		err := s.db.QueryRow(ctx, `
			SELECT progress, earned_at FROM user_achievements
			WHERE user_id = $1 AND achievement_id = $2
		`, userID, achievement.AchievementID).Scan(&existingProgress, &earnedAt)

		wasNotEarned := err != nil || earnedAt == nil
		isNowEarned := progress >= achievement.RequirementValue

		err = s.UpdateAchievementProgress(ctx, userID, achievement.AchievementID, progress)
		if err != nil {
			log.Printf("Error updating achievement progress: %v", err)
			continue
		}

		if wasNotEarned && isNowEarned {
			newlyEarned = append(newlyEarned, types.UserAchievement{
				AchievementID:    achievement.AchievementID,
				Name:             achievement.Name,
				Description:      achievement.Description,
				BadgeIcon:        achievement.BadgeIcon,
				BadgeColor:       achievement.BadgeColor,
				Category:         achievement.Category,
				RequirementValue: achievement.RequirementValue,
				Points:           achievement.Points,
				Progress:         progress,
				IsCompleted:      true,
				CompletionRate:   100,
			})
		}
	}

	return newlyEarned, nil
}

func (s *Store) GetAllAchievements(ctx context.Context) ([]types.Achievement, error) {
	query := `
		SELECT achievement_id, name, description, badge_icon, badge_color, 
		       category, requirement_type, requirement_value, points, created_at
		FROM achievements
		ORDER BY category, requirement_value
	`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []types.Achievement
	for rows.Next() {
		var achievement types.Achievement
		err := rows.Scan(
			&achievement.AchievementID,
			&achievement.Name,
			&achievement.Description,
			&achievement.BadgeIcon,
			&achievement.BadgeColor,
			&achievement.Category,
			&achievement.RequirementType,
			&achievement.RequirementValue,
			&achievement.Points,
			&achievement.CreatedAt,
		)
		if err != nil {
			log.Printf("Error scanning achievement: %v", err)
			continue
		}
		achievements = append(achievements, achievement)
	}

	return achievements, nil
}
