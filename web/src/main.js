import App from './App.svelte';
import 'bulma/css/bulma.min.css';

const app = new App({
  target: document.body,
  props: {
    name: 'ETH Testnet Faucet',
  },
});

export default app;
