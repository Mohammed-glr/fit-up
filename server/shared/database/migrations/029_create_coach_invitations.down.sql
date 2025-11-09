-- Drop coach_invitations table
DROP INDEX IF EXISTS idx_coach_invitations_expires_at;
DROP INDEX IF EXISTS idx_coach_invitations_status;
DROP INDEX IF EXISTS idx_coach_invitations_token;
DROP INDEX IF EXISTS idx_coach_invitations_email;
DROP INDEX IF EXISTS idx_coach_invitations_coach_id;

DROP TABLE IF EXISTS coach_invitations;
