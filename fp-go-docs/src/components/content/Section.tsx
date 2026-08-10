import React, {ReactNode} from 'react';
import styles from './content.module.css';

export type SectionProps = {
  /** DOM id for deep-linking. */
  id?: string;
  /** Mono-styled number/label (e.g. "01"). */
  number?: ReactNode;
  title: ReactNode;
  /** Italic-accent fragment appended to the title (renders inside <em>). */
  titleAccent?: ReactNode;
  /** Mono tag rendered on the right of the section head. */
  tag?: ReactNode;
  /** Short two-line intro under the section header. */
  lede?: ReactNode;
  children: ReactNode;
};

export default function Section({id, number, title, titleAccent, tag, lede, children}: SectionProps) {
  return (
    <section className={styles.section}>
      <div className={styles.sectionHead}>
        <div className={styles.sectionHeadLeft}>
          {number && <span className={styles.sectionNum}>{number}</span>}
          <h2 className={styles.sectionTitle} id={id}>
            {title}
            {titleAccent ? <> <em>{titleAccent}</em></> : null}
          </h2>
        </div>
        {tag && <span className={styles.sectionTag}>{tag}</span>}
      </div>
      {lede && <p className={styles.sectionLede}>{lede}</p>}
      {children}
    </section>
  );
}
