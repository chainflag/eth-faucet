<script>
  import 'bulma/css/bulma.min.css';
  import Notifications from 'svelte-notifications';
  import Faucet from './Faucet.svelte';
</script>

<svelte:head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>ETH Testnet Faucet</title>
</svelte:head>

<Notifications>
  <Faucet />
</Notifications>
