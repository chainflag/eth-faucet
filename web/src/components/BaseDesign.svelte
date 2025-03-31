<script>
  export let faucetInfo;
  export let input;
  export let handleRequest;
  export let gweiToEth;

</script>

<main>
  <section class="hero is-info is-fullheight" style='background-image: url({faucetInfo.background_url})'>
    <div class="hero-head">
      <nav class="navbar">
        <div class="container">
          <div class="navbar-brand">
            <a class="navbar-item" href="../..">
              <span class="icon icon-brand">
                <img src={faucetInfo.logo_url} alt="logo"/>
              </span>
              <span><b>{faucetInfo.symbol} Faucet</b></span>
            </a>
          </div>
        </div>
      </nav>
    </div>

    <div class="hero-body">
      <div class="container has-text-centered">
        <div class="column is-6 is-offset-3">
          <h1 class="title">
            Receive {gweiToEth(faucetInfo.payout)}
            {faucetInfo.symbol} per request
          </h1>
          <h2 class="subtitle">
            Serving from {faucetInfo.account}
          </h2>
          <div id="hcaptcha" data-size="invisible"></div>
          <div class="box">
            <div class="field is-grouped">
              <p class="control is-expanded">
                <input
                  bind:value={input}
                  class="input is-rounded"
                  type="text"
                  placeholder="Enter your address"
                />
              </p>
              <p class="control">
                <button
                  on:click={() => handleRequest(input)}
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

  .navbar-item {
    display: flex;
    flex-direction: row;
    align-items: center;
    gap: 8px;
  }
  .icon-brand {
    width: 6rem
  }
 .hero.is-info {
    background:
      linear-gradient(rgba(0, 0, 0, 0.2), rgba(0, 0, 0, 0.2)),
      no-repeat center center fixed;
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
