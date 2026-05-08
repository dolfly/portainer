import {
  asEnum,
  useTableStateFromUrl,
} from '@@/datatables/useTableStateFromUrl';

import { WorkflowStatus, WorkflowType, DeploymentPlatform } from './types';

const DEFAULT_SORT = 'name' as const;

const WORKFLOW_STATUSES = new Set<WorkflowStatus>([
  'healthy',
  'error',
  'syncing',
  'paused',
  'unknown',
]);
const WORKFLOW_TYPES = new Set<WorkflowType>(['stack', 'edgeStack']);
const DEPLOYMENT_PLATFORMS = new Set<DeploymentPlatform>([
  'dockerStandalone',
  'dockerSwarm',
  'kubernetes',
]);

export function useListState() {
  return useTableStateFromUrl({
    localStorageKey: 'workflows',
    defaultSort: DEFAULT_SORT,
    parseExtra: (params) => ({
      status: asEnum(params.status, WORKFLOW_STATUSES),
      type: asEnum(params.type, WORKFLOW_TYPES),
      platform: asEnum(params.platform, DEPLOYMENT_PLATFORMS),
    }),
    buildExtra: (urlState, setUrlState) => {
      return {
        status: urlState.status,
        type: urlState.type,
        platform: urlState.platform,
        setStatus: (v: WorkflowStatus | null) =>
          setUrlState({ status: v, page: 0 }),

        setGroupFilter: (group: string | null, filter: string | null) => {
          if (group === 'status') {
            setUrlState({
              groupBy: group,
              groupFilter: filter,
              status: filter as WorkflowStatus | null,
              type: null,
              platform: null,
              page: 0,
            });
          } else if (group === 'type') {
            setUrlState({
              groupBy: group,
              groupFilter: filter,
              type: filter as WorkflowType | null,
              status: null,
              platform: null,
              page: 0,
            });
          } else if (group === 'platform') {
            setUrlState({
              groupBy: group,
              groupFilter: filter,
              platform: filter as DeploymentPlatform | null,
              type: null,
              status: null,
              page: 0,
            });
          }
        },
        setSortBy: (id: string, desc: boolean) =>
          setUrlState({
            sort: id ?? DEFAULT_SORT,
            order: desc ? 'desc' : 'asc',
            groupBy: null,
            groupFilter: null,
            status: null,
            type: null,
            platform: null,
            page: 0,
          }),
      };
    },
  });
}
