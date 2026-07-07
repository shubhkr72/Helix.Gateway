import type {ReactNode} from 'react';
import Link from '@docusaurus/Link';
import Layout from '@theme/Layout';
import DayGrid from '@site/src/components/DayGrid';

export default function Home(): ReactNode {
  return (
    <Layout
      title={`Study Days`}
      description="Study Days Documentation">
      <main>
        <DayGrid />
      </main>
    </Layout>
  );
}
