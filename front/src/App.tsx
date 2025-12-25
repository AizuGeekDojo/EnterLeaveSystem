import { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import Top from './components/Top';
import Register from './components/Register';
import Welcome from './components/Welcome';
import Goodbye from './components/Goodbye';
import Question from './components/Question';
import Forgot from './components/Forgot';
import styles from './App.module.css';

function App() {
  const [clockText, setClockText] = useState('0000/00/00  00:00:00');

  useEffect(() => {
    const timer = setInterval(() => {
      const da = new Date();
      const year = da.getFullYear();
      const month = da.getMonth() + 1;
      const date = da.getDate();
      const hour = da.getHours();
      const minute = da.getMinutes();
      const second = da.getSeconds();

      let text = `${year}/`;
      if (month < 10) {
        text += '0';
      }
      text += `${month}/`;
      if (date < 10) {
        text += '0';
      }
      text += `${date}  `;
      if (hour < 10) {
        text += '0';
      }
      text += `${hour}:`;
      if (minute < 10) {
        text += '0';
      }
      text += `${minute}:`;
      if (second < 10) {
        text += '0';
      }
      text += `${second}`;

      setClockText(text);
    }, 1000);

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
          </Routes>
        </div>
      </div>
    </Router>
  );
}

export default App;
