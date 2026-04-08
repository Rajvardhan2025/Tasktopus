import axios from 'axios';

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

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
  start_date: string;
  end_date: string;
  status: 'planned' | 'active' | 'completed';
  completed_points: number;
  total_points: number;
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

// API Functions
export const projectsApi = {
  list: () => api.get<{ data: Project[] }>('/projects'),
  get: (id: string) => api.get<{ data: Project }>(`/projects/${id}`),
  create: (data: Partial<Project>) => api.post<{ data: Project }>('/projects', data),
  update: (id: string, data: Partial<Project>) => api.patch<{ data: Project }>(`/projects/${id}`, data),
  delete: (id: string) => api.delete(`/projects/${id}`),
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
  project: (projectId: string, limit = 50, skip = 0) => 
    api.get<{ data: Activity[] }>(`/projects/${projectId}/activity?limit=${limit}&skip=${skip}`),
  issue: (issueId: string) => 
    api.get<{ data: Activity[] }>(`/issues/${issueId}/activity`),
};

export const searchApi = {
  search: (query: string, filters?: Record<string, string>) => {
    const params = new URLSearchParams({ q: query, ...filters });
    return api.get<{ data: { issues: Issue[] } }>(`/search?${params}`);
  },
};

export const workflowApi = {
  get: (projectId: string) => api.get(`/projects/${projectId}/workflow`),
  createDefault: (projectId: string) => api.post(`/projects/${projectId}/workflow`),
};

export default api;
