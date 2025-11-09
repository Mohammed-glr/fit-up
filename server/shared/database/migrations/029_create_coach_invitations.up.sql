-- Create coach_invitations table for client invitation system
CREATE TABLE IF NOT EXISTS coach_invitations (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    coach_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    first_name TEXT,
    last_name TEXT,
    invitation_token TEXT UNIQUE NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'expired', 'cancelled')),
    custom_message TEXT,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    accepted_at TIMESTAMP,
    accepted_by_user_id TEXT REFERENCES users(id) ON DELETE SET NULL,
    
    CONSTRAINT unique_coach_email UNIQUE(coach_id, email)
);

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_coach_invitations_coach_id ON coach_invitations(coach_id);
CREATE INDEX IF NOT EXISTS idx_coach_invitations_email ON coach_invitations(email);
CREATE INDEX IF NOT EXISTS idx_coach_invitations_token ON coach_invitations(invitation_token);
CREATE INDEX IF NOT EXISTS idx_coach_invitations_status ON coach_invitations(status);
CREATE INDEX IF NOT EXISTS idx_coach_invitations_expires_at ON coach_invitations(expires_at);

-- Add comment for documentation
COMMENT ON TABLE coach_invitations IS 'Stores coach invitations sent to potential clients';
COMMENT ON COLUMN coach_invitations.invitation_token IS 'Unique token used to accept the invitation';
COMMENT ON COLUMN coach_invitations.status IS 'Status of invitation: pending, accepted, expired, or cancelled';
COMMENT ON COLUMN coach_invitations.expires_at IS 'Invitation expiration timestamp (default 7 days from creation)';
