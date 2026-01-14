import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import styles from './Goodbye.module.css';

function Goodbye() {
  const navigate = useNavigate();

  useEffect(() => {
    const timer = setTimeout(() => {
      navigate('/');
    }, 3000);

    return () => clearTimeout(timer);
  }, [navigate]);

  return (
    <div className={styles.goodbye}>
      <h1>Bye!</h1>
    </div>
  );
}

export default Goodbye;
