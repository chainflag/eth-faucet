<script>
  export let name;
  let address = '0x0000000000000000000000000000000000000000';

  function handleRequest() {
    alert('hello');
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
        This faucet drips 1 Ether every day. Serving from account <a
          href="https://etherscan.io/address/{address}"
          target="_blank"
          ><span class="tag is-warning is-light is-medium">{address}</span></a
        >
      </h2>
    </div>
  </section>

  <div class="container is-fluid">
    <div class="box">
      <div class="block">
        <input
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
