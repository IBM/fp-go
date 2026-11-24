import React, {type ReactElement} from 'react';
import Navbar from '@theme-original/Navbar';
import type NavbarType from '@theme/Navbar';
import type {WrapperProps} from '@docusaurus/types';

type Props = WrapperProps<typeof NavbarType>;

export default function NavbarWrapper(props: Props): ReactElement {
  return (
    <>
      <Navbar {...props} />
      <style>{`
        .navbar__brand::after {
          content: 'v2';
          font-family: var(--ifm-font-family-monospace);
          font-size: 10px;
          color: #0043ce;
          background: #d0e2ff;
          padding: 2px 6px;
          margin-left: 6px;
          letter-spacing: 0.02em;
          text-transform: uppercase;
          border-radius: 0;
          display: inline-block;
        }
        
        [data-theme='dark'] .navbar__brand::after {
          color: #78a9ff;
          background: #0043ce;
        }
      `}</style>
    </>
  );
}

// Made with Bob
