<script>
  export let name;
  let fundAddress = '0x0000000000000000000000000000000000000000';
  let address = null;

  async function handleRequest() {
    const res = await fetch('/api/', {
      method: 'POST',
      body: JSON.stringify({
        address,
      }),
    });

    const json = await res.json();
    const result = JSON.stringify(json);
    alert(result);
  }
</script>

<svelte:head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>{name}</title>
</svelte:head>

<main>
  <nav class="navbar is-link" role="navigation" aria-label="main navigation">
    <div class="navbar-brand">
      <a class="navbar-item" href="/">
        <strong>{name}</strong>
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

<style>
</style>
