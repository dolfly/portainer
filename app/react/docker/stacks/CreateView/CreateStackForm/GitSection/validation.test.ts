import { GitFormValues } from './types';
import { getGitValidationSchema } from './validation';

describe('Git validation', () => {
  it('should pass validation with valid git form values', async () => {
    const schema = getGitValidationSchema();

    const validData: GitFormValues = {
      RepositoryURL: 'https://github.com/user/repo',
      RepositoryReferenceName: 'refs/heads/main',
      ComposeFilePathInRepository: 'docker-compose.yml',
      RepositoryAuthentication: false,
      RepositoryUsername: '',
      RepositoryPassword: '',
      TLSSkipVerify: false,
      AdditionalFiles: [],
      AutoUpdate: undefined,
      RepositoryAuthorizationType: undefined,
      SupportRelativePath: false,
      FilesystemPath: '',
      SourceId: 1,
    };

    await expect(schema.validate(validData)).resolves.toBeDefined();
  });

  it('should fail validation when repository URL is empty', async () => {
    const schema = getGitValidationSchema();

    const invalidData = {
      RepositoryURL: '',
      RepositoryReferenceName: 'refs/heads/main',
      ComposeFilePathInRepository: 'docker-compose.yml',
    };

    await expect(schema.validate(invalidData)).rejects.toThrow();
  });
});
