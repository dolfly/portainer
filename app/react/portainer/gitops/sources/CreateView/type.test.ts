import { formValuesToCreatePayload, gitFormValuesToTestPayload } from './type';

const baseGit = {
  url: 'https://github.com/org/repo.git',
  tlsSkipVerify: false,
  connectionOk: false,
};

describe('formValuesToCreatePayload', () => {
  it('populates authentication when authEnabled with username and password', () => {
    const payload = formValuesToCreatePayload({
      name: 'my-source',
      type: 'git',
      git: {
        ...baseGit,
        authentication: {
          authEnabled: true,
          username: 'alice',
          password: 'secret',
        },
      },
    });

    expect(payload.git.authentication).toEqual({
      username: 'alice',
      password: 'secret',
    });
  });

  it('omits authentication when authEnabled is false', () => {
    const payload = formValuesToCreatePayload({
      name: 'my-source',
      type: 'git',
      git: {
        ...baseGit,
        authentication: { authEnabled: false },
      },
    });

    expect(payload.git.authentication).toBeUndefined();
  });

  it('omits authentication when authEnabled but username is missing', () => {
    const payload = formValuesToCreatePayload({
      name: 'my-source',
      type: 'git',
      git: {
        ...baseGit,
        authentication: {
          authEnabled: true,
          password: 'secret',
        },
      },
    });

    expect(payload.git.authentication).toBeUndefined();
  });

  it('omits authentication when authEnabled but password is missing', () => {
    const payload = formValuesToCreatePayload({
      name: 'my-source',
      type: 'git',
      git: {
        ...baseGit,
        authentication: {
          authEnabled: true,
          username: 'alice',
        },
      },
    });

    expect(payload.git.authentication).toBeUndefined();
  });

  it('does not include connectionOk in the create payload', () => {
    const payload = formValuesToCreatePayload({
      name: 'my-source',
      type: 'git',
      git: {
        ...baseGit,
        connectionOk: true,
        authentication: { authEnabled: false },
      },
    });

    expect(payload.git).not.toHaveProperty('connectionOk');
  });
});

describe('gitFormValuesToTestPayload', () => {
  it('populates authentication when authEnabled with username and password', () => {
    const payload = gitFormValuesToTestPayload({
      ...baseGit,
      authentication: {
        authEnabled: true,
        username: 'alice',
        password: 'secret',
      },
    });

    expect(payload.authentication).toEqual({
      username: 'alice',
      password: 'secret',
    });
  });

  it('omits authentication when authEnabled is false', () => {
    const payload = gitFormValuesToTestPayload({
      ...baseGit,
      authentication: { authEnabled: false },
    });

    expect(payload.authentication).toBeUndefined();
  });
});
