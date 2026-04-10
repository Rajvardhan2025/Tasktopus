import api from './api';
import { toast } from '@/hooks/use-toast';

const TOKEN_KEY = 'auth_token';
const USER_KEY = 'auth_user';
const SERVER_WARMUP_THRESHOLD = 3000; // 3 seconds

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

  // Helper to handle slow requests with user feedback
  private async handleAuthRequest<T>(
    requestFn: () => Promise<T>,
    isFirstRequest: boolean = false
  ): Promise<T> {
    let warmupToastId: string | undefined;
    
    // Set a timer to show warmup message if request takes too long
    const warmupTimer = setTimeout(() => {
      if (isFirstRequest) {
        const result = toast({
          title: '🚀 Server is waking up...',
          description: 'This may take 30-60 seconds on first request. Please wait.',
          duration: 60000, // Show for 60 seconds
        });
        warmupToastId = result.id;
      }
    }, SERVER_WARMUP_THRESHOLD);

    try {
      const result = await requestFn();
      clearTimeout(warmupTimer);
      
      // Dismiss warmup toast if it was shown
      if (warmupToastId) {
        toast({ id: warmupToastId, title: '', description: '', duration: 0 });
      }
      
      return result;
    } catch (error) {
      clearTimeout(warmupTimer);
      
      // Dismiss warmup toast if it was shown
      if (warmupToastId) {
        toast({ id: warmupToastId, title: '', description: '', duration: 0 });
      }
      
      throw error;
    }
  }

  // Check if this might be the first request (server cold start)
  private isLikelyFirstRequest(): boolean {
    // If no token exists, this is likely a first-time user or fresh session
    return !this.getToken();
  }

  // Register new user
  async register(data: RegisterData): Promise<AuthResponse> {
    const isFirstRequest = this.isLikelyFirstRequest();
    
    return this.handleAuthRequest(async () => {
      const response = await api.post<{ data: AuthResponse }>('/auth/register', data);
      const authData = response.data.data;
      this.setAuthData(authData.token, authData.user);
      return authData;
    }, isFirstRequest);
  }

  // Login user
  async login(data: LoginData): Promise<AuthResponse> {
    const isFirstRequest = this.isLikelyFirstRequest();
    
    return this.handleAuthRequest(async () => {
      const response = await api.post<{ data: AuthResponse }>('/auth/login', data);
      const authData = response.data.data;
      this.setAuthData(authData.token, authData.user);
      return authData;
    }, isFirstRequest);
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
