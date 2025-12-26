import { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import dayjs from 'dayjs';
import Top from './components/Top';
import Register from './components/Register';
import Welcome from './components/Welcome';
import Goodbye from './components/Goodbye';
import Question from './components/Question';
import Forgot from './components/Forgot';
import Error from './components/Error';
import styles from './App.module.css';

function App() {
  const [clockText, setClockText] = useState('0000/00/00  00:00:00');

  useEffect(() => {
    const timer = setInterval(() => {
      const now = dayjs();
      const text = now.format('YYYY/MM/DD  HH:mm:ss');
      setClockText(text);
    }, 500);

    return () => clearInterval(timer);
  }, []);

  return (
    <Router>
      <div className={styles.background}>
        <div className={styles.container}>
          <div className={styles.topNav}>
            <Link to="/" className={styles.totopPos}>
              ←
            </Link>
          </div>
          <h2 className={styles.clockPos}>{clockText}</h2>
          <Routes>
            <Route path="/" element={<Top />} />
            <Route path="/regist" element={<Register />} />
            <Route path="/welcome" element={<Welcome />} />
            <Route path="/goodbye" element={<Goodbye />} />
            <Route path="/question" element={<Question />} />
            <Route path="/forgot" element={<Forgot />} />
            <Route path="/error" element={<Error />} />
          </Routes>
        </div>
      </div>
    </Router>
  );
}

export default App;
