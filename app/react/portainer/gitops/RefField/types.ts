import { GitAuthModel } from '../types';

export interface RefFieldModel extends GitAuthModel {
  RepositoryURL: string;
  TLSSkipVerify?: boolean;
  SourceId?: number;
}
