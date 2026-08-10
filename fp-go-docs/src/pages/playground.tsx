import React from 'react';
import Layout from '@theme/Layout';
import Playground from '@site/src/components/Playground';

export default function PlaygroundPage(): React.ReactElement {
  return (
    <Layout
      title="Playground"
      description="Try fp-go in your browser with interactive code examples">
      <Playground />
    </Layout>
  );
}

// Made with Bob
