import { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { registCardInfo, getUserInfo } from '../utils/api';
import styles from './Register.module.css';

function Register() {
  const [sid, setSid] = useState('');
  const navigate = useNavigate();
  const location = useLocation();
  const cardid = location.state?.cardid || '';

  const handleRegist = async () => {
    try {
      const res = await registCardInfo(cardid, sid);
      if (res.Success !== true) {
        console.log('Card register failed');
        alert('The ID is not found.');
      } else {
        const userinfo = await getUserInfo(sid);
        if (userinfo.IsEnter) {
          navigate('/question', { state: { userinfo } });
        } else {
          navigate('/welcome', { state: { userinfo } });
        }
      }
    } catch (error) {
      console.error(error);
      navigate('/');
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      handleRegist();
    }
  };

  return (
    <div className={styles.regist}>
      <div className={styles.studentNUM}>
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

export default Register;
