import React, {useCallback, useEffect, useMemo, useState} from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import {useLocation} from '@docusaurus/router';
import {isActiveSidebarItem} from '@docusaurus/plugin-content-docs/client';
import type {Props} from '@theme/DocSidebar/Desktop';
import type {
  PropSidebarItem,
  PropSidebarItemCategory,
  PropSidebarItemLink,
} from '@docusaurus/plugin-content-docs';

import styles from './styles.module.css';

const STORAGE_KEY = 'fp-go-docs:visited';
const PADDING = 2;

const SearchIcon = () => (
  <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" strokeWidth={2} aria-hidden="true">
    <path d="M21 21l-4.3-4.3M17 11a6 6 0 1 1-12 0 6 6 0 0 1 12 0z" strokeLinecap="round" />
  </svg>
);

const ChevIcon = ({className}: {className?: string}) => (
  <svg className={className} viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" strokeWidth={2} aria-hidden="true">
    <path d="M9 6l6 6-6 6" strokeLinecap="round" strokeLinejoin="round" />
  </svg>
);

const PlayIcon = () => (
  <svg viewBox="0 0 24 24" width="14" height="14" fill="currentColor" aria-hidden="true">
    <path d="M8 5v14l11-7-11-7z" />
  </svg>
);

const pad = (n: number) => String(n).padStart(PADDING, '0');

function collectLinks(items: readonly PropSidebarItem[]): PropSidebarItemLink[] {
  const out: PropSidebarItemLink[] = [];
  for (const it of items) {
    if (it.type === 'link') out.push(it);
    else if (it.type === 'category') out.push(...collectLinks(it.items));
  }
  return out;
}

function itemMatchesQuery(item: PropSidebarItem, q: string): boolean {
  if (!q) return true;
  const needle = q.toLowerCase();
  if (item.type === 'link') return item.label.toLowerCase().includes(needle);
  if (item.type === 'category') {
    if (item.label.toLowerCase().includes(needle)) return true;
    return item.items.some((c) => itemMatchesQuery(c, needle));
  }
  return false;
}

function useVisited() {
  const [visited, setVisited] = useState<Set<string>>(() => new Set());

  useEffect(() => {
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      if (raw) setVisited(new Set(JSON.parse(raw)));
    } catch {
      // ignore
    }
  }, []);

  const mark = useCallback((href: string) => {
    setVisited((prev) => {
      if (prev.has(href)) return prev;
      const next = new Set(prev);
      next.add(href);
      try {
        localStorage.setItem(STORAGE_KEY, JSON.stringify([...next]));
      } catch {
        // ignore
      }
      return next;
    });
  }, []);

  return {visited, mark};
}

type NumberedItemProps = {
  item: PropSidebarItemLink;
  num: number;
  active: boolean;
  done: boolean;
};

function NumberedItem({item, num, active, done}: NumberedItemProps) {
  return (
    <Link
      className={clsx(styles.item, done && styles.itemDone, active && styles.itemActive)}
      to={item.href}
      isNavLink
      activeClassName={styles.itemActive}>
      <span className={styles.num}>{pad(num)}</span>
      <span className={styles.title}>{item.label}</span>
      {done && !active && <span className={styles.check} aria-hidden="true">✓</span>}
    </Link>
  );
}

type GroupProps = {
  category: PropSidebarItemCategory;
  links: PropSidebarItemLink[];
  numFor: (href: string) => number;
  isActiveHref: (href: string) => boolean;
  visited: Set<string>;
};

function Group({category, links, numFor, isActiveHref, visited}: GroupProps) {
  return (
    <div className={styles.group}>
      <h5 className={styles.groupHead}>
        <span>{category.label}</span>
        <span className={styles.ct}>{pad(links.length)}</span>
      </h5>
      {links.map((leaf) => (
        <NumberedItem
          key={leaf.href}
          item={leaf}
          num={numFor(leaf.href)}
          active={isActiveHref(leaf.href)}
          done={visited.has(leaf.href)}
        />
      ))}
    </div>
  );
}

type CollapsibleGroupProps = GroupProps & {
  hasActive: boolean;
};

function CollapsibleGroup({
  category,
  links,
  numFor,
  isActiveHref,
  visited,
  hasActive,
}: CollapsibleGroupProps) {
  const [open, setOpen] = useState(hasActive);

  useEffect(() => {
    if (hasActive) setOpen(true);
  }, [hasActive]);

  return (
    <>
      <button
        type="button"
        className={clsx(styles.collapsible, open && styles.collapsibleOpen)}
        onClick={() => setOpen((v) => !v)}
        aria-expanded={open}>
        <span>{category.label}</span>
        <span className={styles.right}>
          <span>{pad(links.length)}</span>
          <ChevIcon className={styles.chev} />
        </span>
      </button>
      {open && (
        <div className={styles.subGroup}>
          {links.map((leaf) => (
            <NumberedItem
              key={leaf.href}
              item={leaf}
              num={numFor(leaf.href)}
              active={isActiveHref(leaf.href)}
              done={visited.has(leaf.href)}
            />
          ))}
        </div>
      )}
    </>
  );
}

function DocSidebarDesktop({path, sidebar, isHidden}: Props) {
  const {pathname} = useLocation();
  const activePath = path ?? pathname;

  const [query, setQuery] = useState('');
  const {visited, mark} = useVisited();

  const allLinks = useMemo(() => collectLinks(sidebar), [sidebar]);
  const numByHref = useMemo(() => {
    const m = new Map<string, number>();
    allLinks.forEach((l, i) => m.set(l.href, i + 1));
    return m;
  }, [allLinks]);

  const numFor = useCallback((href: string) => numByHref.get(href) ?? 0, [numByHref]);

  const isActiveHref = useCallback(
    (href: string) => isActiveSidebarItem({type: 'link', href, label: ''}, activePath),
    [activePath],
  );

  // Mark the current page as visited.
  useEffect(() => {
    const link = allLinks.find((l) => isActiveSidebarItem(l, activePath));
    if (link) mark(link.href);
  }, [activePath, allLinks, mark]);

  const totalCount = allLinks.length;
  const doneCount = allLinks.filter((l) => visited.has(l.href)).length;
  const remainingMin = Math.max(0, (totalCount - doneCount) * 4);

  const filteredSidebar = useMemo(() => {
    const q = query.trim();
    if (!q) return sidebar;
    return sidebar.filter((it) => itemMatchesQuery(it, q));
  }, [sidebar, query]);

  // First two top-level categories render as expanded "groups"; the rest as collapsibles.
  const topCategories = filteredSidebar.filter((i) => i.type === 'category') as PropSidebarItemCategory[];
  const topLinks = filteredSidebar.filter((i) => i.type === 'link') as PropSidebarItemLink[];
  const expandedFirst = topCategories.slice(0, 2);
  const collapsedRest = topCategories.slice(2);

  return (
    <aside className={clsx(styles.sidebar, isHidden && styles.sidebarHidden)}>
      <div className={styles.search}>
        <SearchIcon />
        <input
          type="search"
          placeholder="Search the docs…"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          aria-label="Filter sidebar"
        />
        <kbd>⌘K</kbd>
      </div>

      <div className={styles.header}>
        <div className={styles.eyebrow}>— Getting started</div>
        <h4>
          Learn <em>fp-go</em>
        </h4>
        <div className={styles.meta}>
          <span className={styles.dot}></span>
          <span>{doneCount} of {totalCount} done</span>
          {remainingMin > 0 && (
            <>
              <span>·</span>
              <span>~{remainingMin} min left</span>
            </>
          )}
        </div>
      </div>

      <nav className={styles.list} aria-label="Docs sidebar">
        {expandedFirst.map((cat) => {
          const links = collectLinks(cat.items);
          return (
            <Group
              key={cat.label}
              category={cat}
              links={links}
              numFor={numFor}
              isActiveHref={isActiveHref}
              visited={visited}
            />
          );
        })}

        {topLinks.length > 0 && (
          <div className={styles.group}>
            <h5 className={styles.groupHead}>
              <span>More</span>
              <span className={styles.ct}>{pad(topLinks.length)}</span>
            </h5>
            {topLinks.map((leaf) => (
              <NumberedItem
                key={leaf.href}
                item={leaf}
                num={numFor(leaf.href)}
                active={isActiveHref(leaf.href)}
                done={visited.has(leaf.href)}
              />
            ))}
          </div>
        )}

        {collapsedRest.map((cat) => {
          const links = collectLinks(cat.items);
          const hasActive = links.some((l) => isActiveHref(l.href));
          return (
            <CollapsibleGroup
              key={cat.label}
              category={cat}
              links={links}
              numFor={numFor}
              isActiveHref={isActiveHref}
              visited={visited}
              hasActive={hasActive}
            />
          );
        })}

        {filteredSidebar.length === 0 && (
          <div className={styles.empty}>No matches for “{query}”.</div>
        )}
      </nav>

      <Link to="/playground" className={styles.cta}>
        <span className={styles.ctaIcon}>
          <PlayIcon />
        </span>
        <span className={styles.ctaTxt}>
          <b>Interactive Playground</b>
          <span>Run fp-go in your browser</span>
        </span>
      </Link>
    </aside>
  );
}

export default React.memo(DocSidebarDesktop);
