import { useEffect, useRef, useState } from 'react';
import { useNavigate, Link } from 'react-router';
import { getUserInfo, roomName } from '../utils/api';
import type { CardMessage } from '../types';
import styles from './Top.module.css';

function Top() {
  const [room] = useState(roomName());
  const navigate = useNavigate();
  const wsRef = useRef<WebSocket | null>(null);
  const closeFlgRef = useRef(false);
  const reconnectTimerRef = useRef<number | null>(null);

  const connectCardReader = () => {
    const ws = new WebSocket('ws://localhost:3000/socket/readCard');
    wsRef.current = ws;

    ws.onopen = () => {
      console.log('Card reader standby');
    };

    ws.onmessage = (e) => {
      const message: CardMessage = JSON.parse(e.data);
      console.log('Read card data:', message);

      if (message.IsCard === true) {
        if (message.IsNew === false && message.SID) {
          getUserInfo(message.SID)
            .then((res) => {
              if (res.IsEnter) {
                navigate('/question', { state: { userinfo: res } });
              } else {
                navigate('/welcome', { state: { userinfo: res } });
              }
            })
            .catch((error) => {
              console.error(error);
            });
        } else if (message.CardID) {
          navigate('/regist', { state: { cardid: message.CardID } });
        }
      }
    };

    ws.onerror = (e) => {
      console.log('Card reader communication error', e);
    };

    ws.onclose = () => {
      console.log('Card reader stopped');
      if (!closeFlgRef.current) {
        reconnectTimerRef.current = window.setTimeout(() => {
          window.location.reload();
        }, 3000);
      }
    };
  };

  useEffect(() => {
    connectCardReader();

    return () => {
      closeFlgRef.current = true;
      if (reconnectTimerRef.current) {
        clearTimeout(reconnectTimerRef.current);
      }
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, []);

  return (
    <div className={styles.top}>
      <div className={styles.message}>
        <h1>
          Put your card over the reader
          <br />
        </h1>
        <Link to="/forgot" className={styles.forgotPos}>
          Forgot card?
        </Link>
      </div>
      <h4 className={styles.roomPos}>{room}</h4>
    </div>
  );
}

export default Top;
