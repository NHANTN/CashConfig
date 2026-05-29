import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { User } from '../api/auth';
import * as authApi from '../api/auth';

interface AuthState {
  token: string | null;
  user: User | null;
  permissions: string[];
  loading: boolean;
  initialized: boolean;
  login: (username: string, password: string) => Promise<void>;
  loginLDAP: (username: string, password: string) => Promise<void>;
  loginSSO: (code: string, state: string) => Promise<void>;
  logout: () => void;
  fetchPermissions: () => Promise<void>;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      user: null,
      permissions: [],
      loading: false,
      initialized: false,

      login: async (username, password) => {
        set({ loading: true });
        try {
          const data = await authApi.login(username, password);
          set({ token: data.token, user: data.user, loading: false });
        } catch {
          set({ loading: false });
          throw new Error('登录失败');
        }
      },

      loginLDAP: async (username, password) => {
        set({ loading: true });
        try {
          const data = await authApi.loginLDAP(username, password);
          set({ token: data.token, user: data.user, loading: false });
        } catch {
          set({ loading: false });
          throw new Error('LDAP 登录失败');
        }
      },

      loginSSO: async (code, state) => {
        set({ loading: true });
        try {
          const data = await authApi.handleSSOCallback(code, state);
          set({ token: data.token, user: data.user, loading: false });
        } catch {
          set({ loading: false });
          throw new Error('SSO 登录失败');
        }
      },

      logout: () => {
        set({ token: null, user: null, permissions: [], initialized: true });
      },

      fetchPermissions: async () => {
        try {
          const perms = await authApi.getPermissions();
          set({ permissions: perms });
        } catch {
          set({ permissions: [] });
        }
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({ token: state.token, user: state.user }),
      onRehydrateStorage: () => (state) => {
        state!.initialized = true;
      },
    }
  )
);