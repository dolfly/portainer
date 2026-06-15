import { FormikErrors } from 'formik';
import { boolean, mixed, object, SchemaOf, string } from 'yup';
import { useState } from 'react';

import { GitAuthModel } from '@/react/portainer/gitops/types';
import { AuthTypeOption } from '@/react/portainer/account/git-credentials/types';

import { SwitchField } from '@@/form-components/SwitchField';
import { TextTip } from '@@/Tip/TextTip';

import { isBE } from '../../feature-flags/feature-flags.service';

import { CredentialsSection } from './CredentialsSection';

interface Props {
  value: GitAuthModel;
  onChange: (value: Partial<GitAuthModel>) => void;
  isAuthExplanationVisible?: boolean;
  errors?: FormikErrors<GitAuthModel>;
}

export function AuthFieldset({
  value: initialValue,
  onChange,
  isAuthExplanationVisible,
  errors,
}: Props) {
  const [value, setValue] = useState(initialValue); // TODO: remove this state when form is not inside angularjs

  return (
    <>
      <div className="form-group">
        <div className="col-sm-12">
          <SwitchField
            label="Authentication"
            labelClass="col-sm-3 col-lg-2"
            name="authentication"
            checked={value.RepositoryAuthentication || false}
            onChange={(value) =>
              handleChange({ RepositoryAuthentication: value })
            }
            data-cy="component-gitAuthToggle"
          />
        </div>
      </div>

      {value.RepositoryAuthentication && (
        <>
          {isAuthExplanationVisible && (
            <TextTip color="orange" className="mb-2">
              Enabling authentication will store the credentials and it is
              advisable to use a git service account
            </TextTip>
          )}

          <CredentialsSection
            value={value}
            onChange={handleChange}
            errors={errors}
          />
        </>
      )}
    </>
  );

  function handleChange(partialValue: Partial<GitAuthModel>) {
    onChange(partialValue);
    setValue((value) => ({ ...value, ...partialValue }));
  }
}

export function gitAuthValidation(
  isAuthEdit: boolean,
  isCreatedFromCustomTemplate: boolean
): SchemaOf<GitAuthModel> {
  return object({
    RepositoryAuthentication: boolean().default(false),
    RepositoryUsername: string()
      .when(['RepositoryAuthentication', 'SourceId'], {
        is: (auth: boolean, sourceId?: number) => auth && !sourceId,
        then: string().required('Username is required'),
      })
      .default(''),
    RepositoryPassword: string()
      .when(['RepositoryAuthentication', 'SourceId'], {
        is: (auth: boolean, sourceId?: number) =>
          auth && !sourceId && !isAuthEdit && !isCreatedFromCustomTemplate,
        then: string().required('Personal Access Token is required'),
      })
      .default(''),
    RepositoryAuthorizationType: mixed()
      .oneOf(Object.values(AuthTypeOption))
      .when(['RepositoryAuthentication', 'SourceId'], {
        is: (auth: boolean, sourceId?: number) =>
          isBE &&
          auth &&
          !sourceId &&
          !isAuthEdit &&
          !isCreatedFromCustomTemplate,
        then: mixed().required('Authorization type is required'),
      })
      .default(AuthTypeOption.Basic),
  });
}
