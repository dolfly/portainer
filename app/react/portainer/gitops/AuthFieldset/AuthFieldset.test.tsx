import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';

import { AuthTypeOption } from '@/react/portainer/account/git-credentials/types';
import { GitAuthModel } from '@/react/portainer/gitops/types';

import { AuthFieldset, gitAuthValidation } from './AuthFieldset';

vi.mock('../../feature-flags/feature-flags.service', () => ({
  isBE: true,
  isLimitedToBE: () => false,
}));

vi.mock('@/react/hooks/useDebounce', () => ({
  useDebounce: (value: unknown, callback: (value: unknown) => void) => [
    value,
    callback,
  ],
}));

const defaultGitAuthModel: GitAuthModel = {
  RepositoryAuthentication: false,
  RepositoryUsername: '',
  RepositoryPassword: '',
  RepositoryAuthorizationType: AuthTypeOption.Basic,
};

function renderAuthFieldset({
  value = defaultGitAuthModel,
  onChange = vi.fn(),
  isAuthExplanationVisible = false,
  errors = {},
}: {
  value?: GitAuthModel;
  onChange?: (value: Partial<GitAuthModel>) => void;
  isAuthExplanationVisible?: boolean;
  errors?: Record<string, string>;
} = {}) {
  return render(
    <AuthFieldset
      value={value}
      onChange={onChange}
      isAuthExplanationVisible={isAuthExplanationVisible}
      errors={errors}
    />
  );
}

describe('AuthFieldset', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('key component rendering', () => {
    it('should render authentication toggle', () => {
      renderAuthFieldset();

      expect(screen.getByTestId('component-gitAuthToggle')).toBeInTheDocument();
    });

    it('should render username input when authentication is enabled', () => {
      renderAuthFieldset({
        value: { ...defaultGitAuthModel, RepositoryAuthentication: true },
      });

      expect(
        screen.getByTestId('component-gitUsernameInput')
      ).toBeInTheDocument();
    });

    it('should render password input when authentication is enabled', () => {
      renderAuthFieldset({
        value: { ...defaultGitAuthModel, RepositoryAuthentication: true },
      });

      expect(
        screen.getByTestId('component-gitPasswordInput')
      ).toBeInTheDocument();
    });

    it('should not render interactive fields when authentication is disabled', () => {
      renderAuthFieldset({
        value: { ...defaultGitAuthModel, RepositoryAuthentication: false },
      });

      expect(
        screen.queryByTestId('component-gitUsernameInput')
      ).not.toBeInTheDocument();
      expect(
        screen.queryByTestId('component-gitPasswordInput')
      ).not.toBeInTheDocument();
    });
  });

  describe('props handling', () => {
    it('should handle onChange prop', () => {
      const onChange = vi.fn();
      renderAuthFieldset({ onChange });

      expect(onChange).toBeDefined();
    });

    it('should handle errors prop for username', () => {
      const errors = { RepositoryUsername: 'Username is required' };
      renderAuthFieldset({
        value: { ...defaultGitAuthModel, RepositoryAuthentication: true },
        errors,
      });

      expect(screen.getByText('Username is required')).toBeInTheDocument();
    });

    it('should handle errors prop for password', () => {
      const errors = { RepositoryPassword: 'Password is required' };
      renderAuthFieldset({
        value: { ...defaultGitAuthModel, RepositoryAuthentication: true },
        errors,
      });

      expect(screen.getByText('Password is required')).toBeInTheDocument();
    });

    it('should handle multiple errors', () => {
      const errors = {
        RepositoryUsername: 'Username is required',
        RepositoryPassword: 'Password is required',
      };
      renderAuthFieldset({
        value: { ...defaultGitAuthModel, RepositoryAuthentication: true },
        errors,
      });

      expect(screen.getByText('Username is required')).toBeInTheDocument();
      expect(screen.getByText('Password is required')).toBeInTheDocument();
    });

    it('should handle empty errors object', () => {
      const errors = {};
      renderAuthFieldset({
        value: { ...defaultGitAuthModel, RepositoryAuthentication: true },
        errors,
      });

      expect(
        screen.getByTestId('component-gitUsernameInput')
      ).toBeInTheDocument();
      expect(
        screen.getByTestId('component-gitPasswordInput')
      ).toBeInTheDocument();
    });

    it('should handle isAuthExplanationVisible prop when true', () => {
      renderAuthFieldset({
        value: { ...defaultGitAuthModel, RepositoryAuthentication: true },
        isAuthExplanationVisible: true,
      });

      expect(
        screen.getByText(
          'Enabling authentication will store the credentials and it is advisable to use a git service account'
        )
      ).toBeInTheDocument();
    });

    it('should handle isAuthExplanationVisible prop when false', () => {
      renderAuthFieldset({
        value: { ...defaultGitAuthModel, RepositoryAuthentication: true },
        isAuthExplanationVisible: false,
      });

      expect(
        screen.queryByText(
          'Enabling authentication will store the credentials and it is advisable to use a git service account'
        )
      ).not.toBeInTheDocument();
    });

    it('should handle value prop with all fields populated', () => {
      const value: GitAuthModel = {
        RepositoryAuthentication: true,
        RepositoryUsername: 'testuser',
        RepositoryPassword: 'testpass',
        RepositoryAuthorizationType: AuthTypeOption.Token,
      };

      renderAuthFieldset({ value });

      expect(
        screen.getByTestId('component-gitUsernameInput')
      ).toBeInTheDocument();
      expect(
        screen.getByTestId('component-gitPasswordInput')
      ).toBeInTheDocument();
    });
  });
});

describe('gitAuthValidation', () => {
  describe('default values', () => {
    it('should provide correct default values', async () => {
      const schema = gitAuthValidation(false, false);
      const result = await schema.validate({});

      expect(result).toEqual({
        RepositoryAuthentication: false,
        RepositoryUsername: '',
        RepositoryPassword: '',
        RepositoryAuthorizationType: AuthTypeOption.Basic,
      });
    });
  });

  describe('authentication disabled', () => {
    it('should allow empty values when authentication is disabled', async () => {
      const schema = gitAuthValidation(false, false);
      const data = {
        RepositoryAuthentication: false,
        RepositoryUsername: '',
        RepositoryPassword: '',
        RepositoryAuthorizationType: AuthTypeOption.Basic,
      };

      const result = await schema.validate(data);
      expect(result.RepositoryAuthentication).toBe(false);
    });
  });

  describe('authentication enabled', () => {
    it('should require username when authentication is enabled', async () => {
      const schema = gitAuthValidation(false, false);
      const data = {
        RepositoryAuthentication: true,
        RepositoryUsername: '',
        RepositoryPassword: 'password',
        RepositoryAuthorizationType: AuthTypeOption.Basic,
      };

      await expect(schema.validate(data)).rejects.toThrow(
        'Username is required'
      );
    });

    it('should require password when authentication is enabled, not auth edit, and not from custom template', async () => {
      const schema = gitAuthValidation(false, false);
      const data = {
        RepositoryAuthentication: true,
        RepositoryUsername: 'username',
        RepositoryPassword: '',
        RepositoryAuthorizationType: AuthTypeOption.Basic,
      };

      await expect(schema.validate(data)).rejects.toThrow(
        'Personal Access Token is required'
      );
    });

    it('should set default authorization type when authentication is enabled', async () => {
      const schema = gitAuthValidation(false, false);
      const data = {
        RepositoryAuthentication: true,
        RepositoryUsername: 'username',
        RepositoryPassword: 'password',
        RepositoryAuthorizationType: undefined,
      };

      const result = await schema.validate(data);
      expect(result.RepositoryAuthorizationType).toBe(AuthTypeOption.Basic);
    });

    it('should accept valid authorization types', async () => {
      const schema = gitAuthValidation(false, false);
      const data = {
        RepositoryAuthentication: true,
        RepositoryUsername: 'username',
        RepositoryPassword: 'password',
        RepositoryAuthorizationType: AuthTypeOption.Token,
      };

      const result = await schema.validate(data);
      expect(result.RepositoryAuthorizationType).toBe(AuthTypeOption.Token);
    });

    it('should reject invalid authorization types', async () => {
      const schema = gitAuthValidation(false, false);
      const data = {
        RepositoryAuthentication: true,
        RepositoryUsername: 'username',
        RepositoryPassword: 'password',
        RepositoryAuthorizationType: 999,
      };

      await expect(schema.validate(data)).rejects.toThrow();
    });
  });

  describe('auth edit mode', () => {
    it('should not require password when in auth edit mode', async () => {
      const schema = gitAuthValidation(true, false);
      const data = {
        RepositoryAuthentication: true,
        RepositoryUsername: 'username',
        RepositoryPassword: '',
        RepositoryAuthorizationType: AuthTypeOption.Basic,
      };

      const result = await schema.validate(data);
      expect(result.RepositoryPassword).toBe('');
    });

    it('should not require authorization type when in auth edit mode', async () => {
      const schema = gitAuthValidation(true, false);
      const data = {
        RepositoryAuthentication: true,
        RepositoryUsername: 'username',
        RepositoryPassword: 'password',
        RepositoryAuthorizationType: undefined,
      };

      const result = await schema.validate(data);
      expect(result.RepositoryAuthorizationType).toBe(AuthTypeOption.Basic);
    });
  });

  describe('created from custom template', () => {
    it('should not require password when created from custom template', async () => {
      const schema = gitAuthValidation(false, true);
      const data = {
        RepositoryAuthentication: true,
        RepositoryUsername: 'username',
        RepositoryPassword: '',
        RepositoryAuthorizationType: AuthTypeOption.Basic,
      };

      const result = await schema.validate(data);
      expect(result.RepositoryPassword).toBe('');
    });

    it('should not require authorization type when created from custom template', async () => {
      const schema = gitAuthValidation(false, true);
      const data = {
        RepositoryAuthentication: true,
        RepositoryUsername: 'username',
        RepositoryPassword: 'password',
        RepositoryAuthorizationType: undefined,
      };

      const result = await schema.validate(data);
      expect(result.RepositoryAuthorizationType).toBe(AuthTypeOption.Basic);
    });
  });

  describe('complex scenarios', () => {
    it('should handle complete valid data', async () => {
      const schema = gitAuthValidation(false, false);
      const data = {
        RepositoryAuthentication: true,
        RepositoryUsername: 'testuser',
        RepositoryPassword: 'testpassword',
        RepositoryAuthorizationType: AuthTypeOption.Token,
      };

      const result = await schema.validate(data);
      expect(result).toEqual(data);
    });

    it('should handle auth edit mode', async () => {
      const schema = gitAuthValidation(true, false);
      const data = {
        RepositoryAuthentication: true,
        RepositoryUsername: 'testuser',
        RepositoryPassword: '',
        RepositoryAuthorizationType: AuthTypeOption.Basic,
      };

      const result = await schema.validate(data);
      expect(result.RepositoryPassword).toBe('');
    });

    it('should handle custom template creation', async () => {
      const schema = gitAuthValidation(false, true);
      const data = {
        RepositoryAuthentication: true,
        RepositoryUsername: 'testuser',
        RepositoryPassword: '',
        RepositoryAuthorizationType: AuthTypeOption.Basic,
      };

      const result = await schema.validate(data);
      expect(result.RepositoryPassword).toBe('');
    });
  });
});
