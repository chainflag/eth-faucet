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
      window.scrollTo(
        0,
        document.documentElement.scrollHeight -
          document.documentElement.clientHeight
      );
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
  <body>
    <section class="hero is-info is-fullheight">
      <div class="hero-head">
        <nav class="navbar">
          <div class="container">
            <div class="navbar-brand">
              <a class="navbar-item" href="../">
                <img src="./favicon.png" alt="Logo" />
              </a>
              <span class="navbar-burger burger" data-target="navbarMenu">
                <span />
                <span />
                <span />
              </span>
            </div>
            <div id="navbarMenu" class="navbar-menu">
              <div class="navbar-end">
                <span class="navbar-item">
                  <a class="button is-white is-outlined" href="#">
                    <span class="icon">
                      <i class="fa fa-home" />
                    </span>
                    <span>Home</span>
                  </a>
                </span>
                <span class="navbar-item">
                  <a class="button is-white is-outlined" href="#">
                    <span class="icon">
                      <i class="fa fa-superpowers" />
                    </span>
                    <span>Examples</span>
                  </a>
                </span>
                <span class="navbar-item">
                  <a class="button is-white is-outlined" href="#">
                    <span class="icon">
                      <i class="fa fa-book" />
                    </span>
                    <span>Documentation</span>
                  </a>
                </span>
                <span class="navbar-item">
                  <a
                    class="button is-white is-outlined"
                    href="https://github.com/chainflag/eth-faucet"
                  >
                    <span class="icon">
                      <i class="fa fa-github" />
                    </span>
                    <span>View Source</span>
                  </a>
                </span>
              </div>
            </div>
          </div>
        </nav>
      </div>
      <div class="hero-body">
        <div class="container has-text-centered">
          <div class="column is-6 is-offset-3">
            <h1 class="title">
              Receive {faucetInfo.payout} ETH per request
            </h1>
            <h2 class="subtitle">
              $Serving from account {faucetInfo.account}
            </h2>
            <div class="box">
              <div class="field is-grouped">
                <p class="control is-expanded">
                  <input
                    bind:value={address}
                    class="input"
                    type="text"
                    placeholder="Enter your address"
                  />
                </p>
                <p class="control">
                  <button on:click={handleRequest} class="button is-info"
                    >Request</button
                  >
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  </body>
  <script async type="text/javascript" src="../build/bulma.js"></script>
</main>
