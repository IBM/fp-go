import React from 'react';
import Link from '@docusaurus/Link';
import styles from './styles.module.css';

const ExtIcon = () => (
  <svg className={styles.ext} width="11" height="11" viewBox="0 0 32 32" fill="currentColor" aria-hidden="true">
    <path d="M10 6v2h12.59L6 24.59 7.41 26 24 9.41V22h2V6z" />
  </svg>
);

const CopyIcon = () => (
  <svg width="14" height="14" viewBox="0 0 32 32" fill="currentColor" aria-hidden="true">
    <path d="M28 10v18H10V10h18m0-2H10a2 2 0 0 0-2 2v18a2 2 0 0 0 2 2h18a2 2 0 0 0 2-2V10a2 2 0 0 0-2-2z" />
    <path d="M4 18H2V4a2 2 0 0 1 2-2h14v2H4z" />
  </svg>
);

export default function Footer() {
  const year = new Date().getFullYear();

  const copyInstallCommand = () => {
    if (typeof navigator !== 'undefined' && navigator.clipboard) {
      navigator.clipboard.writeText('go get github.com/IBM/fp-go/v2');
    }
  };

  return (
    <footer className={styles.footer}>
      <div className={styles.top}>
        <div className={styles.grid}>
          <div className={styles.topGrid}>
            <div className={styles.lockup}>
              <div className={styles.wordmark}>
                fp-go <em>· functional Go</em>
              </div>
              <p>
                A practical, generics-powered functional toolkit for Go. Apache 2.0 — open source and IBM-maintained.
              </p>
              <div className={styles.install} title="Copy install command">
                <span className={styles.cmd}>
                  <span className={styles.p}>$</span>go get github.com/IBM/fp-go/v2
                </span>
                <button
                  type="button"
                  className={styles.cp}
                  onClick={copyInstallCommand}
                  aria-label="Copy install command">
                  <CopyIcon />
                </button>
              </div>
            </div>

            <nav className={styles.ftCol}>
              <h5>Documentation</h5>
              <Link to="/docs/intro">Getting Started</Link>
              <Link to="/docs/quickstart">Quick Start</Link>
              <Link to="/docs/installation">Installation</Link>
              <Link to="/docs/migration">Migration Guide</Link>
            </nav>

            <nav className={styles.ftCol}>
              <h5>Learn</h5>
              <Link to="/docs/concepts">Core Concepts</Link>
              <Link to="/docs/recipes">Recipes</Link>
              <Link to="/docs/v2/result">API Reference</Link>
              <Link to="/docs/faq">FAQ</Link>
            </nav>

            <nav className={styles.ftCol}>
              <h5>Community</h5>
              <Link to="https://github.com/IBM/fp-go">GitHub<ExtIcon /></Link>
              <Link to="https://github.com/IBM/fp-go/issues">Issues<ExtIcon /></Link>
              <Link to="https://github.com/IBM/fp-go/discussions">Discussions<ExtIcon /></Link>
              <Link to="https://github.com/IBM/fp-go/blob/main/CONTRIBUTING.md">Contributing<ExtIcon /></Link>
            </nav>

            <nav className={styles.ftCol}>
              <h5>More</h5>
              <Link to="https://ibm.github.io/">IBM Open Source<ExtIcon /></Link>
              <Link to="https://pkg.go.dev/github.com/IBM/fp-go/v2">Go Package<ExtIcon /></Link>
              <Link to="https://github.com/IBM/fp-go/blob/main/LICENSE">License<ExtIcon /></Link>
              <Link to="https://github.com/IBM/fp-go/security">Security</Link>
            </nav>
          </div>
        </div>
      </div>

      <div className={styles.bot}>
        <div className={styles.grid}>
          <div className={styles.botGrid}>
            <div className={styles.copy}>
              © {year} IBM Corporation · Licensed under <b>Apache 2.0</b> · Built with Docusaurus
            </div>
            <div className={styles.badges}>
              <span className={styles.badge}>
                <span className={styles.dot}></span> v2.2.82 stable
              </span>
              <span className={styles.badge}>Go 1.21+</span>
            </div>
            <div className={styles.legal}>
              <Link to="https://www.ibm.com/privacy">Privacy</Link>
              <Link to="https://www.ibm.com/legal">Terms</Link>
              <Link to="https://www.ibm.com/legal/copytrade">Trademarks</Link>
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
}
