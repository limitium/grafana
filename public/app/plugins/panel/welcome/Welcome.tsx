import React, { FC } from 'react';
import { css } from 'emotion';
import { GrafanaTheme } from '@grafana/data';
import { stylesFactory, useTheme } from '@grafana/ui';
import lightBackground from './img/background_light.svg';
import { config } from '@grafana/runtime';

const userSignature = config.bootData.user.signature;

export const WelcomeBanner: FC = () => {
  const styles = getStyles(useTheme());

  return (
    <div className={styles.container}>
      <h1 className={styles.title}>Welcome to bloomy.space</h1>
      <h2 className={styles.nextStep}>Next steps to start</h2>
      <ul className={styles.helpText}>
        <li>
          Download android application to transfer measuremetns from sensors{' '}
          <a className={styles.a} href="https://play.google.com/store/apps/details?id=space.bloomy.ward">
            play.google
          </a>
        </li>
        <li>
          Enter your <b>{userSignature}</b> signature in the app.
        </li>
        <li>
          Create a new dashboard for your data.{' '}
          <a className={styles.a} href="https://grafana.com/docs/grafana/latest/dashboards/">
            how to?
          </a>
        </li>
        <li>
          Add alerts for important metrics.{' '}
          <a className={styles.a} href="https://grafana.com/docs/grafana/latest/alerting/">
            how to?
          </a>
        </li>
        <li>
          Join{' '}
          <a className={styles.a} href="https://t.me/bloomy_space_chat">
            telegram
          </a>{' '}
          or{' '}
          <a className={styles.a} href="https://discord.gg/xC3YNRDk">
            discord
          </a>{' '}
          chat to share your ideas.
        </li>
      </ul>
      <h3 className={styles.recommended}>Recommended setup</h3>
      <ul className={styles.helpText}>
        <li>
          <a className={styles.a} href="https://s.click.aliexpress.com/e/_AquuvN">
            CGG1
          </a>{' '}
          - Air humidity and temperature
        </li>
        <li>
          <a className={styles.a} href="https://s.click.aliexpress.com/e/_Afi1vD">
            HHCCJCY01
          </a>{' '}
          - Ground moisture, conductivity, temperature and illumination
        </li>
        <li>
          <a className={styles.a} href="https://github.com/limitium/feeler">
            WRDSNR01
          </a>{' '}
          - CO2 sensor
        </li>
      </ul>
      <p className={styles.list}>
        List of all{' '}
        <a className={styles.a} href="https://github.com/limitium/wand/blob/master/README.md#supported-sensors">
          supported sensors
        </a>
      </p>
    </div>
  );
};

const getStyles = stylesFactory((theme: GrafanaTheme) => {
  const backgroundImage = theme.isDark ? 'public/img/login_background_dark.svg' : lightBackground;

  return {
    container: css`
      display: flex;
      background: url(${backgroundImage}) no-repeat;
      background-size: cover;
      height: 100%;
      align-items: center;
      justify-content: space-between;
      padding: ${theme.spacing.lg};
      flex-direction: column;
      align-items: flex-start;
      justify-content: center;

      @media only screen and (max-width: ${theme.breakpoints.lg}) {
        background-position: 0px;
        flex-direction: column;
        align-items: flex-start;
        justify-content: center;
      }

      @media only screen and (max-width: ${theme.breakpoints.sm}) {
        padding: ${theme.spacing.sm};
      }
    `,
    title: css`
      margin-bottom: 0;
      @media only screen and (max-width: ${theme.breakpoints.lg}) {
        margin-bottom: ${theme.spacing.sm};
        padding-bottom: ${theme.spacing.lg};
      }

      @media only screen and (max-width: ${theme.breakpoints.md}) {
        font-size: ${theme.typography.heading.h2};
        padding-bottom: ${theme.spacing.sm};
      }
      @media only screen and (max-width: ${theme.breakpoints.sm}) {
        font-size: ${theme.typography.heading.h3};
        padding-bottom: ${theme.spacing.sm};
      }
    `,
    nextStep: css`
      padding-top: ${theme.spacing.lg};
    `,
    a: css`
      text-decoration: underline;
    `,
    list: css`
      padding-top: ${theme.spacing.xl};
    `,
    help: css`
      display: flex;
      align-items: baseline;
    `,
    helpText: css`
      font-size: ${theme.typography.heading.h4};
      padding-left: ${theme.spacing.xl};

      @media only screen and (max-width: ${theme.breakpoints.md}) {
        font-size: ${theme.typography.heading.h4};
        padding-left: ${theme.spacing.lg};
      }

      @media only screen and (max-width: ${theme.breakpoints.sm}) {
        display: none;
        padding-left: ${theme.spacing.sm};
      }
    `,
    recommended: css`
      padding-top: ${theme.spacing.lg};
    `,
    helpLinks: css`
      display: flex;
      flex-wrap: wrap;
    `,
    helpLink: css`
      margin-right: ${theme.spacing.md};
      text-decoration: underline;
      text-wrap: no-wrap;

      @media only screen and (max-width: ${theme.breakpoints.sm}) {
        margin-right: 8px;
      }
    `,
  };
});
