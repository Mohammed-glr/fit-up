import type { ClientSummary } from '@/types/schema';

export type ClientStatus = 'active' | 'needs_attention' | 'inactive' | 'no_schema';

export interface ClientStatusInfo {
  status: ClientStatus;
  label: string;
  color: string;
  icon: string;
  description: string;
}

/**
 * Calculate client status based on last workout date and schema presence
 * 
 * Rules:
 * - No Schema: Client has no active workout plan
 * - Active: Worked out in last 3 days
 * - Needs Attention: No workout in 3-7 days
 * - Inactive: No workout in 7+ days or never worked out
 */
export function getClientStatus(client: ClientSummary): ClientStatusInfo {
  // Check if client has an active schema
  if (!client.current_schema_id) {
    return {
      status: 'no_schema',
      label: 'No Schema',
      color: '#8B5CF6', // Purple
      icon: 'alert-circle-outline',
      description: 'Client needs a workout plan assigned',
    };
  }

  // Check last workout date
  if (!client.last_workout_date) {
    return {
      status: 'inactive',
      label: 'Inactive',
      color: '#EF4444', // Red
      icon: 'close-circle-outline',
      description: 'Client has never worked out',
    };
  }

  const lastWorkout = new Date(client.last_workout_date);
  const now = new Date();
  const daysSinceLastWorkout = Math.floor(
    (now.getTime() - lastWorkout.getTime()) / (1000 * 60 * 60 * 24)
  );

  if (daysSinceLastWorkout <= 3) {
    return {
      status: 'active',
      label: 'Active',
      color: '#10B981', // Green
      icon: 'checkmark-circle',
      description: 'Client is actively training',
    };
  }

  if (daysSinceLastWorkout <= 7) {
    return {
      status: 'needs_attention',
      label: 'Needs Attention',
      color: '#F59E0B', // Orange/Yellow
      icon: 'warning-outline',
      description: `Last workout ${daysSinceLastWorkout} days ago`,
    };
  }

  return {
    status: 'inactive',
    label: 'Inactive',
    color: '#EF4444', // Red
    icon: 'close-circle-outline',
    description: `Last workout ${daysSinceLastWorkout} days ago`,
  };
}

/**
 * Get a human-readable description of how long ago the last workout was
 */
export function getLastWorkoutDescription(lastWorkoutDate: string | null): string {
  if (!lastWorkoutDate) {
    return 'Never worked out';
  }

  const lastWorkout = new Date(lastWorkoutDate);
  const now = new Date();
  const daysSince = Math.floor(
    (now.getTime() - lastWorkout.getTime()) / (1000 * 60 * 60 * 24)
  );

  if (daysSince === 0) {
    return 'Today';
  }

  if (daysSince === 1) {
    return 'Yesterday';
  }

  if (daysSince <= 7) {
    return `${daysSince} days ago`;
  }

  if (daysSince <= 30) {
    const weeks = Math.floor(daysSince / 7);
    return `${weeks} ${weeks === 1 ? 'week' : 'weeks'} ago`;
  }

  const months = Math.floor(daysSince / 30);
  return `${months} ${months === 1 ? 'month' : 'months'} ago`;
}

/**
 * Filter clients by status
 */
export function filterClientsByStatus(
  clients: ClientSummary[],
  statusFilter: ClientStatus | 'all'
): ClientSummary[] {
  if (statusFilter === 'all') {
    return clients;
  }

  return clients.filter((client) => {
    const { status } = getClientStatus(client);
    return status === statusFilter;
  });
}

/**
 * Get status counts for filter chips
 */
export function getStatusCounts(clients: ClientSummary[]): Record<ClientStatus | 'all', number> {
  const counts = {
    all: clients.length,
    active: 0,
    needs_attention: 0,
    inactive: 0,
    no_schema: 0,
  };

  clients.forEach((client) => {
    const { status } = getClientStatus(client);
    counts[status]++;
  });

  return counts;
}
