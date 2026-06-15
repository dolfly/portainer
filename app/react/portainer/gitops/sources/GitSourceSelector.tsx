import { FormControl } from '@@/form-components/FormControl';
import { Select } from '@@/form-components/ReactSelect';

import { useSources } from './queries/useSources';
import { Source } from './types';

export function GitSourceSelector({
  value,
  onChange,
  error,
}: {
  value?: Source['id'];
  onChange(source?: Source | null): void;
  error?: string;
}) {
  const sourcesQuery = useSources({ type: 'git' });
  const sources = sourcesQuery.data?.data ?? [];

  return (
    <div className="form-group">
      <div className="col-sm-12">
        <FormControl label="Source" inputId="source-selector" errors={error}>
          <Select
            placeholder="Select a source"
            value={sources.find((s) => s.id === value) ?? null}
            options={sources}
            getOptionLabel={(s) => s.name}
            getOptionValue={(s) => String(s.id)}
            onChange={onChange}
            isClearable
            isLoading={sourcesQuery.isLoading}
            noOptionsMessage={() => 'No git sources available'}
            inputId="source-selector"
            data-cy="source-selector"
          />
        </FormControl>
      </div>
    </div>
  );
}
