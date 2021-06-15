<script>
  import { getNotificationsContext } from 'svelte-notifications';

  const { addNotification } = getNotificationsContext();

  let fundAddress = '0x04daa5C20d3278Ce47241805b1572d4a6ab95Db3';
  let address = null;
  async function handleRequest() {
    const res = await fetch('/api/', {
      method: 'POST',
      body: JSON.stringify({
        address,
      }),
    });

    let text = res.status === 429 ? res.statusText : await res.text();
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
      <h1 class="title">Get Testnet Ether</h1>
      <h2 class="subtitle">
        This faucet drips 1 Ether every day. Serving from account <span
          class="tag is-warning is-light is-medium">{fundAddress}</span
        >
      </h2>
    </div>
  </section>

  <div class="container is-fluid">
    <div class="box">
      <div class="block">
        <input
          bind:value={address}
          class="input is-dark"
          type="text"
          placeholder="Enter your account address"
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
