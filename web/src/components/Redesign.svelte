<script>
  import CopyButton from './CopyButton.svelte';
  import Navigation from './Navigation.svelte';
  import networkIcon from './../assets/networkIcon.svg';
  import arrowRight from './../assets/arrowRight.svg';
  import logo from './../assets/logo.svg';

  export let faucetInfo;
  export let input;
  export let handleRequest;

  const path = window.location.pathname;
  let network = '';

  network = path.includes('stavanger')
    ? 'Stavanger Testnet'
    : path.substring(1).split('.')[0].replace('-', ' ');

  const openMessageWindow = (subject, email) => {
    const emailSupport = 'support+presto@gateway.fm';
    const mailtoLink = `mailto:${emailSupport}?subject=${encodeURIComponent(subject)}`;

    window.location.href = mailtoLink;
  };

</script>

<main>
  <section class="hero is-info is-fullheight">
    <div class="hero-head">
      <nav class="navbar">
        <div class="container">
          <div class="navbar-brand">
            <a class="navbar-item" href="https://gateway.fm/">
              <img src={logo} alt="logo" />
            </a>

            <Navigation />
          </div>
          <div id="navbarMenu" class="navbar-menu">
            <div class="navbar-end">
              <a
                href="https://presto.gateway.fm/onboarding"
                target="_blank"
                rel="noopener noreferrer"
              >
                <button class="button is-primary is-rounded">
                  Deploy rollup <img
                    src={arrowRight}
                    class="icon arrow-right"
                    alt="arrow"
                  />
                </button></a
              >
            </div>
          </div>
        </div>
      </nav>
    </div>

    <div class="hero-body">
      <div class="container has-text-centered">
        <div class="column is-7 is-offset-3 centered-column">
          <div class="network">
            <img src={networkIcon} alt="logo" />
            <div>{network}</div>
          </div>
          <div class="title">
            Receive <div class="gas-token">
              {faucetInfo.payout}
              {faucetInfo.symbol}
            </div>
          </div>
          <div id="hcaptcha" data-size="invisible"></div>
          <div class="card">
            <div>
              <div class="subtitle">
                <div>Serving from</div>
                <div class="address-from">
                  {faucetInfo.account}
                  <CopyButton text={faucetInfo.account} />
                </div>
              </div>
            </div>
            <div class="field is-grouped">
              <div class="control is-expanded">
                <input
                  bind:value={input}
                  class="input"
                  type="text"
                  placeholder="Enter your address"
                />
              </div>
              <div class="control">
                <button
                  on:click={() => handleRequest(input)}
                  class="button is-primary is-rounded"
                >
                  Request
                </button>

                <div class="box-offer">
                  Claim 1 POL test token for development. If you need additional
                  tokens for extensive testing, please
                  <!-- svelte-ignore a11y-invalid-attribute -->
                  <a
                    class="link"
                    href="#"
                    role="button"
                    on:click={() =>
                      openMessageWindow('Additional tokens request')}
                  >
                    contact support
                  </a>
                </div>
              </div>
            </div>
          </div>
          <div class="box-logo">
            Powered by <a
              class="navbar-item"
              href="https://gateway.fm/"
              target="_blank"><img src={logo} alt="Gatewayfm" /></a
            >
          </div>
        </div>
      </div>
    </div>
  </section>
</main>

<style>
  .hero {
    padding-top: 16px;
  }
  .box-logo {
    display: flex;
    justify-content: center;
    align-items: center;
    color: #303030;
    margin-top: 16px;
    gap: 8px;
    font-size: 14px;
    font-weight: 500;
  }
  .network {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    font-weight: 500;
    padding: 8px 12px;
    background-color: #dcf0fd;
    border-radius: 8px;
    color: #183053;
  }
  .link {
    color: #8950fa !important;
    text-decoration: underline;
    cursor: pointer;
  }
  .button {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    background-color: #8950fa;
    gap: 16px;
  }

  .box-offer {
    font-weight: 400;
    font-size: 12px;
    line-height: 20px;
    letter-spacing: 0px;
  }
  .address-from {
    display: flex;
    align-items: center;
    padding: 8px 12px;
    border-radius: 8px;
    background-color: #eee8ff;
  }
  .field {
    display: flex;
    flex-direction: column;
    gap: 24px;
  }

  .gas-token {
    color: #8950fa;
  }

  .card {
    display: flex;
    flex-direction: column;
    min-width: 100%;
    gap: 24px;
    box-shadow: 0px 8px 16px rgba(0, 0, 0, 0.1);
    border-radius: 19px;
    padding: 32px;
    color: #161718;
  }

  .title {
    display: inline-flex;
    color: #161718;
    gap: 8px;
    font-weight: 500;
    font-size: 72px; /* Adjust size as needed */
    line-height: 80px;
    letter-spacing: 0px;
  }
  .hero.is-info {
    background: url('/background.jpg') no-repeat center center fixed;
    -webkit-background-size: cover;
    -moz-background-size: cover;
    -o-background-size: cover;
    background-size: cover;
  }

  .subtitle {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 14px;
    font-weight: 500;
    /* line-height: 24px; */
    letter-spacing: 0px;
    color: #161718;
  }
  .hero .subtitle {
    line-height: 1.5;
  }
  .centered-column {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .control {
    display: flex;
    flex-direction: column;
    width: 100%;
    gap: 24px;
  }

  .navbar-item {
    color: #8950fa !important;
  }
  .navbar-item:hover {
    background-color: transparent !important;
    cursor: pointer;
  }

  .icon {
    width: 16px;
    height: 16px;
  }
</style>
