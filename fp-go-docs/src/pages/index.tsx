import type {ReactNode} from 'react';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import Heading from '@theme/Heading';

import styles from './index.module.css';

function HomepageHeader() {
  // NOTE: removed `clsx('hero hero--primary', ...)` — those Docusaurus
  // built-in classes paint a dark blue gradient on top of the section,
  // which makes the Carbon-style light hero content unreadable.
  return (
    <header className={styles.heroBanner}>
      <div className={styles.carbonGrid}>
        <div className={styles.heroInner}>
          <div className={styles.heroEyebrow}>v2.2.82 · Stable release</div>
          <Heading as="h1" className={styles.heroTitle}>
            Functional<br/>
            programming<br/>
            <span className={styles.heroAccent}>for Go.</span>
          </Heading>
          <p className={styles.heroTagline}>
            A practical, generics-powered toolkit for writing typed, composable Go.
            Zero dependencies. Production ready.
          </p>
          <div className={styles.heroCtas}>
            <Link className={`button button--primary button--lg ${styles.heroCta}`} to="/docs/intro">
              Get started
              <svg width="16" height="16" viewBox="0 0 32 32" fill="currentColor" aria-hidden="true">
                <path d="M18 6L16.6 7.4 24.2 15H4v2h20.2l-7.6 7.6L18 26l10-10z"/>
              </svg>
            </Link>
            <Link className={`button button--outline button--lg ${styles.heroCta}`} to="/docs/quickstart">
              Quick start
              <svg width="16" height="16" viewBox="0 0 32 32" fill="currentColor" aria-hidden="true">
                <path d="M11 23v-14l11 7-11 7z"/>
              </svg>
            </Link>
          </div>
        </div>
        <aside className={styles.heroMeta}>
          <div className={styles.metaRow}>
            <span className={styles.metaKey}>// Install</span>
            <div className={styles.installBar}>
              <span className={styles.prompt}>$</span>
              <span>go get github.com/IBM/fp-go/v2</span>
            </div>
          </div>
          <div className={styles.metaRow}>
            <span className={styles.metaKey}>// Requirements</span>
            <span className={styles.metaValue}>go 1.21+</span>
          </div>
          <div className={styles.metaRow}>
            <span className={styles.metaKey}>// License</span>
            <span className={styles.metaValue}>Apache 2.0 — free for commercial use</span>
          </div>
          <div className={styles.metaRow}>
            <span className={styles.metaKey}>// Last release</span>
            <span className={styles.metaValue}>v2.2.82 · 2 weeks ago</span>
          </div>
        </aside>
      </div>
    </header>
  );
}

function StatsStrip() {
  return (
    <section className={styles.statsSection}>
      <div className={styles.carbonGrid}>
        <div className={styles.stats}>
          <div className={styles.stat}>
            <div className={styles.statNum}>2<span>k</span></div>
            <div className={styles.statLabel}>GitHub stars</div>
          </div>
          <div className={styles.stat}>
            <div className={styles.statNum}>50<span>+</span></div>
            <div className={styles.statLabel}>Contributors</div>
          </div>
          <div className={styles.stat}>
            <div className={styles.statNum}>95<span>%</span></div>
            <div className={styles.statLabel}>Test coverage</div>
          </div>
          <div className={styles.stat}>
            <div className={styles.statNum}>0</div>
            <div className={styles.statLabel}>Dependencies</div>
          </div>
        </div>
      </div>
    </section>
  );
}

function FeaturesSection() {
  const features = [
    {
      num: '01',
      icon: (
        <svg width="32" height="32" viewBox="0 0 32 32" fill="currentColor" aria-hidden="true">
          <path d="M19 4h-2v8h-7l8 10 8-10h-7zM6 28V14H4v14a2 2 0 0 0 2 2h20v-2z"/>
        </svg>
      ),
      title: 'Fast & idiomatic',
      description: 'Generics-native. No reflection, no runtime cost. Feels like plain Go.',
    },
    {
      num: '02',
      icon: (
        <svg width="32" height="32" viewBox="0 0 32 32" fill="currentColor" aria-hidden="true">
          <path d="M16 30L7 25.5V11h2v13.3l7 3.5 7-3.5V11h2v14.5zM26 9h-2V5h-3V3h3V0h2v3h3v2h-3zM12 9h-2V5H7V3h3V0h2v3h3v2h-3z"/>
        </svg>
      ),
      title: 'Type-safe',
      description: 'Option, Either, Result, Try — full inference. Never interface{} again.',
    },
    {
      num: '03',
      icon: (
        <svg width="32" height="32" viewBox="0 0 32 32" fill="currentColor" aria-hidden="true">
          <path d="M28 22a3.99 3.99 0 0 0-3.86 3H17v-4h6v-8h2v-2h-6v2h2v6h-9V13a4 4 0 1 0-2 0v6H4v2h2v8h7.14a4 4 0 1 0 0-2H8v-6h7v4h-2v2h6v-2h-2v-4h7.14A4 4 0 1 0 28 22zM12 7a2 2 0 1 1-2 2 2 2 0 0 1 2-2zM10 26a2 2 0 1 1 2 2 2 2 0 0 1-2-2zm18 2a2 2 0 1 1 2-2 2 2 0 0 1-2 2z"/>
        </svg>
      ),
      title: 'Composable',
      description: 'Pipe, compose, fold. Build readable pipelines from one-screen functions.',
    },
    {
      num: '04',
      icon: (
        <svg width="32" height="32" viewBox="0 0 32 32" fill="currentColor" aria-hidden="true">
          <path d="M16 30A14 14 0 1 1 30 16 14 14 0 0 1 16 30zm0-26a12 12 0 1 0 12 12A12 12 0 0 0 16 4z"/>
          <path d="M14 21.59L8.41 16 7 17.41 14 24.41 25.41 13 24 11.59 14 21.59z"/>
        </svg>
      ),
      title: 'Battle-tested',
      description: 'Powers data pipelines at IBM and beyond. 95% coverage, fuzz-tested.',
    },
  ];

  return (
    <section className={styles.featuresSection}>
      <div className={styles.carbonGrid}>
        <div className={styles.sectionHead}>
          <div className={styles.sectionEyebrow}>Capabilities</div>
          <Heading as="h2" className={styles.sectionTitle}>
            Functional primitives, built for the <em>Go developer experience.</em>
          </Heading>
        </div>
        <div className={styles.featureRow}>
          {features.map((feature, idx) => (
            <article key={idx} className={styles.tile}>
              <div className={styles.tileNum}>{feature.num}</div>
              <div className={styles.tileIcon}>{feature.icon}</div>
              <Heading as="h3" className={styles.tileTitle}>{feature.title}</Heading>
              <p className={styles.tileDesc}>{feature.description}</p>
              <Link className={styles.tileMore} to="/docs/intro">
                Learn more
                <svg width="16" height="16" viewBox="0 0 32 32" fill="currentColor" aria-hidden="true">
                  <path d="M18 6L16.6 7.4 24.2 15H4v2h20.2l-7.6 7.6L18 26l10-10z"/>
                </svg>
              </Link>
            </article>
          ))}
        </div>
      </div>
    </section>
  );
}

function CodeSection() {
  return (
    <section className={styles.codeSection}>
      <div className={styles.carbonGrid}>
        <div className={styles.sectionHead}>
          <div className={styles.sectionEyebrow}>Example</div>
          <Heading as="h2" className={styles.sectionTitle}>
            Write <em>less.</em> Express more.
          </Heading>
        </div>
        <div className={styles.codeCopy}>
          <p>
            Replace <code>if err != nil</code> ladders with typed flows.
            Compose small, total functions into pipelines that are easy to read, refactor, and test.
          </p>
          <ul className={styles.bullets}>
            <li>
              <span className={styles.bulletNum}>01</span>
              <span><strong>Option[T]</strong> — nullable values without nil checks.</span>
            </li>
            <li>
              <span className={styles.bulletNum}>02</span>
              <span><strong>Either[E,A]</strong> & <strong>Result[T]</strong> — typed error flows.</span>
            </li>
            <li>
              <span className={styles.bulletNum}>03</span>
              <span><strong>Pipe</strong> & <strong>Compose</strong> — left-to-right data transforms.</span>
            </li>
            <li>
              <span className={styles.bulletNum}>04</span>
              <span><strong>Slice & Map</strong> helpers: filter, fold, group, partition, zip.</span>
            </li>
          </ul>
        </div>
        <div className={styles.codeCard}>
          <div className={styles.codeTabs}>
            <div className={`${styles.codeTab} ${styles.codeTabActive}`}>main.go</div>
            <div className={styles.codeTab}>option.go</div>
            <div className={styles.codeTab}>pipe_test.go</div>
          </div>
          <div className={styles.codeBar}>
            <span className={styles.codeDot}></span>
            <span>fp-go · greet.go</span>
          </div>
          <pre className={styles.codeBlock}><code>{`package main

import (
    "fmt"
    "strings"
    F "github.com/IBM/fp-go/v2"
    O "github.com/IBM/fp-go/v2/option"
)

// Parse → validate → format. Total, typed, nil-safe.
func greet(raw string) string {
    return F.Pipe3(
        O.FromString(raw),
        O.Filter(func(s string) bool {
            return len(s) > 1
        }),
        O.Map(strings.Title),
        O.GetOrElse(F.Constant("friend")),
    )
}

func main() {
    fmt.Println("Hello, " + greet("ada"))
    // → Hello, Ada
}`}</code></pre>
        </div>
      </div>
    </section>
  );
}

export default function Home(): ReactNode {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title={`${siteConfig.title} - Functional Programming for Go`}
      description="Type-safe functional programming library for Go with monads, functors, and more">
      <HomepageHeader />
      <main>
        <StatsStrip />
        <FeaturesSection />
        <CodeSection />
      </main>
    </Layout>
  );
}
