<script>
  import { onMount } from 'svelte';
  import { getAddress } from '@ethersproject/address';
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
    try {
      address = getAddress(address);
    } catch (error) {
      addNotification({
        text: error.reason,
        type: 'warning',
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
    let type = res.ok ? 'success' : 'warning';
    addNotification({
      text,
      type,
      removeAfter: 4000,
      position: 'bottom-center',
    });
  }
</script>

<main>
  <section class="hero is-info is-fullheight">
    <div class="hero-head">
      <nav class="navbar">
        <div class="container">
          <div class="navbar-brand">
            <a class="navbar-item" href=".">
              <span class="icon">
                <i class="fa fa-bath" />
              </span>
              <span><b>ETH Faucet</b></span>
            </a>
          </div>
          <div id="navbarMenu" class="navbar-menu">
            <div class="navbar-end">
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
            Serving from {faucetInfo.account}
          </h2>
          <div class="box">
            <div class="field is-grouped">
              <p class="control is-expanded">
                <input
                  bind:value={address}
                  class="input is-rounded"
                  type="text"
                  placeholder="Enter your address"
                />
              </p>
              <p class="control">
                <button
                  on:click={handleRequest}
                  class="button is-primary is-rounded"
                >
                  Request
                </button>
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</main>

<style>
  .hero.is-info {
    background: linear-gradient(rgba(0, 0, 0, 0.5), rgba(0, 0, 0, 0.5)),
      url('https://picsum.photos/1200/900') no-repeat center center fixed;
    -webkit-background-size: cover;
    -moz-background-size: cover;
    -o-background-size: cover;
    background-size: cover;
  }
  .hero .subtitle {
    padding: 3rem 0;
    line-height: 1.5;
  }
  .box {
    border-radius: 19px;
  }
</style>
