import React, {ReactNode} from 'react';
import styles from './content.module.css';

export type ApiRow = {
  symbol: ReactNode;
  signature: ReactNode;
  description?: ReactNode;
  since?: ReactNode;
};

export type ApiTableProps = {
  columns?: [ReactNode, ReactNode, ReactNode];
  rows?: ApiRow[];
  children?: ReactNode;
};

export default function ApiTable({columns, rows, children}: ApiTableProps) {
  // Support both rows prop and markdown table children
  if (children) {
    // Pass through markdown table as-is
    return <div className={styles.api}>{children}</div>;
  }
  
  // Original array-based rendering
  const [c1, c2, c3] = columns ?? ['Symbol', 'Signature', 'Since'];
  return (
    <div className={styles.api}>
      <table className={styles.apiTable}>
        <thead>
          <tr>
            <th>{c1}</th>
            <th>{c2}</th>
            <th>{c3}</th>
          </tr>
        </thead>
        <tbody>
          {rows?.map((r, i) => (
            <tr key={i}>
              <td><code>{r.symbol}</code></td>
              <td>
                <code>{r.signature}</code>
                {r.description && <div className={styles.apiDesc}>{r.description}</div>}
              </td>
              <td>{r.since && <code>{r.since}</code>}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
