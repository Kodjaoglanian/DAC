import { User } from './auth';

export type ProjectStatus = 'planning' | 'in_progress' | 'completed' | 'cancelled';
export type Priority = 'low' | 'medium' | 'high' | 'critical';

export interface Project {
  id: string;
  name: string;
  description: string;
  status: ProjectStatus;
  priority: Priority;
  start_date: string | null;
  end_date: string | null;
  created_by: User;
  members: ProjectMemberSummary[];
  tasks_count: number;
  created_at: string;
  updated_at: string;
}

export interface ProjectMemberSummary {
  id: string;
  user_id: string;
  name: string;
  role: string;
  joined_at: string;
}

export interface CreateProjectRequest {
  name: string;
  description: string;
  priority: Priority;
  start_date?: string;
  end_date?: string;
}

export interface UpdateProjectRequest {
  name?: string;
  description?: string;
  status?: ProjectStatus;
  priority?: Priority;
  start_date?: string;
  end_date?: string;
}

export interface ProjectFilters {
  page?: number;
  page_size?: number;
  status?: ProjectStatus;
  priority?: Priority;
  search?: string;
  created_by?: string;
  sort?: 'created_at' | 'name' | 'priority' | 'status';
  order?: 'ASC' | 'DESC';
}
