import apiClient from './client';
import { ApiResponse } from '@/types/api';
import { Task, CreateTaskRequest, UpdateTaskRequest } from '@/types/task';

export const tasksApi = {
  list: (projectId: string, params?: object) =>
    apiClient.get<ApiResponse<Task[]>>(`/projects/${projectId}/tasks`, { params }),
  getById: (id: string) =>
    apiClient.get<ApiResponse<Task>>(`/tasks/${id}`),
  create: (projectId: string, data: CreateTaskRequest) =>
    apiClient.post<ApiResponse<Task>>(`/projects/${projectId}/tasks`, data),
  update: (id: string, data: UpdateTaskRequest) =>
    apiClient.put<ApiResponse<Task>>(`/tasks/${id}`, data),
  delete: (id: string) =>
    apiClient.delete(`/tasks/${id}`),
  updateStatus: (id: string, status: string) =>
    apiClient.patch(`/tasks/${id}/status`, { status }),
};
