<script>
  import 'bulma/css/bulma.min.css';
  import Notifications from 'svelte-notifications';
  import Faucet from './Faucet.svelte';
</script>

<svelte:head>
  <title>ETH Testnet Faucet</title>
  <link
    href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css"
    rel="stylesheet"
  />
</svelte:head>

<Notifications>
  <Faucet />
</Notifications>
