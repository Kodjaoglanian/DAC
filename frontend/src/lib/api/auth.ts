import apiClient from './client';
import { ApiResponse } from '@/types/api';
import { LoginRequest, RegisterRequest, AuthResponse, User } from '@/types/auth';

export const authApi = {
  login: (data: LoginRequest) =>
    apiClient.post<ApiResponse<AuthResponse>>('/auth/login', data),
  register: (data: RegisterRequest) =>
    apiClient.post<ApiResponse<AuthResponse>>('/auth/register', data),
  me: () =>
    apiClient.get<ApiResponse<User>>('/auth/me'),
};
