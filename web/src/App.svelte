<script>
  import 'bulma/css/bulma.min.css';
  import Notifications from 'svelte-notifications';
  import Faucet from './Faucet.svelte';
</script>

<svelte:head>
  <title>ETH Testnet Faucet</title>
</svelte:head>

<Notifications>
  <Faucet />
</Notifications>
