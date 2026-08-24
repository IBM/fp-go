import React, {ReactNode} from 'react';
import styles from './content.module.css';

export type CompareProps = {children: ReactNode};

export function Compare({children}: CompareProps) {
  return <div className={styles.compare}>{children}</div>;
}

export type CompareColProps = {
  kind: 'bad' | 'good';
  title?: ReactNode;
  /** Right-aligned pill in the column header. */
  pill?: ReactNode;
  children: ReactNode;
};

export function CompareCol({kind, title, pill, children}: CompareColProps) {
  return (
    <div className={`${styles.compareCol} ${kind === 'bad' ? styles.bad : styles.good}`}>
      <div className={styles.compareHead}>
        <span>{title ?? (kind === 'bad' ? 'Before' : 'After')}</span>
        {pill && <span className={styles.comparePill}>{pill}</span>}
      </div>
      <div className={styles.compareBody}>{children}</div>
    </div>
  );
}

export default Compare;
