<script>
  import CopyButton from './CopyButton.svelte';
  import Navigation from './Navigation.svelte';
  import networkIcon from './../assets/networkIcon.svg';
  import arrowRight from './../assets/arrowRight.svg';
  import logo from './../assets/logo.svg';

  export let faucetInfo;
  export let input;
  export let handleRequest;
  export let gweiToEth;

  $: paidCustomer = faucetInfo.paid_customer;

  const openMessageWindow = (subject, email) => {
    const emailSupport = 'support+presto@gateway.fm';
    const mailtoLink = `mailto:${emailSupport}?subject=${encodeURIComponent(subject)}`;

    window.location.href = mailtoLink;
  };

  function autoResize(event) {
    const textarea = event.target;
    textarea.style.height = 'auto'; // Reset height
    textarea.style.height = `${textarea.scrollHeight}px`; // Set new height
  }

  
</script>

<main>
  <section
    class="hero is-info is-fullheight"
    style="background-image: url({faucetInfo.background_url})"
  >
    <div class="hero-head">
      <nav class="navbar">
        <div class="header-container">
          <div class="navbar-brand">
            <a class="navbar-item" href="https://gateway.fm/">
              <span class="icon icon-brand">
                <img src={faucetInfo.logo_url} alt="logo" />
              </span>
            </a>
            {#if !paidCustomer}
              <Navigation />
            {/if}
          </div>
          <div>
            <div class="navbar-end">
              {#if !paidCustomer}
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
              {/if}
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
            <div>{faucetInfo.network}</div>
          </div>
          <div class="title">
            Receive <div class="gas-token">
              {gweiToEth(faucetInfo.payout)}
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
                <textarea
                  bind:value={input}
                  class="input"
                       rows="1"
      on:input={autoResize}
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
                {#if !paidCustomer}
                  <div class="box-offer">
                    Claim 1 {faucetInfo.symbol} test token for development.<br/> If you
                    need additional tokens for extensive testing,<br/> please
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
                {/if}
              </div>
            </div>
          </div>
          {#if !paidCustomer}
            <div class="box-logo">
              Powered by <a
                class="navbar-item"
                href="https://gateway.fm/"
                target="_blank"><img src={logo} alt="Gatewayfm" /></a
              >
            </div>
          {/if}
        </div>
      </div>
    </div>
  </section>
</main>

<style>
.input {
  resize: none;
    overflow: hidden;
}
  .header-container {
    display: flex;
    width: 100%;
    padding-inline: 16px;
    justify-content: space-between;
    align-items: center;
    gap: 16px;
  }
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

  .button.is-primary {
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

  @media (max-width: 768px) {
    .title {
      font-size: 38px;
      line-height: 56px;
    }
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

  .subtitle {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 14px;
    font-weight: 500;
    letter-spacing: 0px;
    color: #161718;
  }

  @media (max-width: 992px) {
    .subtitle {
      font-size: 12px;
    }
  }
  @media (max-width: 768px) {
    .subtitle {
      flex-direction: column;
      font-size: 16px;
      gap: 8px;
    }
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

  .icon-brand {
    width: 5rem;
  }
</style>
