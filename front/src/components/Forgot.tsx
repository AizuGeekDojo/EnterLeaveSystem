import { useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { getUserInfo } from '../utils/api';
import type { ErrorInfo } from '../types';
import styles from './Forgot.module.css';

function Forgot() {
  const [sid, setSid] = useState('');
  const navigate = useNavigate();
  const location = useLocation();

  let errorInfo = location.state?.errorInfo as ErrorInfo | undefined;

  const handleForgot = async () => {
    try {
      const res = await getUserInfo(sid);
      console.log(res);

      if (res.UserName === '') {
        console.log('Your ID is incorrect.');
        alert('Your ID is incorrect.');
      } else {
        if (res.IsEnter) {
          navigate('/question', { state: { userinfo: res } });
        } else {
          navigate('/welcome', { state: { userinfo: res } });
        }
      }
    } catch (error) {
      console.log(error);
      navigate('/');
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      handleForgot();
    }
  };

  return (
    <div className={styles.forgot}>
      <div className={styles.studentNUM}>
        {errorInfo && <p>{errorInfo.message}</p>}
        <h1>Input your student number</h1>
        <input
          type="text"
          value={sid}
          onChange={(e) => setSid(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="s13xxxxx"
        />
        <h4 style={{ color: '#aaaaaa' }}>
          <br />
          <span style={{ color: '#88aace' }}>Enter</span> key to continue -&gt;
        </h4>
      </div>
    </div>
  );
}

export default Forgot;
