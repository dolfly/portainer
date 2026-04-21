import clsx from 'clsx';

import { DropdownMenu, DropdownOption } from '../DropdownMenu/DropdownMenu';

export interface SortOption<TSortKey extends string = string> {
  key: TSortKey;
  label: string;
  grouped?: boolean;
}

export interface SortByGroupProps<TSortKey extends string> {
  sortBy: TSortKey;
  onSortChange: (key: TSortKey) => void;
  sortOptions: SortOption<TSortKey>[];
  groupFilter: string | null;
  groupOptions?: Record<string, DropdownOption[]>;
  onGroupFilterChange: (value: string | null) => void;
  dataCy?: string;
}

export function SortByGroup<TSortKey extends string>({
  sortBy,
  onSortChange,
  sortOptions,
  groupFilter,
  groupOptions,
  onGroupFilterChange,
  dataCy,
}: SortByGroupProps<TSortKey>) {
  return (
    <>
      <span
        className="text-xs font-semibold tracking-wider text-gray-11 th-highcontrast:text-white th-dark:text-white"
        data-cy="sort-by-label"
      >
        SORT BY:
      </span>
      <div
        className={clsx(
          'flex',
          'bg-gray-4 th-highcontrast:bg-black th-dark:bg-gray-iron-11',
          'gap-1 rounded-md p-1 th-highcontrast:border th-highcontrast:border-solid th-highcontrast:border-white'
        )}
        role="group"
        aria-label="Sort by"
      >
        {sortOptions.map((option, index) => (
          <SortOptionItem
            key={option.key}
            option={option}
            isActive={sortBy === option.key}
            isFirst={index === 0}
            isLast={index === sortOptions.length - 1}
            onSortChange={onSortChange}
            groupFilter={groupFilter}
            groupOptions={groupOptions}
            onGroupFilterChange={onGroupFilterChange}
            dataCy={dataCy}
          />
        ))}
      </div>
    </>
  );
}

const baseBtn =
  'px-4 py-1.5 text-xs align-middle font-medium transition-colors';
const activeBtn =
  'z-10 border-none rounded-md font-medium ' +
  'bg-white text-gray-8 ' +
  'th-dark:bg-gray-iron-10 th-dark:text-white ' +
  'th-highcontrast:bg-white th-highcontrast:text-black';
const inactiveBtn =
  'text-muted border-none rounded-md ' +
  'bg-gray-4 hover:bg-gray-2 ' +
  'th-dark:bg-gray-iron-11 th-dark:text-white th-dark:hover:bg-gray-iron-9 ' +
  'th-highcontrast:bg-black th-highcontrast:text-white th-highcontrast:hover:bg-white th-highcontrast:hover:text-black';

interface SortOptionItemProps<TSortKey extends string> {
  option: SortOption<TSortKey>;
  isActive: boolean;
  isFirst: boolean;
  isLast: boolean;
  onSortChange: (key: TSortKey) => void;
  groupFilter: string | null;
  groupOptions?: Record<string, DropdownOption[]>;
  onGroupFilterChange: (value: string | null) => void;
  dataCy?: string;
}

function SortOptionItem<TSortKey extends string>({
  option,
  isActive,
  isFirst,
  isLast,
  onSortChange,
  groupFilter,
  groupOptions,
  onGroupFilterChange,
  dataCy,
}: SortOptionItemProps<TSortKey>) {
  const className = clsx(
    baseBtn,
    isActive ? activeBtn : inactiveBtn,
    isFirst && 'rounded-l-md',
    isLast && 'rounded-r-md'
  );

  if (option.grouped) {
    return (
      <DropdownMenu
        label={option.label}
        options={groupOptions?.[option.key]}
        selected={groupFilter}
        onSelect={onGroupFilterChange}
        badge={isActive ? groupFilter : undefined}
        className={className}
        data-cy={`${dataCy}-sort-by-${option.key.toLowerCase()}-button`}
        onClick={() => {
          if (!isActive) {
            onSortChange(option.key);
          }
        }}
      />
    );
  }

  return (
    <button
      type="button"
      className={className}
      onClick={() => {
        if (!isActive) {
          onSortChange(option.key);
        }
      }}
      data-cy={`${dataCy}-sort-by-${option.key.toLowerCase()}-button`}
    >
      {option.label}
    </button>
  );
}
