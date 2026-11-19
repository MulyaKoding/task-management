export type TaskStatus = 'todo' | 'in_progress' | 'done';

export interface Task {
  id: number;
  title: string;
  description?: string | null;
  status: TaskStatus;
  assignee?: string | null;
  created_at: string;
  updated_at: string;
  deleted_at?: string | null;
}

export interface PaginationMeta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface TaskFilterParams {
  status?: string;
  keyword?: string;
  assignee?: string;
  page?: number;
  limit?: number;
  sort?: string;
}

export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  meta?: PaginationMeta;
  message?: string;
  error?: {
    code: string;
    message: string;
    details?: any;
  };
}

export class ApiError extends Error {
  code: string;
  status: number;
  details?: any;

  constructor(status: number, message: string, code = 'ERROR', details?: any) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.code = code;
    this.details = details;
  }
}

// In Expo web/mobile, localhost is used for local backend
const API_BASE_URL = process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080/api';

export const TaskApi = {
  async fetchTasks(params: TaskFilterParams = {}): Promise<{ data: Task[]; meta: PaginationMeta }> {
    const query = new URLSearchParams();
    if (params.status) query.append('status', params.status);
    if (params.keyword) query.append('keyword', params.keyword);
    if (params.assignee) query.append('assignee', params.assignee);
    if (params.page) query.append('page', params.page.toString());
    if (params.limit) query.append('limit', params.limit.toString());
    if (params.sort) query.append('sort', params.sort);

    const url = `${API_BASE_URL}/tasks?${query.toString()}`;
    const res = await fetch(url);
    const json: ApiResponse<Task[]> = await res.json();

    if (!res.ok || !json.success) {
      throw new ApiError(
        res.status,
        json.error?.message || 'Failed to fetch tasks',
        json.error?.code || 'FETCH_ERROR',
        json.error?.details
      );
    }

    return {
      data: json.data || [],
      meta: json.meta || { page: 1, limit: 10, total: 0, total_pages: 0 },
    };
  },

  async createTask(data: {
    title: string;
    description?: string;
    status?: TaskStatus;
    assignee?: string;
  }): Promise<Task> {
    const res = await fetch(`${API_BASE_URL}/tasks`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });

    const json: ApiResponse<Task> = await res.json();
    if (!res.ok || !json.success) {
      throw new ApiError(
        res.status,
        json.error?.message || 'Failed to create task',
        json.error?.code || 'CREATE_ERROR',
        json.error?.details
      );
    }

    return json.data!;
  },

  async updateTask(
    id: number,
    data: {
      title?: string;
      description?: string;
      status?: TaskStatus;
      assignee?: string;
    }
  ): Promise<Task> {
    const res = await fetch(`${API_BASE_URL}/tasks/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });

    const json: ApiResponse<Task> = await res.json();
    if (!res.ok || !json.success) {
      throw new ApiError(
        res.status,
        json.error?.message || 'Failed to update task',
        json.error?.code || 'UPDATE_ERROR',
        json.error?.details
      );
    }

    return json.data!;
  },

  async deleteTask(id: number): Promise<void> {
    const res = await fetch(`${API_BASE_URL}/tasks/${id}`, {
      method: 'DELETE',
    });

    const json: ApiResponse<any> = await res.json();
    if (!res.ok || !json.success) {
      throw new ApiError(
        res.status,
        json.error?.message || 'Failed to delete task',
        json.error?.code || 'DELETE_ERROR',
        json.error?.details
      );
    }
  },
};
