import { PanelPlugin } from '@grafana/data';
import { SignatureBanner } from './Signature';

export const plugin = new PanelPlugin(SignatureBanner).setNoPadding();
