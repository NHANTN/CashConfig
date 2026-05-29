import { create } from 'zustand';
import type { Store, StoreListParams } from '../api/store';
import * as api from '../api/store';

interface State {
  list: Store[];
  loading: boolean;
  params: StoreListParams;
  page: number;
  pageSize: number;
  fetch: (params?: StoreListParams) => Promise<void>;
  remove: (id: number) => Promise<void>;
  setParams: (params: StoreListParams) => void;
  setPage: (page: number, pageSize: number) => void;
}

export const useStoreStore = create<State>((set, get) => ({
  list: [],
  loading: false,
  params: {},
  page: 1,
  pageSize: 20,
  setParams: (params) => set({ params }),
  setPage: (page, pageSize) => set({ page, pageSize }),
  fetch: async (params) => {
    set({ loading: true });
    try { set({ list: await api.listStores(params || get().params), loading: false }); }
    catch { set({ loading: false }); }
  },
  remove: async (id) => {
    await api.deleteStore(id);
    const { list, page } = get();
    const newList = list.filter((m) => m.id !== id);
    if (newList.length === 0 && page > 1) {
      set({ list: newList, page: page - 1 });
    } else {
      set({ list: newList });
    }
  },
}));