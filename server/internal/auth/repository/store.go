package repository

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tdmdh/fit-up-server/internal/auth/types"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	query := `
		SELECT id, username, name, bio, email, email_verified, image, password, role, is_two_factor_enabled, created_at, updated_at
		FROM users 
		WHERE email = $1
	`

	var user types.User
	err := s.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Name,
		&user.Bio,
		&user.Email,
		&user.EmailVerified,
		&user.Image,
		&user.PasswordHash,
		&user.Role,
		&user.IsTwoFactorEnabled,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, types.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *Store) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	query := `
		SELECT id, username, name, bio, email, email_verified, image, password, role, is_two_factor_enabled, created_at, updated_at
		FROM users 
		WHERE id = $1
	`

	var user types.User
	err := s.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Name,
		&user.Bio,
		&user.Email,
		&user.EmailVerified,
		&user.Image,
		&user.PasswordHash,
		&user.Role,
		&user.IsTwoFactorEnabled,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, types.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *Store) GetUserByUsername(ctx context.Context, username string) (*types.User, error) {
	query := `
		SELECT id, username, name, bio, email, email_verified, image, password, role, is_two_factor_enabled, created_at, updated_at
		FROM users 
		WHERE username = $1
	`

	var user types.User
	err := s.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Name,
		&user.Bio,
		&user.Email,
		&user.EmailVerified,
		&user.Image,
		&user.PasswordHash,
		&user.Role,
		&user.IsTwoFactorEnabled,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, types.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *Store) CreateUser(ctx context.Context, user *types.User) error {
	query := `
		INSERT INTO users (id, username, name, bio, email, email_verified, image, password, role, is_two_factor_enabled)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := s.db.Exec(ctx, query,
		user.ID,
		user.Username,
		user.Name,
		user.Bio,
		user.Email,
		user.EmailVerified,
		user.Image,
		user.PasswordHash,
		user.Role,
		user.IsTwoFactorEnabled,
	)

	return err
}

func (s *Store) UpdateUser(ctx context.Context, id string, updates *types.UpdateUserRequest) error {
	query := `
		UPDATE users 
		SET name = COALESCE($2, name), 
		    bio = COALESCE($3, bio), 
		    image = COALESCE($4, image),
		    updated_at = NOW()
		WHERE id = $1
	`

	log.Printf("💾 Executing UPDATE for user %s - Name: %v, Bio: %v, Image present: %v",
		id,
		updates.Name,
		updates.Bio,
		updates.Image != nil,
	)

	result, err := s.db.Exec(ctx, query, id, updates.Name, updates.Bio, updates.Image)
	if err != nil {
		log.Printf("❌ Database error: %v", err)
		return err
	}

	rowsAffected := result.RowsAffected()
	log.Printf("✅ Rows affected: %d", rowsAffected)

	if rowsAffected == 0 {
		return types.ErrUserNotFound
	}

	return nil
}

func (s *Store) UpdateUserPassword(ctx context.Context, userID string, hashedPassword string) error {
	query := `
		UPDATE users 
		SET password = $2, updated_at = NOW()
		WHERE id = $1
	`

	_, err := s.db.Exec(ctx, query, userID, hashedPassword)
	return err
}

func (s *Store) UpdateUserRole(ctx context.Context, userID string, role types.UserRole) error {
	query := `
		UPDATE users 
		SET role = $2, updated_at = NOW()
		WHERE id = $1
	`

	_, err := s.db.Exec(ctx, query, userID, role)
	return err
}

func (s *Store) CreatePasswordResetToken(ctx context.Context, email string, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO password_reset_tokens (email, token, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) 
		DO UPDATE SET token = $2, expires_at = $3, created_at = NOW()
	`

	_, err := s.db.Exec(ctx, query, email, token, expiresAt)
	return err
}

func (s *Store) GetPasswordResetToken(ctx context.Context, token string) (*types.PasswordResetToken, error) {
	query := `
		SELECT email, token, expires_at, used
		FROM password_reset_tokens 
		WHERE token = $1
	`

	var resetToken types.PasswordResetToken
	err := s.db.QueryRow(ctx, query, token).Scan(
		&resetToken.Email,
		&resetToken.Token,
		&resetToken.Expires,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, types.ErrPasswordResetTokenNotFound
		}
		return nil, err
	}

	return &resetToken, nil
}

func (s *Store) GetUserByPasswordResetToken(ctx context.Context, token string) (*types.User, error) {
	query := `
		SELECT u.id, u.username, u.name, u.bio, u.email, u.email_verified, u.image, u.password, u.role, u.is_two_factor_enabled, u.created_at, u.updated_at
		FROM users u
		INNER JOIN password_reset_tokens prt ON u.email = prt.email
		WHERE prt.token = $1 AND prt.expires_at > NOW() AND prt.used = false
	`

	var user types.User
	err := s.db.QueryRow(ctx, query, token).Scan(
		&user.ID,
		&user.Username,
		&user.Name,
		&user.Bio,
		&user.Email,
		&user.EmailVerified,
		&user.Image,
		&user.PasswordHash,
		&user.Role,
		&user.IsTwoFactorEnabled,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, types.ErrPasswordResetTokenNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *Store) DeletePasswordResetToken(ctx context.Context, token string) error {
	query := `DELETE FROM password_reset_tokens WHERE token = $1`
	_, err := s.db.Exec(ctx, query, token)
	return err
}

func (s *Store) CreateVerificationToken(ctx context.Context, email, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO verification_tokens (email, token, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (email)
		DO UPDATE SET token = EXCLUDED.token, expires_at = EXCLUDED.expires_at, updated_at = NOW(), consumed_at = NULL
	`

	_, err := s.db.Exec(ctx, query, email, token, expiresAt)
	return err
}

func (s *Store) GetVerificationToken(ctx context.Context, token string) (*types.VerificationToken, error) {
	query := `
		SELECT id, email, token, expires_at
		FROM verification_tokens
		WHERE token = $1 AND consumed_at IS NULL
	`

	var verificationToken types.VerificationToken
	err := s.db.QueryRow(ctx, query, token).Scan(
		&verificationToken.ID,
		&verificationToken.Email,
		&verificationToken.Token,
		&verificationToken.Expires,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, types.ErrVerificationTokenNotFound
		}
		return nil, err
	}

	return &verificationToken, nil
}

func (s *Store) DeleteVerificationToken(ctx context.Context, token string) error {
	query := `DELETE FROM verification_tokens WHERE token = $1`
	_, err := s.db.Exec(ctx, query, token)
	return err
}

func (s *Store) MarkEmailVerified(ctx context.Context, userID string, verifiedAt time.Time) error {
	query := `
		UPDATE users
		SET email_verified = $2, updated_at = NOW()
		WHERE id = $1
	`

	result, err := s.db.Exec(ctx, query, userID, verifiedAt)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return types.ErrUserNotFound
	}

	return nil
}

func (s *Store) MarkPasswordResetTokenAsUsed(ctx context.Context, token string) error {
	query := `
		UPDATE password_reset_tokens 
		SET used = true 
		WHERE token = $1
	`

	_, err := s.db.Exec(ctx, query, token)
	return err
}

func (s *Store) CreateRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time, accessTokenJTI string) error {
	query := `
		INSERT INTO jwt_refresh_tokens (user_id, token_hash, access_token_jti, expires_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := s.db.Exec(ctx, query, userID, token, accessTokenJTI, expiresAt)
	return err
}

func (s *Store) GetRefreshToken(ctx context.Context, token string) (*types.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, access_token_jti, expires_at, created_at, last_used_at, is_revoked, revoked_at
		FROM jwt_refresh_tokens 
		WHERE token_hash = $1
	`

	var refreshToken types.RefreshToken
	err := s.db.QueryRow(ctx, query, token).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.TokenHash,
		&refreshToken.AccessTokenJTI,
		&refreshToken.ExpiresAt,
		&refreshToken.CreatedAt,
		&refreshToken.LastUsedAt,
		&refreshToken.IsRevoked,
		&refreshToken.RevokedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, types.ErrRefreshTokenNotFound
		}
		return nil, err
	}

	return &refreshToken, nil
}

func (s *Store) DeleteRefreshToken(ctx context.Context, token string) error {
	query := `DELETE FROM jwt_refresh_tokens WHERE token_hash = $1`
	_, err := s.db.Exec(ctx, query, token)
	return err
}

func (s *Store) CleanupExpiredRefreshTokens(ctx context.Context) error {
	query := `DELETE FROM jwt_refresh_tokens WHERE expires_at < NOW()`
	_, err := s.db.Exec(ctx, query)
	return err
}

func (s *Store) RevokeRefreshToken(ctx context.Context, token string) error {
	query := `
		UPDATE jwt_refresh_tokens 
		SET is_revoked = true, revoked_at = NOW() 
		WHERE token_hash = $1
	`

	_, err := s.db.Exec(ctx, query, token)
	return err
}

func (s *Store) RevokeAllUserRefreshTokens(ctx context.Context, userID string) error {
	query := `
		UPDATE jwt_refresh_tokens 
		SET is_revoked = true, revoked_at = NOW() 
		WHERE user_id = $1 AND is_revoked = false
	`

	_, err := s.db.Exec(ctx, query, userID)
	return err
}

func (s *Store) UpdateRefreshTokenLastUsed(ctx context.Context, token string) error {
	query := `
		UPDATE jwt_refresh_tokens 
		SET last_used_at = NOW() 
		WHERE token_hash = $1
	`

	_, err := s.db.Exec(ctx, query, token)
	return err
}

func (s *Store) GetUserStats(ctx context.Context, userID string) (*types.UserStats, error) {
	stats := &types.UserStats{
		UserID: userID,
	}

	workoutQuery := `
		SELECT 
			COUNT(DISTINCT DATE(date)) as total_workouts,
			MIN(DATE(date)) as first_workout,
			MAX(DATE(date)) as last_workout
		FROM progress_logs
		WHERE user_id = $1
	`

	var firstWorkout, lastWorkout *time.Time
	err := s.db.QueryRow(ctx, workoutQuery, userID).Scan(
		&stats.TotalWorkouts,
		&firstWorkout,
		&lastWorkout,
	)
	if err != nil && err != pgx.ErrNoRows {
		log.Printf("Error fetching workout stats: %v", err)
	}
	stats.FirstWorkoutDate = firstWorkout
	stats.LastWorkoutDate = lastWorkout

	programQuery := `
		SELECT COUNT(*) 
		FROM generated_plans
		WHERE user_id = $1
		AND is_active = true
	`
	err = s.db.QueryRow(ctx, programQuery, userID).Scan(&stats.ActivePrograms)
	if err != nil && err != pgx.ErrNoRows {
		log.Printf("Error fetching active programs: %v", err)
	}

	daysActiveQuery := `
		SELECT COUNT(DISTINCT DATE(date))
		FROM progress_logs
		WHERE user_id = $1
	`
	err = s.db.QueryRow(ctx, daysActiveQuery, userID).Scan(&stats.DaysActive)
	if err != nil && err != pgx.ErrNoRows {
		log.Printf("Error fetching days active: %v", err)
	}

	currentStreak, longestStreak := s.calculateStreaks(ctx, userID)
	stats.CurrentStreak = currentStreak
	stats.LongestStreak = longestStreak

	if firstWorkout != nil && lastWorkout != nil {
		daysDiff := lastWorkout.Sub(*firstWorkout).Hours() / 24
		stats.TotalWeeks = int(daysDiff / 7)
	}

	if stats.TotalWeeks > 0 {
		expectedWorkouts := stats.TotalWeeks * 4 // Assuming 4 workouts per week average
		if expectedWorkouts > 0 {
			stats.CompletionRate = float64(stats.TotalWorkouts) / float64(expectedWorkouts) * 100
			if stats.CompletionRate > 100 {
				stats.CompletionRate = 100
			}
		}
	}

	coachInfo, err := s.getAssignedCoach(ctx, userID)
	if err == nil && coachInfo != nil {
		stats.AssignedCoach = coachInfo
	}

	return stats, nil
}

func (s *Store) calculateStreaks(ctx context.Context, userID string) (int, int) {
	query := `
		SELECT DISTINCT DATE(date) as workout_date
		FROM progress_logs
		WHERE user_id = $1
		ORDER BY workout_date DESC
	`

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		log.Printf("Error fetching workout dates for streak: %v", err)
		return 0, 0
	}
	defer rows.Close()

	var dates []time.Time
	for rows.Next() {
		var date time.Time
		if err := rows.Scan(&date); err != nil {
			continue
		}
		dates = append(dates, date)
	}

	if len(dates) == 0 {
		return 0, 0
	}

	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	today := time.Now().Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)

	mostRecentDate := dates[0].Truncate(24 * time.Hour)
	if !mostRecentDate.Equal(today) && !mostRecentDate.Equal(yesterday) {
		currentStreak = 0
	}

	for i := 0; i < len(dates)-1; i++ {
		daysDiff := dates[i].Sub(dates[i+1]).Hours() / 24

		if daysDiff <= 1 {
			tempStreak++
			if currentStreak > 0 && i == 0 {
				currentStreak = tempStreak
			}
		} else {
			if tempStreak > longestStreak {
				longestStreak = tempStreak
			}
			tempStreak = 1
		}
	}

	if tempStreak > longestStreak {
		longestStreak = tempStreak
	}

	if currentStreak == 1 && len(dates) > 0 {
		if !mostRecentDate.Equal(today) && !mostRecentDate.Equal(yesterday) {
			currentStreak = 0
		}
	}

	return currentStreak, longestStreak
}

func (s *Store) getAssignedCoach(ctx context.Context, userID string) (*types.CoachInfo, error) {
	
	query := `
		SELECT 
			u.id as coach_id,
			u.display_name as name,
			u.image_url as image,
			u.specialty,
			ca.created_at as assigned_at,
			COUNT(m.message_id) as total_messages
		FROM workout_profiles wp
		INNER JOIN coach_assignments ca ON ca.user_id = wp.workout_profile_id
		INNER JOIN users u ON u.id = ca.coach_id
		LEFT JOIN messages m ON m.sender_id = u.id AND m.receiver_id = $1
		WHERE wp.auth_user_id = $1
		AND ca.is_active = true
		GROUP BY u.id, u.display_name, u.image_url, u.specialty, ca.created_at
		LIMIT 1
	`

	var coachInfo types.CoachInfo
	var totalMessages int

	err := s.db.QueryRow(ctx, query, userID).Scan(
		&coachInfo.CoachID,
		&coachInfo.Name,
		&coachInfo.Image,
		&coachInfo.Specialty,
		&coachInfo.AssignedAt,
		&totalMessages,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No coach assigned
		}
		log.Printf("Error fetching assigned coach: %v", err)
		return nil, err
	}

	coachInfo.TotalMessages = totalMessages
	return &coachInfo, nil
}

func (s *Store) GetTodayWorkout(ctx context.Context, userID string) (*types.TodayWorkout, error) {

	planInfoQuery := `
		SELECT 
			gp.plan_id,
			COALESCE(pgm.algorithm_version, 'Generated Plan') as plan_name,
			gp.generated_at,
			COUNT(gpd.plan_day_id) as total_days
		FROM generated_plans gp
		LEFT JOIN plan_generation_metadata pgm ON pgm.plan_id = gp.plan_id
		LEFT JOIN generated_plan_days gpd ON gpd.plan_id = gp.plan_id
		WHERE gp.user_id = $1
		AND gp.is_active = true
		GROUP BY gp.plan_id, pgm.algorithm_version, gp.generated_at
		ORDER BY gp.generated_at DESC
		LIMIT 1
	`

	var planID int
	var planName string
	var generatedAt time.Time
	var totalDays int

	err := s.db.QueryRow(ctx, planInfoQuery, userID).Scan(&planID, &planName, &generatedAt, &totalDays)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if totalDays == 0 {
		return nil, nil // Plan has no days
	}

	
	now := time.Now()
	daysSinceStart := int(now.Sub(generatedAt).Hours() / 24)

	currentDayIndex := (daysSinceStart % totalDays) + 1

	query := `
		SELECT 
			gpd.plan_day_id,
			gpd.day_index,
			gpd.day_title,
			gpd.focus,
			gpd.is_rest
		FROM generated_plan_days gpd
		WHERE gpd.plan_id = $1
		AND gpd.day_index = $2
		LIMIT 1
	`

	var workout types.TodayWorkout
	var planDayID int

	workout.PlanID = planID
	workout.PlanName = planName

	err = s.db.QueryRow(ctx, query, planID, currentDayIndex).Scan(
		&planDayID,
		&workout.DayIndex,
		&workout.DayTitle,
		&workout.Focus,
		&workout.IsRest,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil 
		}
		return nil, err
	}

	if workout.IsRest {
		workout.TotalExercises = 0
		workout.EstimatedMinutes = 0
		workout.Exercises = []types.TodayExercise{}
		return &workout, nil
	}

	exerciseQuery := `
		SELECT 
			gpe.exercise_id,
			gpe.name,
			gpe.sets,
			gpe.reps,
			gpe.rest_seconds,
			gpe.notes
		FROM generated_plan_exercises gpe
		WHERE gpe.plan_day_id = $1
		ORDER BY gpe.exercise_order ASC
	`

	rows, err := s.db.Query(ctx, exerciseQuery, planDayID)
	if err != nil {
		log.Printf("Error fetching workout exercises: %v", err)
		return &workout, nil
	}
	defer rows.Close()

	var exercises []types.TodayExercise
	totalRestTime := 0

	for rows.Next() {
		var ex types.TodayExercise
		err := rows.Scan(
			&ex.ExerciseID,
			&ex.Name,
			&ex.Sets,
			&ex.Reps,
			&ex.RestSeconds,
			&ex.Notes,
		)
		if err != nil {
			log.Printf("Error scanning exercise: %v", err)
			continue
		}
		exercises = append(exercises, ex)
		totalRestTime += ex.RestSeconds * ex.Sets
	}

	workout.Exercises = exercises
	workout.TotalExercises = len(exercises)

	if len(exercises) > 0 {
		totalSets := 0
		for _, ex := range exercises {
			totalSets += ex.Sets
		}
		estimatedWorkTime := totalSets * 45 // 45 seconds per set
		workout.EstimatedMinutes = (estimatedWorkTime + totalRestTime) / 60
	}

	completionQuery := `
		SELECT date
		FROM progress_logs
		WHERE user_id = $1
		AND DATE(date) = CURRENT_DATE
		LIMIT 1
	`

	var completedDate time.Time
	err = s.db.QueryRow(ctx, completionQuery, userID).Scan(&completedDate)
	if err == nil {
		workout.IsCompleted = true
		workout.CompletedAt = &completedDate
	}

	return &workout, nil
}

func (s *Store) SaveWorkoutCompletion(ctx context.Context, userID string, completion *types.WorkoutCompletionRequest) (*types.WorkoutCompletionResponse, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	totalSets := len(completion.Exercises)
	completedSets := 0
	totalVolume := 0.0

	for _, exercise := range completion.Exercises {
		if exercise.Completed {
			completedSets++
			totalVolume += float64(exercise.Reps) * exercise.Weight
		}

		_, err := tx.Exec(ctx, `
			INSERT INTO progress_logs (user_id, exercise_id, date, sets_completed, reps_completed, weight_used, duration_seconds, notes)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, userID, exercise.ExerciseID, completion.CompletedAt, 1, exercise.Reps, exercise.Weight, 0, exercise.Notes)

		if err != nil {
			log.Printf("Error inserting progress log for exercise %s: %v", exercise.ExerciseName, err)
		}
	}

	completionRate := 0.0
	if totalSets > 0 {
		completionRate = (float64(completedSets) / float64(totalSets)) * 100
	}

	currentStreak, _ := s.calculateStreaks(ctx, userID)

	isPersonalBest := false
	var maxVolume float64
	err = s.db.QueryRow(ctx, `
		SELECT COALESCE(MAX(total_volume), 0) as max_volume
		FROM (
			SELECT SUM(reps_completed * weight_used) as total_volume
			FROM progress_logs
			WHERE user_id = $1
			AND date < $2
			GROUP BY date
		) daily_volumes
	`, userID, completion.CompletedAt).Scan(&maxVolume)

	if err == nil && totalVolume > maxVolume {
		isPersonalBest = true
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &types.WorkoutCompletionResponse{
		Success:         true,
		Message:         "Workout saved successfully",
		WorkoutDate:     completion.CompletedAt,
		TotalSets:       totalSets,
		CompletedSets:   completedSets,
		CompletionRate:  completionRate,
		TotalVolume:     totalVolume,
		DurationMinutes: completion.DurationSeconds / 60,
		NewStreak:       currentStreak,
		IsPersonalBest:  isPersonalBest,
	}, nil
}

func (s *Store) GetActivityFeed(ctx context.Context, userID string, limit int) ([]types.ActivityFeedItem, error) {
	if limit <= 0 {
		limit = 10
	}

	activities := []types.ActivityFeedItem{}

	workoutQuery := `
		SELECT DISTINCT
			date,
			COUNT(DISTINCT exercise_id) as exercise_count,
			SUM(reps_completed * weight_used) as total_volume
		FROM progress_logs
		WHERE user_id = $1
		GROUP BY date
		ORDER BY date DESC
		LIMIT $2
	`

	rows, err := s.db.Query(ctx, workoutQuery, userID, limit)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var date time.Time
			var exerciseCount int
			var totalVolume float64

			if err := rows.Scan(&date, &exerciseCount, &totalVolume); err == nil {
				activities = append(activities, types.ActivityFeedItem{
					ID:          date.Format("2006-01-02") + "-workout",
					Type:        types.ActivityWorkoutCompleted,
					Title:       "Workout Completed",
					Description: fmt.Sprintf("Completed %d exercises with %.0f lbs total volume", exerciseCount, totalVolume),
					Timestamp:   date,
					Icon:        "fitness",
					Metadata: map[string]interface{}{
						"exercise_count": exerciseCount,
						"total_volume":   totalVolume,
					},
				})
			}
		}
	}

	currentStreak, longestStreak := s.calculateStreaks(ctx, userID)

	if currentStreak > 0 && (currentStreak%7 == 0 || currentStreak%30 == 0 || currentStreak == 100) {
		milestoneDate := time.Now()
		activities = append(activities, types.ActivityFeedItem{
			ID:          fmt.Sprintf("streak-%d", currentStreak),
			Type:        types.ActivityStreakMilestone,
			Title:       "Streak Milestone!",
			Description: fmt.Sprintf("You've reached a %d-day workout streak! Keep it up! 🔥", currentStreak),
			Timestamp:   milestoneDate,
			Icon:        "flame",
			Metadata: map[string]interface{}{
				"current_streak": currentStreak,
				"longest_streak": longestStreak,
			},
		})
	}

	planQuery := `
		SELECT 
			gp.plan_id,
			gp.generated_at,
			pgm.algorithm_version,
			pgm.weekly_frequency
		FROM generated_plans gp
		LEFT JOIN plan_generation_metadata pgm ON pgm.plan_id = gp.plan_id
		WHERE gp.user_id = $1
		ORDER BY gp.generated_at DESC
		LIMIT 3
	`

	planRows, err := s.db.Query(ctx, planQuery, userID)
	if err == nil {
		defer planRows.Close()
		for planRows.Next() {
			var planID int
			var generatedAt time.Time
			var algorithmVersion *string
			var weeklyFrequency *int

			if err := planRows.Scan(&planID, &generatedAt, &algorithmVersion, &weeklyFrequency); err == nil {
				version := "Custom Plan"
				if algorithmVersion != nil {
					version = *algorithmVersion
				}

				freq := 3
				if weeklyFrequency != nil {
					freq = *weeklyFrequency
				}

				activities = append(activities, types.ActivityFeedItem{
					ID:          fmt.Sprintf("plan-%d", planID),
					Type:        types.ActivityNewPlan,
					Title:       "New Workout Plan",
					Description: fmt.Sprintf("Generated %s with %d workouts per week", version, freq),
					Timestamp:   generatedAt,
					Icon:        "calendar",
					Metadata: map[string]interface{}{
						"plan_id":          planID,
						"algorithm":        version,
						"weekly_frequency": freq,
					},
				})
			}
		}
	}

	prQuery := `
		WITH daily_volumes AS (
			SELECT 
				date,
				exercise_id,
				MAX(weight_used) as max_weight,
				SUM(reps_completed * weight_used) as total_volume
			FROM progress_logs
			WHERE user_id = $1
			AND date >= NOW() - INTERVAL '30 days'
			GROUP BY date, exercise_id
		),
		historical_max AS (
			SELECT 
				exercise_id,
				MAX(weight_used) as historical_max_weight
			FROM progress_logs
			WHERE user_id = $1
			AND date < NOW() - INTERVAL '30 days'
			GROUP BY exercise_id
		)
		SELECT 
			dv.date,
			dv.exercise_id,
			dv.max_weight,
			COALESCE(hm.historical_max_weight, 0) as prev_max
		FROM daily_volumes dv
		LEFT JOIN historical_max hm ON hm.exercise_id = dv.exercise_id
		WHERE dv.max_weight > COALESCE(hm.historical_max_weight, 0)
		ORDER BY dv.date DESC
		LIMIT 3
	`

	prRows, err := s.db.Query(ctx, prQuery, userID)
	if err == nil {
		defer prRows.Close()
		for prRows.Next() {
			var date time.Time
			var exerciseID *int
			var maxWeight, prevMax float64

			if err := prRows.Scan(&date, &exerciseID, &maxWeight, &prevMax); err == nil {
				activities = append(activities, types.ActivityFeedItem{
					ID:          fmt.Sprintf("pr-%s-%d", date.Format("2006-01-02"), exerciseID),
					Type:        types.ActivityPRChieved,
					Title:       "Personal Record!",
					Description: fmt.Sprintf("New PR: %.0f lbs (previous: %.0f lbs) 🏆", maxWeight, prevMax),
					Timestamp:   date,
					Icon:        "trophy",
					Metadata: map[string]interface{}{
						"exercise_id": exerciseID,
						"new_weight":  maxWeight,
						"prev_weight": prevMax,
					},
				})
			}
		}
	}

	sort.Slice(activities, func(i, j int) bool {
		return activities[i].Timestamp.After(activities[j].Timestamp)
	})

	if len(activities) > limit {
		activities = activities[:limit]
	}

	return activities, nil
}

func (s *Store) GetWorkoutHistory(ctx context.Context, userID string, startDate, endDate *time.Time, page, pageSize int) (*types.WorkoutHistoryResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	dateFilter := ""
	args := []interface{}{userID}
	argIndex := 2

	if startDate != nil {
		dateFilter += fmt.Sprintf(" AND DATE(pl.date) >= $%d", argIndex)
		args = append(args, startDate.Format("2006-01-02"))
		argIndex++
	}
	if endDate != nil {
		dateFilter += fmt.Sprintf(" AND DATE(pl.date) <= $%d", argIndex)
		args = append(args, endDate.Format("2006-01-02"))
		argIndex++
	}

	query := fmt.Sprintf(`
		SELECT 
			DATE(pl.date) as workout_date,
			COUNT(DISTINCT pl.exercise_id) as total_exercises,
			COUNT(*) as completed_sets,
			COALESCE(SUM(pl.reps_completed * pl.weight_used), 0) as total_volume,
			COALESCE(ROUND(AVG(pl.duration_seconds)::numeric / 60), 0) as avg_duration
		FROM progress_logs pl
		WHERE pl.user_id = $1
		%s
		GROUP BY DATE(pl.date)
		ORDER BY DATE(pl.date) DESC
		LIMIT $%d OFFSET $%d
	`, dateFilter, argIndex, argIndex+1)

	args = append(args, pageSize, offset)

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		log.Printf("Error querying workout history: %v", err)
		return nil, err
	}
	defer rows.Close()

	workouts := []types.WorkoutHistoryItem{}
	for rows.Next() {
		var workout types.WorkoutHistoryItem
		var durationMinutes float64

		err := rows.Scan(
			&workout.Date,
			&workout.TotalExercises,
			&workout.CompletedSets,
			&workout.TotalVolume,
			&durationMinutes,
		)
		if err != nil {
			log.Printf("Error scanning workout history: %v", err)
			continue
		}

		workout.DurationMinutes = int(durationMinutes)

		exerciseQuery := `
			SELECT DISTINCT e.name
			FROM progress_logs pl
			LEFT JOIN exercises e ON e.exercise_id = pl.exercise_id
			WHERE pl.user_id = $1 AND pl.date = $2
			ORDER BY e.name
		`
		exerciseRows, err := s.db.Query(ctx, exerciseQuery, userID, workout.Date)
		if err == nil {
			exercises := []string{}
			for exerciseRows.Next() {
				var name *string
				if err := exerciseRows.Scan(&name); err == nil && name != nil {
					exercises = append(exercises, *name)
				}
			}
			exerciseRows.Close()
			workout.Exercises = exercises
		}

		workouts = append(workouts, workout)
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT date)
		FROM progress_logs
		WHERE user_id = $1
		%s
	`, dateFilter)

	var totalCount int
	countArgs := []interface{}{userID}
	if startDate != nil {
		countArgs = append(countArgs, *startDate)
	}
	if endDate != nil {
		countArgs = append(countArgs, *endDate)
	}

	err = s.db.QueryRow(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		totalCount = len(workouts)
	}

	return &types.WorkoutHistoryResponse{
		Workouts:   workouts,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		HasMore:    (page * pageSize) < totalCount,
	}, nil
}

func (s *Store) GetExerciseProgress(ctx context.Context, userID string, exerciseName string, startDate, endDate *time.Time) (*types.ExerciseProgressData, error) {
	dateFilter := ""
	args := []interface{}{userID, exerciseName}
	argIndex := 3

	if startDate != nil {
		dateFilter += fmt.Sprintf(" AND pl.date >= $%d", argIndex)
		args = append(args, *startDate)
		argIndex++
	}
	if endDate != nil {
		dateFilter += fmt.Sprintf(" AND pl.date <= $%d", argIndex)
		args = append(args, *endDate)
		argIndex++
	}

	query := fmt.Sprintf(`
		WITH exercise_data AS (
			SELECT 
				pl.date,
				pl.exercise_id,
				MAX(pl.weight_used) as max_weight,
				MAX(pl.reps_completed) as max_reps,
				COUNT(*) as sets,
				SUM(pl.reps_completed * pl.weight_used) as volume
			FROM progress_logs pl
			LEFT JOIN exercises e ON e.exercise_id = pl.exercise_id
			WHERE pl.user_id = $1
			AND (e.name = $2 OR e.name IS NULL)
			%s
			GROUP BY pl.date, pl.exercise_id
			ORDER BY pl.date ASC
		),
		max_values AS (
			SELECT 
				MAX(max_weight) as overall_max_weight,
				MAX(volume) as overall_max_volume,
				SUM(sets) as total_sets
			FROM exercise_data
		)
		SELECT 
			ed.date,
			ed.exercise_id,
			ed.max_weight,
			ed.max_reps,
			ed.sets,
			ed.volume,
			mv.overall_max_weight,
			mv.overall_max_volume,
			mv.total_sets
		FROM exercise_data ed
		CROSS JOIN max_values mv
		ORDER BY ed.date ASC
	`, dateFilter)

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progressData types.ExerciseProgressData
	progressData.ExerciseName = exerciseName
	progressData.DataPoints = []types.ExerciseProgressDataPoint{}

	var previousMaxWeight float64 = 0

	for rows.Next() {
		var point types.ExerciseProgressDataPoint
		var exerciseID *int

		err := rows.Scan(
			&point.Date,
			&exerciseID,
			&point.Weight,
			&point.Reps,
			&point.Sets,
			&point.Volume,
			&progressData.MaxWeight,
			&progressData.MaxVolume,
			&progressData.TotalSets,
		)
		if err != nil {
			log.Printf("Error scanning exercise progress: %v", err)
			continue
		}

		if exerciseID != nil && progressData.ExerciseID == nil {
			progressData.ExerciseID = exerciseID
		}

		if point.Weight > previousMaxWeight {
			point.IsPersonalRecord = true
			previousMaxWeight = point.Weight
		}

		progressData.DataPoints = append(progressData.DataPoints, point)
	}

	return &progressData, nil
}
