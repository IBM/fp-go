import React, {ReactNode} from 'react';
import styles from './content.module.css';

export type TLDRProps = {children: ReactNode};

export function TLDR({children}: TLDRProps) {
  return <div className={styles.tldr}>{children}</div>;
}

export type TLDRCardProps = {
  label?: ReactNode;
  value: ReactNode;
  /** Italic accent rendered after the value (small serif italic). */
  accent?: ReactNode;
  /** Mono unit rendered after the value (e.g. "per op"). */
  unit?: ReactNode;
  description?: ReactNode;
  /** Tint the value green (up) or red (down). */
  variant?: 'default' | 'up' | 'down';
  /** Render the value as prose-sized text instead of a big number. */
  prose?: boolean;
};

export function TLDRCard({label, value, accent, unit, description, variant = 'default', prose}: TLDRCardProps) {
  const variantCls = variant === 'up' ? styles.up : variant === 'down' ? styles.down : '';
  return (
    <div className={`${styles.tldrCard} ${variantCls}`}>
      {label && <div className={styles.tldrK}>{label}</div>}
      <div className={`${styles.tldrV} ${prose ? styles.tldrVProse : ''}`}>
        {value}
        {accent ? <em>{accent}</em> : null}
        {unit ? <small>{unit}</small> : null}
      </div>
      {description && <div className={styles.tldrD}>{description}</div>}
    </div>
  );
}

export default TLDR;
