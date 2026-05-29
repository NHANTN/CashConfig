import { create } from 'zustand';
import type { Var, VarListParams } from '../api/var';
import * as api from '../api/var';

interface State {
  list: Var[];
  loading: boolean;
  params: VarListParams;
  page: number;
  pageSize: number;
  fetch: (params?: VarListParams) => Promise<void>;
  remove: (id: number) => Promise<void>;
  setParams: (params: VarListParams) => void;
  setPage: (page: number, pageSize: number) => void;
}

export const useVarStore = create<State>((set, get) => ({
  list: [],
  loading: false,
  params: {},
  page: 1,
  pageSize: 20,
  setParams: (params) => set({ params }),
  setPage: (page, pageSize) => set({ page, pageSize }),
  fetch: async (params) => {
    set({ loading: true });
    try { set({ list: await api.listVars(params || get().params), loading: false }); }
    catch { set({ loading: false }); }
  },
  remove: async (id) => {
    await api.deleteVar(id);
    const { list, page } = get();
    const newList = list.filter((m) => m.id !== id);
    if (newList.length === 0 && page > 1) {
      set({ list: newList, page: page - 1 });
    } else {
      set({ list: newList });
    }
  },
}));