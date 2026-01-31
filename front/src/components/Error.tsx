import { useLocation } from 'react-router';
import type { ErrorInfo } from '../types';
import styles from './Error.module.css';

function Error() {
    const location = useLocation();

    let errorInfo = location.state?.errorInfo as ErrorInfo | undefined;

    return (
        <div className={styles.error}>
            <h1>Under maintenance</h1>
            <div className={styles.message}>
                <p>{errorInfo ? errorInfo.message : 'Unknown error'}</p>
            </div>
        </div>
    );
}

export default Error;
