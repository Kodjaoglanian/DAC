import apiClient from './client';
import { ApiResponse, PaginatedResponse } from '@/types/api';
import { Project, CreateProjectRequest, UpdateProjectRequest, ProjectFilters } from '@/types/project';

export const projectsApi = {
  list: (filters?: ProjectFilters) =>
    apiClient.get<ApiResponse<PaginatedResponse<Project>>>('/projects', { params: filters }),
  getById: (id: string) =>
    apiClient.get<ApiResponse<Project>>(`/projects/${id}`),
  create: (data: CreateProjectRequest) =>
    apiClient.post<ApiResponse<Project>>('/projects', data),
  update: (id: string, data: UpdateProjectRequest) =>
    apiClient.put<ApiResponse<Project>>(`/projects/${id}`, data),
  delete: (id: string) =>
    apiClient.delete(`/projects/${id}`),
  updateStatus: (id: string, status: string) =>
    apiClient.patch(`/projects/${id}/status`, { status }),
};
