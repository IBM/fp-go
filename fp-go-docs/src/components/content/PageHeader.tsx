import React, {ReactNode} from 'react';
import styles from './content.module.css';

export type PageHeaderMeta = {label: string; value: ReactNode};

export type PageHeaderProps = {
  eyebrow?: ReactNode;
  title: ReactNode;
  /** Italic-accent fragment appended to the title (renders inside <em>). */
  titleAccent?: ReactNode;
  lede?: ReactNode;
  meta?: PageHeaderMeta[];
};

export default function PageHeader({eyebrow, title, titleAccent, lede, meta}: PageHeaderProps) {
  return (
    <header className={styles.head}>
      <div className={styles.headLeft}>
        {eyebrow && <div className={styles.eyebrow}>{eyebrow}</div>}
        <h1 className={styles.headTitle}>
          {title}
          {titleAccent ? <> <em>{titleAccent}</em></> : null}
        </h1>
        {lede && <p className={styles.lede}>{lede}</p>}
      </div>
      {meta && meta.length > 0 && (
        <aside className={styles.headMeta}>
          {meta.map((m, i) => (
            <div key={i} className={styles.metaRow}>
              <span className={styles.metaK}>{m.label}</span>
              <span className={styles.metaV}>{m.value}</span>
            </div>
          ))}
        </aside>
      )}
    </header>
  );
}

export function MetaPill({children}: {children: ReactNode}) {
  return <span className={styles.pill}>{children}</span>;
}
