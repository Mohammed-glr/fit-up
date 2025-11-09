import type { ClientSummary } from '@/types/schema';

export type SortOption = 'name' | 'last_active' | 'join_date' | 'workouts' | 'streak';
export type SortOrder = 'asc' | 'desc';

export interface ClientSort {
  option: SortOption;
  order: SortOrder;
}

export const SORT_OPTIONS = [
  { key: 'name' as SortOption, label: 'Name', icon: 'person' },
  { key: 'last_active' as SortOption, label: 'Last Active', icon: 'time' },
  { key: 'join_date' as SortOption, label: 'Join Date', icon: 'calendar' },
  { key: 'workouts' as SortOption, label: 'Total Workouts', icon: 'barbell' },
  { key: 'streak' as SortOption, label: 'Current Streak', icon: 'flame' },
] as const;

/**
 * Sort clients based on the selected option and order
 */
export function sortClients(
  clients: ClientSummary[],
  sortBy: SortOption,
  sortOrder: SortOrder
): ClientSummary[] {
  const sorted = [...clients].sort((a, b) => {
    let comparison = 0;

    switch (sortBy) {
      case 'name':
        const nameA = `${a.first_name} ${a.last_name}`.toLowerCase();
        const nameB = `${b.first_name} ${b.last_name}`.toLowerCase();
        comparison = nameA.localeCompare(nameB);
        break;

      case 'last_active':
        // Sort by last workout date (most recent first when desc)
        const dateA = a.last_workout_date ? new Date(a.last_workout_date).getTime() : 0;
        const dateB = b.last_workout_date ? new Date(b.last_workout_date).getTime() : 0;
        comparison = dateA - dateB;
        break;

      case 'join_date':
        // Sort by assigned date
        const assignedA = new Date(a.assigned_at).getTime();
        const assignedB = new Date(b.assigned_at).getTime();
        comparison = assignedA - assignedB;
        break;

      case 'workouts':
        comparison = a.total_workouts - b.total_workouts;
        break;

      case 'streak':
        comparison = a.current_streak - b.current_streak;
        break;

      default:
        comparison = 0;
    }

    return sortOrder === 'asc' ? comparison : -comparison;
  });

  return sorted;
}

/**
 * Get a human-readable description of the current sort
 */
export function getSortDescription(sortBy: SortOption, sortOrder: SortOrder): string {
  const option = SORT_OPTIONS.find((opt) => opt.key === sortBy);
  if (!option) return '';

  const direction = sortOrder === 'asc' ? '↑' : '↓';
  return `${option.label} ${direction}`;
}
