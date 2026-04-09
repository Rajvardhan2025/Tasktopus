import { appConfig } from '@/lib/config';

interface ApiResponse<T> {
  data: T;
  status: number;
  headers: Headers;
}

interface ApiErrorResponse {
  data: unknown;
  status: number;
  headers: Headers;
}

class ApiError extends Error {
  response: ApiErrorResponse;

  constructor(message: string, response: ApiErrorResponse) {
    super(message);
    this.name = 'ApiError';
    this.response = response;
  }
}

function joinUrl(base: string, path: string): string {
  if (path.startsWith('http://') || path.startsWith('https://')) {
    return path;
  }

  if (path.startsWith('/')) {
    return `${base}${path}`;
  }

  return `${base}/${path}`;
}

async function parseResponseBody(response: Response): Promise<unknown> {
  const contentType = response.headers.get('content-type') || '';

  if (contentType.includes('application/json')) {
    return response.json();
  }

  const text = await response.text();
  return text;
}

async function request<T>(path: string, init: RequestInit = {}): Promise<ApiResponse<T>> {
  const url = joinUrl(appConfig.apiBaseUrl, path);
  const headers = new Headers(init.headers);

  if (init.body && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }

  const response = await fetch(url, {
    ...init,
    headers,
  });

  const body = await parseResponseBody(response);

  if (!response.ok) {
    const message =
      (body as { error?: { message?: string }; message?: string })?.error?.message ||
      (body as { message?: string })?.message ||
      `Request failed with status ${response.status}`;

    throw new ApiError(message, {
      data: body,
      status: response.status,
      headers: response.headers,
    });
  }

  return {
    data: body as T,
    status: response.status,
    headers: response.headers,
  };
}

const api = {
  get: <T>(path: string) => request<T>(path),
  post: <T>(path: string, data?: unknown) =>
    request<T>(path, {
      method: 'POST',
      body: data === undefined ? undefined : JSON.stringify(data),
    }),
  patch: <T>(path: string, data?: unknown) =>
    request<T>(path, {
      method: 'PATCH',
      body: data === undefined ? undefined : JSON.stringify(data),
    }),
  delete: <T>(path: string) =>
    request<T>(path, {
      method: 'DELETE',
    }),
};

// Types
export interface Project {
  id: string;
  name: string;
  key: string;
  description: string;
  workflow_id: string;
  custom_fields: CustomField[];
  members: string[];
  created_at: string;
  updated_at: string;
}

export interface CustomField {
  name: string;
  type: string;
  options?: string[];
}

export interface Issue {
  id: string;
  issue_key: string;
  project_id: string;
  type: 'epic' | 'story' | 'task' | 'bug' | 'subtask';
  title: string;
  description: string;
  status: string;
  priority: 'lowest' | 'low' | 'medium' | 'high' | 'highest';
  assignee_id?: string;
  reporter_id: string;
  sprint_id?: string;
  parent_id?: string;
  labels: string[];
  story_points?: number;
  custom_fields?: Record<string, any>;
  watchers: string[];
  version: number;
  created_at: string;
  updated_at: string;
}

export interface Sprint {
  id: string;
  project_id: string;
  name: string;
  goal: string;
  start_date?: string;
  end_date?: string;
  status: 'future' | 'active' | 'closed';
  completed_points: number;
  total_points: number;
  velocity: number;
  is_default: boolean;
  issue_count: number;
  complete_issue_count: number;
  created_at: string;
  updated_at: string;
}

export interface Comment {
  id: string;
  issue_id: string;
  user_id: string;
  content: string;
  mentions: string[];
  parent_id?: string;
  created_at: string;
  updated_at: string;
}

export interface Activity {
  id: string;
  project_id: string;
  issue_id?: string;
  user_id: string;
  action: string;
  changes?: Record<string, any>;
  timestamp: string;
}

export interface CursorResponse<T> {
  items: T[];
  next_cursor: string;
}

export interface Notification {
  id: string;
  user_id: string;
  type: 'assigned' | 'mentioned' | 'watched' | 'commented';
  issue_id: string;
  actor_id: string;
  message: string;
  read: boolean;
  timestamp: string;
}

export interface User {
  id: string;
  email: string;
  display_name: string;
  avatar_url?: string;
  created_at: string;
  updated_at: string;
}

// API Functions
export const projectsApi = {
  list: () => api.get<{ data: Project[] }>('/projects'),
  get: (id: string) => api.get<{ data: Project }>(`/projects/${id}`),
  create: (data: Partial<Project>) => api.post<{ data: Project }>('/projects', data),
  update: (id: string, data: Partial<Project>) => api.patch<{ data: Project }>(`/projects/${id}`, data),
  delete: (id: string) => api.delete(`/projects/${id}`),
  members: (projectId: string) => api.get<{ data: User[] }>(`/projects/${projectId}/members`),
  addMember: (projectId: string, userId: string) =>
    api.post(`/projects/${projectId}/members`, { user_id: userId }),
  removeMember: (projectId: string, userId: string) =>
    api.delete(`/projects/${projectId}/members/${userId}`),
};

export const issuesApi = {
  list: (projectId: string) => api.get<{ data: { issues: Issue[] } }>(`/projects/${projectId}/board`),
  get: (id: string) => api.get<{ data: Issue }>(`/issues/${id}`),
  create: (projectId: string, data: Partial<Issue>) =>
    api.post<{ data: Issue }>(`/projects/${projectId}/issues`, data),
  update: (id: string, data: Partial<Issue> & { version: number }) =>
    api.patch<{ data: Issue }>(`/issues/${id}`, data),
  transition: (id: string, toStatus: string, version: number) =>
    api.post<{ data: Issue }>(`/issues/${id}/transitions`, { to_status: toStatus, version }),
  watch: (id: string) => api.post(`/issues/${id}/watch`),
  unwatch: (id: string) => api.delete(`/issues/${id}/watch`),
};

export const sprintsApi = {
  list: (projectId: string) => api.get<{ data: Sprint[] }>(`/projects/${projectId}/sprints`),
  get: (id: string) => api.get<{ data: Sprint }>(`/sprints/${id}`),
  create: (data: Partial<Sprint>) => api.post<{ data: Sprint }>('/sprints', data),
  update: (id: string, data: Partial<Sprint>) => api.patch<{ data: Sprint }>(`/sprints/${id}`, data),
  start: (id: string) => api.post<{ data: Sprint }>(`/sprints/${id}/start`),
  complete: (id: string, carryOverIssues: string[]) =>
    api.post<{ data: Sprint }>(`/sprints/${id}/complete`, { carry_over_issues: carryOverIssues }),
};

export const commentsApi = {
  list: (issueId: string) => api.get<{ data: Comment[] }>(`/issues/${issueId}/comments`),
  create: (issueId: string, content: string, parentId?: string) =>
    api.post<{ data: Comment }>(`/issues/${issueId}/comments`, { content, parent_id: parentId }),
  update: (id: string, content: string) =>
    api.patch<{ data: Comment }>(`/comments/${id}`, { content }),
  delete: (id: string) => api.delete(`/comments/${id}`),
};

export const activityApi = {
  project: (projectId: string, limit = 50, cursor?: string, action?: string, userId?: string) => {
    const params = new URLSearchParams({ limit: String(limit) });
    if (cursor) params.set('cursor', cursor);
    if (action) params.set('action', action);
    if (userId) params.set('user_id', userId);
    return api.get<{ data: CursorResponse<Activity> }>(`/projects/${projectId}/activity?${params}`);
  },
  issue: (issueId: string) =>
    api.get<{ data: Activity[] }>(`/issues/${issueId}/activity`),
};

export const searchApi = {
  search: (query: string, filters?: Record<string, string>) => {
    const params = new URLSearchParams({ q: query, ...filters });
    return api.get<{ data: CursorResponse<Issue> & { query: Record<string, unknown> } }>(`/search?${params}`);
  },
};

export const workflowApi = {
  get: (projectId: string) => api.get(`/projects/${projectId}/workflow`),
  createDefault: (projectId: string) => api.post(`/projects/${projectId}/workflow`),
};

export const notificationsApi = {
  list: (limit = 50) => api.get<{ data: Notification[] }>(`/notifications?limit=${limit}`),
  markRead: (id: string) => api.post(`/notifications/${id}/read`),
  markAllRead: () => api.post('/notifications/read-all'),
};

export const usersApi = {
  list: () => api.get<{ data: User[] }>('/users'),
  create: (data: Pick<User, 'email' | 'display_name' | 'avatar_url'>) => api.post<{ data: User }>('/users', data),
};

export default api;
