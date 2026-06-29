import React, {ReactNode} from 'react';
import Link from '@docusaurus/Link';
import styles from './content.module.css';
import {ArrowLeft, ArrowRight} from './icons';

export type PagerLink = {to: string; title: ReactNode; label?: ReactNode};

export type PagerProps = {
  prev?: PagerLink;
  next?: PagerLink;
};

export default function Pager({prev, next}: PagerProps) {
  return (
    <nav className={styles.pager}>
      {prev ? (
        <Link to={prev.to} className={styles.pagerLink}>
          <span className={styles.pagerLbl}>
            <ArrowLeft /> {prev.label ?? 'Previous'}
          </span>
          <span className={styles.pagerTtl}>{prev.title}</span>
        </Link>
      ) : (
        <span />
      )}
      {next ? (
        <Link to={next.to} className={`${styles.pagerLink} ${styles.next}`}>
          <span className={styles.pagerLbl}>
            {next.label ?? 'Next'} <ArrowRight />
          </span>
          <span className={styles.pagerTtl}>{next.title}</span>
        </Link>
      ) : (
        <span />
      )}
    </nav>
  );
}
