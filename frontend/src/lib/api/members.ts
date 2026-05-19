import apiClient from './client';
import { ApiResponse } from '@/types/api';
import { Member, AddMemberRequest, UpdateMemberRoleRequest } from '@/types/member';

export const membersApi = {
  list: (projectId: string) =>
    apiClient.get<ApiResponse<Member[]>>(`/projects/${projectId}/members`),
  add: (projectId: string, data: AddMemberRequest) =>
    apiClient.post<ApiResponse<Member>>(`/projects/${projectId}/members`, data),
  updateRole: (projectId: string, userId: string, data: UpdateMemberRoleRequest) =>
    apiClient.patch(`/projects/${projectId}/members/${userId}`, data),
  remove: (projectId: string, userId: string) =>
    apiClient.delete(`/projects/${projectId}/members/${userId}`),
};
