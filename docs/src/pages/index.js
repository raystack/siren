import React from 'react';
import Layout from '@theme/Layout';
import clsx from 'clsx';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Container from '../core/Container';
import GridBlock from '../core/GridBlock';
import useBaseUrl from '@docusaurus/useBaseUrl';

const Hero = () => {
  const { siteConfig } = useDocusaurusContext();
  return (
    <div className="homeHero">
      <div className="logo"><img src={useBaseUrl('img/pattern.svg')} /></div>
      <div className="container banner">
        <div className="row">
          <div className={clsx('col col--5')}>
            <div className="homeTitle">{siteConfig.tagline}</div>
            <small className="homeSubTitle">Siren provides an easy-to-use universal alert, notification, channels management framework for the entire observability infrastructure.</small>
            <a className="button" href="docs/introduction">Documentation</a>
          </div>
          <div className={clsx('col col--1')}></div>
          <div className={clsx('col col--6')}>
            <div className="text--right"><img src={useBaseUrl('img/banner.svg')} /></div>
          </div>
        </div>
      </div>
    </div >
  );
};

export default function Home() {
  const { siteConfig } = useDocusaurusContext();
  return (
    <Layout
      title={siteConfig.tagline}
      description="Siren provides an easy-to-use universal alert, notification, channels management framework for the entire observability infrastructure.">
      <Hero />
      <main>
        <Container className="textSection wrapper" background="light">
          <h1>Built for scale</h1>
          <p>
            Siren provides an easy-to-use universal alert, notification, channels management framework for the entire observability infrastructure..
          </p>
          <GridBlock
            layout="threeColumn"
            contents={[
              {
                title: 'Rule Templates',
                content: (
                  <div>
                    Siren provides a way to define templates over prometheus Rule, which can be reused to create multiple instances of same rule with configurable thresholds.
                  </div>
                ),
              },
              {
                title: 'Multi-tenancy',
                content: (
                  <div>
                    Rules created with Siren are by default multi-tenancy aware.
                  </div>
                ),
              },
              {
                title: 'DIY Interface',
                content: (
                  <div>
                    Siren can be used to easily create/edit prometheus rules. It also provides soft delete(disable) so that you can preserve thresholds in case you need to reuse the same alert.
                  </div>
                ),
              },
              {
                title: 'Managing bulk rules',
                content: (
                  <div>
                    Siren enables users to manage bulk alerts using YAML files in specified format using simple CLI.
                  </div>
                ),
              },
              {
                title: 'Credentials Management',
                content: (
                  <div>
                    Siren can store slack and pagerduty credentials, sync them with Cortex alertmanager to deliver alerts on proper channels, in a multi-tenant fashion. It gives a simple interface to rotate the credentials on demand via HTTP API.
                  </div>
                ),
              },
              {
                title: 'Alert History',
                content: (
                  <div>
                    Siren can store alerts triggered via Cortex Alertmanager, which can be used for audit purposes.
                  </div>
                ),
              },
            ]}
          />
        </Container>

        <Container className="textSection wrapper" background="light">
          <h1>Trusted by</h1>
          <p>
            Siren was originally created for the Gojek data processing platform,
            and it has been used, adapted and improved by other teams internally and externally.
          </p>
          <GridBlock className="logos"
            layout="fourColumn"
            contents={[
              {
                content: (
                  <img src={useBaseUrl('users/gojek.png')} />
                ),
              },
              {
                content: (
                  <img src={useBaseUrl('users/midtrans.png')} />
                ),
              },
              {
                content: (
                  <img src={useBaseUrl('users/mapan.png')} />
                ),
              },
              {
                content: (
                  <img src={useBaseUrl('users/moka.png')} />
                ),
              },
            ]}>
          </GridBlock>
        </Container>
      </main>
    </Layout >
  );
}
