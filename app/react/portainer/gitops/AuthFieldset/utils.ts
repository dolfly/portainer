import { GitAuthenticationResponse, GitAuthModel } from '../types';

export function parseAuthResponse(
  auth?: GitAuthenticationResponse
): GitAuthModel {
  if (!auth) {
    return {
      RepositoryAuthentication: false,
      RepositoryPassword: '',
      RepositoryUsername: '',
    };
  }

  return {
    RepositoryAuthentication: true,
    RepositoryPassword: '',
    RepositoryUsername: auth.Username,
  };
}

export function transformGitAuthenticationViewModel(
  auth?: GitAuthModel
): GitAuthenticationResponse | null {
  if (!auth || !auth.RepositoryAuthentication) {
    return null;
  }

  if (!auth.RepositoryUsername && !auth.RepositoryPassword) {
    return null;
  }

  return {
    Username: auth.RepositoryUsername,
    Password: auth.RepositoryPassword,
  };
}
