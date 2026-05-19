import { Project } from './project';

export interface DashboardData {
  total_projects: number;
  active_projects: number;
  completed_projects: number;
  total_tasks: number;
  tasks_by_status: Record<string, number>;
  projects_by_status: Record<string, number>;
  projects_by_priority: Record<string, number>;
  recent_projects: Project[];
}

export interface ProjectReport {
  project: Project;
  total_tasks: number;
  completed_tasks: number;
  progress: number;
  status_history: StatusHistoryEntry[];
}

export interface StatusHistoryEntry {
  old_status: string;
  new_status: string;
  changed_by: string;
  changed_at: string;
}
