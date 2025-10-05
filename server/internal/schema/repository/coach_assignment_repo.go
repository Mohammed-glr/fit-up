package repository

import (
	"context"

	"github.com/tdmdh/fit-up-server/internal/schema/types"
)

func (s *Store) CreateCoachAssignment(ctx context.Context, req *types.CoachAssignmentRequest) (*types.CoachAssignment, error) {
	query := `
		INSERT INTO coach_assignments (coach_id, user_id, assigned_by, notes)
		VALUES ($1, $2, $3, $4)
		RETURNING assignment_id, coach_id, user_id, assigned_at, assigned_by, is_active, notes
	`

	var assignment types.CoachAssignment
	err := s.db.QueryRow(ctx, query,
		req.CoachID,
		req.UserID,
		req.CoachID,
		req.Notes,
	).Scan(
		&assignment.AssignmentID,
		&assignment.CoachID,
		&assignment.UserID,
		&assignment.AssignedAt,
		&assignment.AssignedBy,
		&assignment.IsActive,
		&assignment.Notes,
	)

	return &assignment, err
}

func (s *Store) GetClientsByCoachID(ctx context.Context, coachID string) ([]types.ClientSummary, error) {
	query := `
		SELECT 
			ca.user_id,
			wp.auth_user_id,
			'FirstName' as first_name,  -- TODO: Add user profile table
			'LastName' as last_name,
			'email@example.com' as email,
			ca.assigned_at,
			ws.schema_id,
			COUNT(DISTINCT fg.goal_id) as active_goals,
			AVG(CASE WHEN sess.status = 'completed' THEN 1.0 ELSE 0.0 END) as completion_rate,
			MAX(sess.start_time) as last_workout_date,
			COUNT(DISTINCT sess.session_id) as total_workouts,
			0 as current_streak,  -- TODO: Calculate streak
			wp.level as fitness_level
		FROM coach_assignments ca
		LEFT JOIN workout_profiles wp ON ca.user_id = wp.workout_profile_id
		LEFT JOIN weekly_schemas ws ON ca.user_id = ws.user_id AND ws.active = true
		LEFT JOIN fitness_goals fg ON ca.user_id = fg.user_id AND fg.is_active = true
		LEFT JOIN workout_sessions sess ON ca.user_id = sess.user_id
		WHERE ca.coach_id = $1 AND ca.is_active = true
		GROUP BY ca.user_id, wp.auth_user_id, ca.assigned_at, ws.schema_id, wp.level
		ORDER BY ca.assigned_at DESC
	`

	rows, err := s.db.Query(ctx, query, coachID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []types.ClientSummary
	for rows.Next() {
		var client types.ClientSummary
		err := rows.Scan(
			&client.UserID,
			&client.AuthID,
			&client.FirstName,
			&client.LastName,
			&client.Email,
			&client.AssignedAt,
			&client.CurrentSchemaID,
			&client.ActiveGoals,
			&client.CompletionRate,
			&client.LastWorkoutDate,
			&client.TotalWorkouts,
			&client.CurrentStreak,
			&client.FitnessLevel,
		)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	return clients, nil
}

func (s *Store) IsCoachForUser(ctx context.Context, coachID string, userID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM coach_assignments 
			WHERE coach_id = $1 AND user_id = $2 AND is_active = true
		)
	`

	var exists bool
	err := s.db.QueryRow(ctx, query, coachID, userID).Scan(&exists)
	return exists, err
}

func (s *Store) GetCoachDashboard(ctx context.Context, coachID string) (*types.CoachDashboard, error) {
	// Implementation for complete dashboard data
	// TODO: Implement full dashboard aggregation

	clients, err := s.GetClientsByCoachID(ctx, coachID)
	if err != nil {
		return nil, err
	}

	dashboard := &types.CoachDashboard{
		CoachID:      coachID,
		TotalClients: len(clients),
		Clients:      clients,
	}

	return dashboard, nil
}

func (s *Store) LogCoachActivity(ctx context.Context, activity *types.CoachActivity) error {
	query := `
		INSERT INTO coach_activity_log (coach_id, user_id, activity_type, description)
		VALUES ($1, $2, $3, $4)
	`

	_, err := s.db.Exec(ctx, query,
		activity.CoachID,
		activity.UserID,
		activity.ActivityType,
		activity.Description,
	)

	return err
}

func (s *Store) DeactivateAssignment(ctx context.Context, assignmentID int) error {
	query := `
		UPDATE coach_assignments
		SET is_active = false
		WHERE assignment_id = $1
	`
	_, err := s.db.Exec(ctx, query, assignmentID)
	return err
}

func (s *Store) GetCoachActivityLog(ctx context.Context, coachID string, limit int) ([]types.CoachActivity, error) {
	query := `
		SELECT activity_id, coach_id, user_id, activity_type, description, created_at
		FROM coach_activity_log
		WHERE coach_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := s.db.Query(ctx, query, coachID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []types.CoachActivity
	for rows.Next() {
		var activity types.CoachActivity
		err := rows.Scan(
			&activity.ActivityID,
			&activity.CoachID,
			&activity.UserID,
			&activity.ActivityType,
			&activity.Description,
			&activity.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

func (s *Store) GetCoachAssignment(ctx context.Context, assignmentID int) (*types.CoachAssignment, error) {
	query := `
		SELECT assignment_id, coach_id, user_id, assigned_at, assigned_by, is_active, notes, deactivated_at
		FROM coach_assignments
		WHERE assignment_id = $1
	`
	var assignment types.CoachAssignment
	err := s.db.QueryRow(ctx, query, assignmentID).Scan(
		&assignment.AssignmentID,
		&assignment.CoachID,
		&assignment.UserID,
		&assignment.AssignedAt,
		&assignment.AssignedBy,
		&assignment.IsActive,
		&assignment.Notes,
		&assignment.DeactivatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

func (s *Store) GetCoachByUserID(ctx context.Context, userID int) (*types.CoachAssignment, error) {
	query := `
		SELECT assignment_id, coach_id, user_id, assigned_at, assigned_by, is_active, notes, deactivated_at
		FROM coach_assignments
		WHERE user_id = $1 AND is_active = true
	`
	var assignment types.CoachAssignment
	err := s.db.QueryRow(ctx, query, userID).Scan(
		&assignment.AssignmentID,
		&assignment.CoachID,
		&assignment.UserID,
		&assignment.AssignedAt,
		&assignment.AssignedBy,
		&assignment.IsActive,
		&assignment.Notes,
		&assignment.DeactivatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}
