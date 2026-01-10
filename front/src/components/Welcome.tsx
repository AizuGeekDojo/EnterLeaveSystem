import { useEffect, useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { addLog, roomName } from '../utils/api';
import type { UserInfo } from '../types';
import styles from './Welcome.module.css';

function Welcome() {
  const [username, setUsername] = useState('');
  const [room] = useState(roomName());
  const navigate = useNavigate();
  const location = useLocation();

  let userinfo = location.state?.userinfo as UserInfo | undefined | null;
  useEffect(() => {
    if (userinfo) {
      setUsername(userinfo.UserName);
      addLog(userinfo.SID, true, '');
    }

    const timer = setTimeout(() => {
      navigate('/');
    }, 5000);

    return () => { userinfo = null; clearTimeout(timer) };
  }, [location.state.userinfo, navigate]);

  return (
    <div className={styles.welcome}>
      {username && // usernameがある
        <h1>
          Welcome to {room},<br/>{username}!
        </h1>
      }
      {username === '' && // usernameがない: 新入生等DBに未登録のユーザー
        <h1>
          Welcome to {room}!
        </h1>
      }
    </div>
  );
}

export default Welcome;
