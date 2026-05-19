import apiClient from './client';
import { ApiResponse } from '@/types/api';
import { DashboardData, ProjectReport } from '@/types/report';

export const reportsApi = {
  dashboard: () =>
    apiClient.get<ApiResponse<DashboardData>>('/reports/dashboard'),
  projectReport: (projectId: string) =>
    apiClient.get<ApiResponse<ProjectReport>>(`/reports/projects/${projectId}`),
  projectsByStatus: () =>
    apiClient.get<ApiResponse<Record<string, number>>>('/reports/projects/by-status'),
  tasksByStatus: () =>
    apiClient.get<ApiResponse<Record<string, number>>>('/reports/tasks/by-status'),
};
