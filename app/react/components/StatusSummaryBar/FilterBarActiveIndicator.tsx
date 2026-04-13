interface Props {
  label: string;
  onClear: () => void;
}

export function FilterBarActiveIndicator({ label, onClear }: Props) {
  return (
    <div
      className="flex items-center gap-4 px-5 bg-gray-2 th-dark:bg-gray-iron-10 th-highcontrast:bg-transparent whitespace-nowrap"
      data-cy="active-filter-indicator"
    >
      <span className="text-sm text-[var(--text-muted-color)]">
        Showing:{' '}
        <span className="font-semibold text-[var(--text-summary-color)]">
          {label}
        </span>
      </span>
      <button
        type="button"
        className="border-0 bg-transparent cursor-pointer p-0 text-[21px] font-bold leading-none opacity-[var(--button-opacity)] hover:opacity-[var(--button-opacity-hover)] text-[var(--button-close-color)]"
        onClick={onClear}
        aria-label="Clear filter"
        data-cy="clear-filter-button"
      >
        &times;
      </button>
    </div>
  );
}
