import { IFormController } from 'angular';
import { FormikErrors } from 'formik';

import { GitAuthModel } from '@/react/portainer/gitops/types';
import { gitAuthValidation } from '@/react/portainer/gitops/AuthFieldset';

import { validateForm } from '@@/form-components/validate-form';

export default class GitFormAuthFieldsetController {
  errors?: FormikErrors<GitAuthModel> = {};

  $async: <T>(fn: () => Promise<T>) => Promise<T>;

  gitFormAuthFieldset?: IFormController;

  value?: GitAuthModel;

  isAuthEdit: boolean;

  onChange?: (value: GitAuthModel) => void;

  /* @ngInject */
  constructor($async: <T>(fn: () => Promise<T>) => Promise<T>) {
    this.$async = $async;

    this.isAuthEdit = false;
    this.handleChange = this.handleChange.bind(this);
    this.runGitValidation = this.runGitValidation.bind(this);
  }

  async handleChange(newValues: Partial<GitAuthModel>) {
    // this should never happen, but just in case
    if (!this.value) {
      throw new Error('GitFormController: value is required');
    }

    const value = {
      ...this.value,
      ...newValues,
    };
    this.onChange?.(value);
    await this.runGitValidation(value, this.isAuthEdit);
  }

  async runGitValidation(value: GitAuthModel, isAuthEdit: boolean) {
    return this.$async(async () => {
      this.errors = {};
      this.gitFormAuthFieldset?.$setValidity(
        'gitFormAuth',
        true,
        this.gitFormAuthFieldset
      );

      this.errors = await validateForm<GitAuthModel>(
        () => gitAuthValidation(isAuthEdit, false),
        value
      );
      if (this.errors && Object.keys(this.errors).length > 0) {
        this.gitFormAuthFieldset?.$setValidity(
          'gitFormAuth',
          false,
          this.gitFormAuthFieldset
        );
      }
    });
  }

  async $onInit() {
    // this should never happen, but just in case
    if (!this.value) {
      throw new Error('GitFormController: value is required');
    }

    await this.runGitValidation(this.value, this.isAuthEdit);
  }
}
