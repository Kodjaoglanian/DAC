import { User } from './auth';

export type ProjectMemberRole = 'owner' | 'manager' | 'member';

export interface Member {
  id: string;
  user: User;
  role: ProjectMemberRole;
  joined_at: string;
}

export interface AddMemberRequest {
  user_id: string;
  role: ProjectMemberRole;
}

export interface UpdateMemberRoleRequest {
  role: ProjectMemberRole;
}
