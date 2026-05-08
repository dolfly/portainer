import { useState } from 'react';

import {
  BasicTableSettings,
  BackendPaginationTableSettings,
  backendPaginationSettings,
  createPersistedStore,
  createTableStore,
  ZustandSetFunc,
} from '@@/datatables/types';
import { useTableState, TableState } from '@@/datatables/useTableState';

interface SortableListSettings
  extends BasicTableSettings,
    BackendPaginationTableSettings {
  groupBy: string | null;
  setGroupBy: (group: string | null) => void;
  groupFilter: string | null;
  setGroupFilter: (group: string | null, filter: string | null) => void;
}

export type SortableListState = TableState<SortableListSettings>;

function sortableListExtras(
  set: ZustandSetFunc<SortableListSettings>
): Partial<SortableListSettings> {
  return {
    ...backendPaginationSettings<SortableListSettings>(set),
    groupBy: null,
    setGroupBy: (group) => set((s) => ({ ...s, groupBy: group })),
    groupFilter: null,
    setGroupFilter: (group: string | null, filter: string | null) =>
      set((s) => ({ ...s, groupBy: group, groupFilter: filter, page: 0 })),
  };
}

export function createSortableListStore(
  storageKey: string,
  defaultSort?: string
) {
  return createPersistedStore<SortableListSettings>(
    storageKey,
    defaultSort,
    sortableListExtras
  );
}

export function createSortableListStoreUnpersisted(defaultSort?: string) {
  return createTableStore<SortableListSettings>(
    defaultSort,
    sortableListExtras
  );
}

export function useSortableListState(
  storageKey: string,
  defaultSort?: string
): SortableListState {
  const [store] = useState(() =>
    createSortableListStore(storageKey, defaultSort)
  );
  return useTableState(store, storageKey);
}
