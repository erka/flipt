import { createContext, useEffect, useMemo, useState } from 'react';
import { getAuthSelf, getInfo } from '~/data/api';
import { IAuthGithubInternal } from '~/types/auth/Github';
import { IAuthOIDCInternal } from '~/types/auth/OIDC';
import { expireAuthSelf } from '~/data/api';
import { useError } from '~/data/hooks/error';
import { redirect } from 'react-router-dom';

type Session = {
  authenticated: boolean;
  self?: IAuthOIDCInternal | IAuthGithubInternal;
};

interface SessionContextType {
  session?: Session;
  login: (uri: string) => void;
  logout: () => void;
}

export const SessionContext = createContext({} as SessionContextType);

const blank = {
  authenticated: false
};

export default function SessionProvider({
  children
}: {
  children: React.ReactNode;
}) {
  const [session, setSession] = useState<Session>(blank);
  const { setError, clearError } = useError();
  useEffect(() => {
    const loadSession = async () => {
      let session = {
        authenticated: false
      } as Session;

      try {
        await getInfo();
        session.authenticated = true;
        session.self = await getAuthSelf();
      } catch (err) {
        // TODO: check the response if it's 401
      }
      setSession(session);
    };
    loadSession();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const logout = async () => {
    expireAuthSelf()
      .then(() => {
        setSession(blank);
        redirect('/login');
      })
      .catch((err) => {
        setError(err);
      });
  };

  const login = async (uri: string) => {
    const res = await fetch(uri, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    });

    if (!res.ok || res.status !== 200) {
      const { message } = await res.json();
      setError('Unable to authenticate: ' + message);
      return;
    }

    clearError();
    const body = await res.json();
    window.location.href = body.authorizeUrl;
  };

  const value = useMemo(
    () => ({
      session,
      login,
      logout
    }),
    [session, setSession, logout]
  );

  return (
    <SessionContext.Provider value={value}>{children}</SessionContext.Provider>
  );
}
