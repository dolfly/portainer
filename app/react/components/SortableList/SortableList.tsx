import '@@/datatables/datatable.css';

import { ReactNode } from 'react';

import { AutomationTestingProps } from '@/types';

import { DropdownOption } from '../DropdownMenu/DropdownMenu';

import { SortOption, SortableListHeader } from './SortableListHeader';
import { SortableGroup } from './SortableListGroup';
import { SortableListBody } from './SortableListBody';
import { SortableListCard } from './SortableListCard';
import { SortableListPager } from './SortableListPager';
import { SortableListState } from './sortable-list.store';

export type { SortableGroup };
export type { SortableListState };
export type { SortOption };

interface Props<T> extends AutomationTestingProps {
  tableState: SortableListState;
  sortOptions: SortOption[];
  groupOptions?: Record<string, DropdownOption[]>;
  groups: SortableGroup<T>[];
  totalCount: number;
  renderItem: (item: T, index: number) => ReactNode;
  renderColumnHeaders?: (groupKey: string, items: T[]) => ReactNode;
  getItemKey?: (item: T, index: number) => string | number;
  showGroupHeaders?: boolean;
  emptyMessage?: string;
  searchPlaceholder?: string;
  headerButtons?: ReactNode;
  actionButton?: ReactNode;
  isLoading?: boolean;
}

export function SortableList<T>({
  tableState,
  sortOptions,
  groupOptions,
  groups,
  totalCount,
  renderItem,
  renderColumnHeaders,
  getItemKey,
  showGroupHeaders = true,
  emptyMessage = 'No items found',
  searchPlaceholder,
  headerButtons,
  actionButton,
  isLoading = false,
  'data-cy': dataCy,
}: Props<T>) {
  const rawSortKey = tableState.sortBy?.id ?? '';
  const activeSortKey =
    sortOptions
      .filter((o) => !o.grouped)
      .find((opt) => opt.key.toLowerCase() === rawSortKey.toLowerCase())?.key ??
    null;
  const rawGroupKey = tableState.groupBy ?? '';
  const activeGroupKey =
    sortOptions
      .filter((o) => o.grouped)
      .find((opt) => opt.key.toLowerCase() === rawGroupKey.toLowerCase())
      ?.key ?? null;
  const activeKey = activeGroupKey ?? activeSortKey ?? '';

  return (
    <SortableListCard>
      <SortableListHeader
        activeKey={activeKey}
        sortDesc={tableState.sortBy?.desc ?? false}
        onSortChange={(key) => {
          tableState.setSortBy(
            key,
            computeSortDesc(key, activeKey, tableState.sortBy?.desc ?? false)
          );
        }}
        searchTerm={tableState.search}
        onSearchChange={tableState.setSearch}
        groupFilter={tableState.groupFilter}
        onGroupFilterChange={tableState.setGroupFilter}
        groupOptions={groupOptions}
        sortOptions={sortOptions}
        searchPlaceholder={searchPlaceholder}
        headerButtons={headerButtons}
        actionButton={actionButton}
        data-cy={`${dataCy}-header`}
      />

      <div className="overflow-y-auto">
        <SortableListBody
          isLoading={isLoading}
          groups={groups}
          showGroupHeaders={
            showGroupHeaders && (groupOptions?.[activeKey]?.length ?? 0) > 0
          }
          renderItem={renderItem}
          renderColumnHeaders={renderColumnHeaders}
          getItemKey={getItemKey}
          emptyMessage={emptyMessage}
          data-cy={`${dataCy}-body`}
        />
      </div>

      <SortableListPager
        page={tableState.page}
        pageSize={tableState.pageSize}
        totalCount={totalCount}
        onPageChange={tableState.setPage}
        onPageSizeChange={tableState.setPageSize}
      />
    </SortableListCard>
  );
}

export function computeSortDesc(
  key: string,
  activeKey: string,
  currentDesc: boolean
): boolean {
  return activeKey === key ? !currentDesc : false;
}
