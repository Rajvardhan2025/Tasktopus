import api from './api';

const TOKEN_KEY = 'auth_token';
const USER_KEY = 'auth_user';

export interface User {
  id: string;
  email: string;
  display_name: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface RegisterData {
  email: string;
  password: string;
  display_name: string;
}

export interface LoginData {
  email: string;
  password: string;
}

export interface ChangePasswordData {
  current_password: string;
  new_password: string;
}

class AuthService {
  // Get stored token
  getToken(): string | null {
    return localStorage.getItem(TOKEN_KEY);
  }

  // Get stored user
  getUser(): User | null {
    const userStr = localStorage.getItem(USER_KEY);
    if (!userStr) return null;
    try {
      return JSON.parse(userStr);
    } catch {
      return null;
    }
  }

  // Check if user is authenticated
  isAuthenticated(): boolean {
    return !!this.getToken();
  }

  // Store auth data
  private setAuthData(token: string, user: User): void {
    localStorage.setItem(TOKEN_KEY, token);
    localStorage.setItem(USER_KEY, JSON.stringify(user));
  }

  // Clear auth data
  private clearAuthData(): void {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
  }

  // Register new user
  async register(data: RegisterData): Promise<AuthResponse> {
    const response = await api.post<{ data: AuthResponse }>('/auth/register', data);
    const authData = response.data.data;
    this.setAuthData(authData.token, authData.user);
    return authData;
  }

  // Login user
  async login(data: LoginData): Promise<AuthResponse> {
    const response = await api.post<{ data: AuthResponse }>('/auth/login', data);
    const authData = response.data.data;
    this.setAuthData(authData.token, authData.user);
    return authData;
  }

  // Logout user
  logout(): void {
    this.clearAuthData();
    window.location.href = '/';
  }

  // Get current user from API
  async getCurrentUser(): Promise<User> {
    const response = await api.get<{ data: User }>('/auth/me');
    const user = response.data.data;
    // Update stored user
    localStorage.setItem(USER_KEY, JSON.stringify(user));
    return user;
  }

  // Change password
  async changePassword(data: ChangePasswordData): Promise<void> {
    await api.post('/auth/change-password', data);
  }

  // Add token to request headers
  getAuthHeader(): Record<string, string> {
    const token = this.getToken();
    if (!token) return {};
    return {
      Authorization: `Bearer ${token}`,
    };
  }
}

export const authService = new AuthService();
