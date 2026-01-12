import { useState, useRef } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { addLog } from '../utils/api';
import type { UserInfo } from '../types';
import styles from './Question.module.css';

function Question() {
  const [checkedUse, setCheckedUse] = useState<string[]>([]);
  const [message, setMessage] = useState('');
  const [sending, setSending] = useState(false);
  const sendBtnRef = useRef<HTMLButtonElement>(null);
  const navigate = useNavigate();
  const location = useLocation();

  const handleCheckboxChange = (value: string) => {
    if (checkedUse.includes(value)) {
      setCheckedUse(checkedUse.filter((item) => item !== value));
    } else {
      setCheckedUse([...checkedUse, value]);
    }
  };

  const handleSend = () => {
    if (sending) {
      return;
    }

    setSending(true);
    if (sendBtnRef.current) {
      sendBtnRef.current.disabled = true;
    }

    const userinfo = location.state?.userinfo as UserInfo | undefined;
    const answer = JSON.stringify({
      Use: checkedUse,
      message: message,
    });

    if (!userinfo) {
      navigate('/');
      return;
    }

    addLog(userinfo.AINSID, false, answer);
    navigate('/goodbye');
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.ctrlKey && e.key === 'Enter') {
      handleSend();
    }
  };

  return (
    <div className={styles.ques}>
      <div className={styles.question}>
        <h2>What is your purpose?</h2>
        <div className={styles.checkboxes}>
          {['3DPrinter', 'LaserCutter', 'Training session', 'Other'].map((item) => (
            <div key={item} className={styles.checkbox}>
              <input
                type="checkbox"
                value={item}
                checked={checkedUse.includes(item)}
                onChange={() => handleCheckboxChange(item)}
              />
              <label style={{ fontSize: '18px' }}>
                {item === 'Training session' ? 'Training' : item}
              </label>
            </div>
          ))}
        </div>
        <h3>If you have any request, please fill in.</h3>
        <br />
        <textarea
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder=""
        />
        <br />
        <button ref={sendBtnRef} className="btn btn-info" onClick={handleSend}>
          send
        </button>
        <br />
        <h5 style={{ color: '#aaaaaa' }}>
          <span style={{ color: '#88aace' }}>Ctrl + Enter</span> to submit your request.
        </h5>
      </div>
    </div>
  );
}

export default Question;
