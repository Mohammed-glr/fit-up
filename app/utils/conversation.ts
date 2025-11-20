import type { UserRole } from '@/types/auth';

export const getInitial = (value: string) => {
  return value.trim().charAt(0)?.toUpperCase() || 'U';
};

export const formatRole = (role?: UserRole) => {
  if (!role) {
    return 'Member';
  }

  switch (role) {
    case 'coach':
      return 'Coach';
    case 'user':
      return 'Client';
    case 'admin':
      return 'Admin';
  }

  return 'Member';
};

export const canMessage = (currentRole?: UserRole, targetRole?: UserRole) => {
  if (!currentRole || !targetRole) {
    return false;
  }

  // Allow anyone to message anyone
  return true;
};

export const roleRestrictionMessage = (role?: UserRole) => {
  if (role === 'coach') {
    return 'Coaches can only message client accounts. Try searching for a client username.';
  }

  if (role === 'user') {
    return 'Clients can only message coach accounts. Try searching for a coach username.';
  }

  return 'Update your role to coach or client to start conversations.';
};
