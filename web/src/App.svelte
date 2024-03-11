<script>
  import 'animate.css';
  import 'bulma/css/bulma.css';
  import Faucet from './Faucet.svelte';

	const faviconPath = import.meta.env.VITE_FAVICON_PATH
</script>

<svelte:head>
  <link
    href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css"
    rel="stylesheet"
  />
	<link rel="icon" type="image/png" href="{faviconPath}" />
</svelte:head>

<Faucet />
