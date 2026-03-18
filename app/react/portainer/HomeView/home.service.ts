import axios, { parseAxiosError } from 'Portainer/services/axios/axios';

import { Motd } from './types';

export async function getMotd() {
  try {
    const { data } = await axios.get<Motd>('/motd');
    return data;
  } catch (err) {
    throw parseAxiosError(
      err as Error,
      'Unable to retrieve information message'
    );
  }
}
