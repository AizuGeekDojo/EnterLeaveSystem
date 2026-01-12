import { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { registerCardInfo, getUserInfo } from '../utils/api';
import styles from './Register.module.css';

function Register() {
  const [ainsID, setAinsID] = useState('');
  const navigate = useNavigate();
  const location = useLocation();
  const cardid = location.state?.cardid || '';

  const handleRegister = async () => {
    try {
      const res = await registerCardInfo(cardid, ainsID);
      if (res.Success !== true) {
        console.log('Card register failed');
        alert(`Your AINS ID (${ainsID}) is not registered in the system\nPlease contact to administrator.`);
      } else {
        const userinfo = await getUserInfo(ainsID);
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
      handleRegister();
    }
  };

  return (
    <div className={styles.register}>
      <div className={styles.ainsID}>
        <h1>Please input your AINS ID</h1>
        <input
          type="text"
          value={ainsID}
          onChange={(e) => setAinsID(e.target.value)}
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
