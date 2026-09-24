import React, {ReactNode, useRef, useState} from 'react';
import styles from './content.module.css';
import {CopyIcon} from './icons';

export type CodeCardStatus = 'tested' | 'warn' | 'note' | null;

export type CodeCardProps = {
  /** File name shown in the title bar. */
  file?: ReactNode;
  /** Right-aligned status chip. */
  status?: CodeCardStatus | ReactNode;
  /** Whether to show the copy button (default true). */
  copy?: boolean;
  children: ReactNode;
};

function renderStatus(status: CodeCardProps['status']) {
  if (!status) return null;
  if (status === 'tested') {
    return <span className={`${styles.codeTag} ${styles.ok}`}>tested</span>;
  }
  if (status === 'warn') {
    return <span className={`${styles.codeTag} ${styles.warn}`}>caveat</span>;
  }
  if (status === 'note') {
    return <span className={styles.codeTag}>note</span>;
  }
  if (typeof status === 'string') {
    return <span className={styles.codeTag}>{status}</span>;
  }
  return status;
}

export default function CodeCard({file, status, copy = true, children}: CodeCardProps) {
  const preRef = useRef<HTMLPreElement>(null);
  const [copied, setCopied] = useState(false);

  const onCopy = () => {
    if (!preRef.current) return;
    const text = preRef.current.innerText;
    if (typeof navigator !== 'undefined' && navigator.clipboard) {
      navigator.clipboard.writeText(text).then(() => {
        setCopied(true);
        setTimeout(() => setCopied(false), 1200);
      });
    }
  };

  return (
    <div className={styles.codeCard}>
      <div className={styles.codeBar}>
        <span className={styles.codeDot} />
        {file && <span className={styles.codeFile}>{file}</span>}
        <span className={styles.codeSpacer} />
        {renderStatus(status)}
        {copy && (
          <button type="button" className={styles.codeCopy} onClick={onCopy} aria-label="Copy code">
            {copied ? 'Copied' : <CopyIcon />}
          </button>
        )}
      </div>
      <pre ref={preRef} className={styles.codeBlock}>{children}</pre>
    </div>
  );
}
