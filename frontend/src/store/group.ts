import { create } from 'zustand';
import type { Group } from '../api/group';
import * as api from '../api/group';

interface GroupState {
  list: Group[];
  all: Group[];
  loading: boolean;
  page: number;
  pageSize: number;
  fetch: (name?: string) => Promise<void>;
  fetchAll: () => Promise<void>;
  remove: (id: number) => Promise<void>;
  setPage: (page: number, pageSize: number) => void;
}

export const useGroupStore = create<GroupState>((set, get) => ({
  list: [],
  all: [],
  loading: false,
  page: 1,
  pageSize: 20,
  fetch: async (name) => {
    set({ loading: true });
    try { set({ list: await api.listGroups(name), loading: false }); }
    catch { set({ loading: false }); }
  },
  fetchAll: async () => {
    try { set({ all: await api.listAllGroups() }); }
    catch {}
  },
  remove: async (id) => {
    await api.deleteGroup(id);
    const { list, page } = get();
    const newList = list.filter((m) => m.id !== id);
    if (newList.length === 0 && page > 1) {
      set({ list: newList, page: page - 1 });
    } else {
      set({ list: newList });
    }
  },
  setPage: (page, pageSize) => set({ page, pageSize }),
}));