import { create } from 'zustand';
import type { Module, ModuleListParams } from '../api/module';
import * as moduleApi from '../api/module';

interface ModuleState {
  list: Module[];
  loading: boolean;
  params: ModuleListParams;
  page: number;
  pageSize: number;
  fetch: (params?: ModuleListParams) => Promise<void>;
  remove: (id: number) => Promise<void>;
  setParams: (params: ModuleListParams) => void;
  setPage: (page: number, pageSize: number) => void;
}

export const useModuleStore = create<ModuleState>((set, get) => ({
  list: [],
  loading: false,
  params: {},
  page: 1,
  pageSize: 20,

  setParams: (params) => set({ params }),

  setPage: (page, pageSize) => set({ page, pageSize }),

  fetch: async (params) => {
    set({ loading: true });
    try {
      const list = await moduleApi.listModules(params || get().params);
      set({ list, loading: false });
    } catch {
      set({ loading: false });
    }
  },

  remove: async (id) => {
    await moduleApi.deleteModule(id);
    const { list, page } = get();
    const newList = list.filter((m) => m.id !== id);
    if (newList.length === 0 && page > 1) {
      set({ list: newList, page: page - 1 });
    } else {
      set({ list: newList });
    }
  },
}));