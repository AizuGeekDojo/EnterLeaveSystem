import type { UserInfo, RegistResponse } from '../types';

export const getUserInfo = async (ainsID: string): Promise<UserInfo> => {
  const response = await fetch(`http://localhost:3000/api/user?AINSID=${ainsID}`, {
    mode: 'cors',
    method: 'GET',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
  });

  if (!response.ok) {
    throw response;
  }

  return response.json();
};

export const registerCardInfo = async (
  cardid: string,
  ainsID: string
): Promise<RegistResponse> => {
  const response = await fetch('http://localhost:3000/api/user', {
    mode: 'cors',
    method: 'POST',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      AINSID: ainsID,
      CardID: cardid,
    }),
  });

  return response.json();
};

export const addLog = async (
  ainsID: string,
  isenter: boolean,
  ext: string
): Promise<Response | null> => {
  console.log({
    ainsID: ainsID,
    IsEnter: isenter,
    Ext: ext,
  });

  try {
    const response = await fetch('http://localhost:3000/api/log', {
      mode: 'cors',
      method: 'POST',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        AINSID: ainsID,
        IsEnter: isenter,
        Ext: ext,
      }),
    });

    if (!response.ok) {
      throw response;
    }

    return response;
  } catch (error) {
    console.error(error);
    return null;
  }
};

export const roomName = (): string => {
  const roomname = import.meta.env.VITE_ROOMNAME;
  if (roomname === undefined) {
    return 'University of Aizu';
  }
  return roomname;
};

export const isShowQuestion = (): boolean => {
  const isshow = import.meta.env.VITE_SHOWQUESTION;
  if (isshow === undefined) {
    return false;
  }
  return isshow === 'true';
};
