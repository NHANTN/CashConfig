import { create } from 'zustand';
import type { TillList, TillListParams } from '../api/till_list';
import * as api from '../api/till_list';

interface State {
  list: TillList[];
  loading: boolean;
  params: TillListParams;
  page: number;
  pageSize: number;
  fetch: (params?: TillListParams) => Promise<void>;
  remove: (id: number) => Promise<void>;
  setParams: (params: TillListParams) => void;
  setPage: (page: number, pageSize: number) => void;
}

export const useTillListStore = create<State>((set, get) => ({
  list: [],
  loading: false,
  params: {},
  page: 1,
  pageSize: 20,
  setParams: (params) => set({ params }),
  setPage: (page, pageSize) => set({ page, pageSize }),
  fetch: async (params) => {
    set({ loading: true });
    try { set({ list: await api.listTillLists(params || get().params), loading: false }); }
    catch { set({ loading: false }); }
  },
  remove: async (id) => {
    await api.deleteTillList(id);
    const { list, page } = get();
    const newList = list.filter((m) => m.id !== id);
    if (newList.length === 0 && page > 1) {
      set({ list: newList, page: page - 1 });
    } else {
      set({ list: newList });
    }
  },
}));