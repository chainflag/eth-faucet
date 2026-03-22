import 'animate.css';
import 'bulma/css/bulma.css';
import { mount } from 'svelte';
import Faucet from './Faucet.svelte';

const app = mount(Faucet, {
  target: document.getElementById('app'),
});

export default app;
