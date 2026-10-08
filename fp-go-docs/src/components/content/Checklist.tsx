import React, {ReactNode, ReactElement, Children} from 'react';
import styles from './content.module.css';
import {CheckIcon} from './icons';

export type ChecklistItemType = {
  label: ReactNode;
  impact?: ReactNode;
  done?: boolean;
};

export type ChecklistItemProps = {
  status?: 'required' | 'recommended' | 'optional';
  children: ReactNode;
};

export function ChecklistItem({status, children}: ChecklistItemProps): ReactElement {
  // This component is only used for JSX syntax, actual rendering happens in Checklist
  return <></>;
}

export type ChecklistProps = {
  title?: ReactNode;
  items?: ChecklistItemType[];
  children?: ReactNode;
};

export default function Checklist({title = 'Steps', items, children}: ChecklistProps) {
  // Support both array-based and JSX children patterns
  let checklistItems: ChecklistItemType[];
  
  if (items) {
    // Array-based pattern
    checklistItems = items;
  } else if (children) {
    // JSX children pattern
    checklistItems = Children.toArray(children)
      .filter((child): child is ReactElement<ChecklistItemProps> => React.isValidElement(child))
      .map((child) => ({
        label: child.props.children,
        impact: child.props.status,
        done: false, // JSX pattern doesn't track completion
      }));
  } else {
    checklistItems = [];
  }

  const done = checklistItems.filter((i) => i.done).length;
  
  return (
    <div className={styles.check}>
      <div className={styles.checkHead}>
        <span>{title}</span>
        {items && <span>{done} / {checklistItems.length} complete</span>}
      </div>
      <ul className={styles.checkList}>
        {checklistItems.map((it, i) => (
          <li key={i} className={`${styles.checkItem} ${it.done ? styles.done : ''}`}>
            <span className={styles.checkBox}>
              {it.done && <CheckIcon size={10} />}
            </span>
            <span className={styles.checkLbl}>{it.label}</span>
            {it.impact && <span className={styles.checkImpact}>{it.impact}</span>}
          </li>
        ))}
      </ul>
    </div>
  );
}
