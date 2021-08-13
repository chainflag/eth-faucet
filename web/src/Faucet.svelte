<script>
  import { onMount } from 'svelte';
  import { isAddress } from '@ethersproject/address';
  import { formatEther } from '@ethersproject/units';
  import { getNotificationsContext } from 'svelte-notifications';

  const { addNotification } = getNotificationsContext();

  let address = null;
  let faucetInfo = {
    account: '0x0000000000000000000000000000000000000000',
    payout: 1,
  };

  onMount(async () => {
    const res = await fetch('/api/info');
    faucetInfo = await res.json();
    faucetInfo.payout = parseInt(formatEther(faucetInfo.payout));
  });

  async function handleRequest() {
    if (!isAddress(address)) {
      addNotification({
        text: 'Invalid address',
        type: 'danger',
        removeAfter: 4000,
        position: 'bottom-center',
      });
      return;
    }

    const res = await fetch('/api/claim', {
      method: 'POST',
      body: JSON.stringify({
        address,
      }),
    });
    let text = await res.text();
    let type = res.ok ? 'success' : 'danger';
    addNotification({
      text,
      type,
      removeAfter: 4000,
      position: 'bottom-center',
    });
  }
</script>

<main>
  <nav class="navbar is-link" role="navigation" aria-label="main navigation">
    <div class="navbar-brand">
      <a class="navbar-item" href="/">
        <strong>Ether Faucet</strong>
      </a>
    </div>

    <div class="navbar-menu">
      <div class="navbar-end">
        <div class="navbar-item">
          <div class="button">
            <a
              class="button is-text is-small"
              href="https://github.com/chainflag/eth-faucet"
              target="_blank"
            >
              <strong>View on Github</strong>
            </a>
          </div>
        </div>
      </div>
    </div>
  </nav>

  <section class="section">
    <div class="container">
      <h1 class="title">Receive {faucetInfo.payout} ETH per request</h1>
      <h2 class="subtitle">
        Serving from account
        <span class="tag is-warning is-light is-medium"
          >{faucetInfo.account}</span
        >
      </h2>
    </div>
  </section>

  <div class="container is-fluid">
    <div class="box">
      <div class="block">
        <label class="label">Enter your account address</label>
        <input
          bind:value={address}
          class="input"
          type="text"
          placeholder="0x..."
        />
      </div>
      <button on:click={handleRequest} class="button is-danger">Request</button>
    </div>
  </div>

  <footer class="footer">
    <div class="content has-text-centered">
      <p>
        Powered by <a href="https://chainflag.org" target="_blank">ChainFlag</a>
      </p>
    </div>
  </footer>
</main>
