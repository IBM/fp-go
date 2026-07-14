import React, {ReactNode} from 'react';
import styles from './content.module.css';

export type BenchRow = {
  label: ReactNode;
  /** 0–1 fill of the speed bar. */
  bar?: number;
  /** 'win' | 'lose' tints the bar. */
  barKind?: 'win' | 'lose' | 'neutral';
  /** Raw ns/op (or any numeric label). */
  nsOp?: ReactNode;
  /** Raw B/op. */
  bOp?: ReactNode;
  /** Delta text (e.g. "-85%" or "baseline"). */
  delta?: ReactNode;
  /** Good/bad colored delta. */
  deltaKind?: 'good' | 'bad';
  /** Highlight the row as the winner. */
  winner?: boolean;
};

export type BenchProps = {
  title?: ReactNode;
  /** Shell command shown in the head, e.g. "go test -bench=. -benchmem". */
  command?: ReactNode;
  /** Column headers. Default: ["Variant", "Speed", "ns/op", "B/op", "Δ"]. */
  columns?: [ReactNode, ReactNode, ReactNode, ReactNode, ReactNode];
  rows: BenchRow[];
};

export default function Bench({title, command, columns, rows}: BenchProps) {
  const [c1, c2, c3, c4, c5] = columns ?? ['Variant', 'Speed', 'ns/op', 'B/op', 'Δ'];
  return (
    <div className={styles.bench}>
      <div className={styles.benchHead}>
        <span>{title}</span>
        {command && (
          <span className={styles.benchCmd}>
            <span className={styles.benchPrompt}>$</span>
            {command}
          </span>
        )}
      </div>
      <table className={styles.benchTable}>
        <thead>
          <tr>
            <th>{c1}</th>
            <th style={{width: '40%'}}>{c2}</th>
            <th className={styles.benchNum}>{c3}</th>
            <th className={styles.benchNum}>{c4}</th>
            <th className={styles.benchNum}>{c5}</th>
          </tr>
        </thead>
        <tbody>
          {rows.map((r, i) => {
            const fillCls =
              r.barKind === 'win' ? styles.win : r.barKind === 'lose' ? styles.lose : '';
            const deltaCls =
              r.deltaKind === 'good' ? styles.deltaGood : r.deltaKind === 'bad' ? styles.deltaBad : styles.benchNum;
            return (
              <tr key={i} className={r.winner ? styles.benchWinner : undefined}>
                <td className={styles.benchLabel}>{r.label}</td>
                <td>
                  {r.bar != null && (
                    <span className={styles.bar}>
                      <span className={styles.barTrack}>
                        <span
                          className={`${styles.barFill} ${fillCls}`}
                          style={{width: `${Math.max(0, Math.min(1, r.bar)) * 100}%`}}
                        />
                      </span>
                    </span>
                  )}
                </td>
                <td className={styles.benchNum}>{r.nsOp}</td>
                <td className={styles.benchNum}>{r.bOp}</td>
                <td className={deltaCls}>{r.delta}</td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}
