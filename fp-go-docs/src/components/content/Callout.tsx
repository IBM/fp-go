import React, {ReactNode} from 'react';
import styles from './content.module.css';
import {InfoIcon, WarnIcon, CheckIcon} from './icons';

export type CalloutType = 'info' | 'warn' | 'success';

export type CalloutProps = {
  type?: CalloutType;
  title?: ReactNode;
  children: ReactNode;
};

const ICONS: Record<CalloutType, React.FC<{size?: number}>> = {
  info: InfoIcon,
  warn: WarnIcon,
  success: CheckIcon,
};

export default function Callout({type = 'info', title, children}: CalloutProps) {
  const Icon = ICONS[type];
  const typeCls = type === 'warn' ? styles.warn : type === 'success' ? styles.success : '';
  return (
    <div className={`${styles.callout} ${typeCls}`}>
      <span className={styles.calloutIcon}>
        <Icon size={18} />
      </span>
      <div className={styles.calloutBody}>
        {title && <strong className={styles.calloutTitle}>{title}</strong>}
        {children}
      </div>
    </div>
  );
}
