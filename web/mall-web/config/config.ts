import { defineConfig } from '@umijs/max';
import proxy from './proxy';
import routes from './routes';
const environment = process.env.UMI_ENV === 'test' ? 'test' : 'dev';
export default defineConfig({ routes, proxy: proxy[environment], hash: true, title: 'Mall 商城', esbuildMinifyIIFE: true });
