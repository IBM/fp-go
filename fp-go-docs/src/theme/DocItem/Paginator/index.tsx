import React, {type ReactNode} from 'react';
import {useDoc} from '@docusaurus/plugin-content-docs/client';
import {Pager} from '@site/src/components/content';

export default function DocItemPaginator(): ReactNode {
  const {metadata} = useDoc();
  const {previous, next} = metadata;
  if (!previous && !next) return null;
  return (
    <Pager
      prev={previous ? {to: previous.permalink, title: previous.title} : undefined}
      next={next ? {to: next.permalink, title: next.title} : undefined}
    />
  );
}
