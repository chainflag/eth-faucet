<script>
  import { onMount } from 'svelte';
  import { toast, setDefaults as setToast } from 'bulma-toast';
  import faucetIcon from './icons/faucet.svg?raw';
  import githubIcon from './icons/github.svg?raw';

  const ETH_ADDRESS_RE = /^(0x)?[0-9a-fA-F]{40}$/;

  let input = $state('');
  let faucetInfo = $state({
    account: '0x0000000000000000000000000000000000000000',
    network: 'testnet',
    payout: 1,
    symbol: 'ETH',
    hcaptcha_sitekey: '',
  });
  let isLoading = $state(false);
  let hcaptchaLoaded = $state(false);
  let widgetID = $state(null);
  let captchaEl;

  const callbackName = `hcaptchaOnLoad_${Date.now()}`;
  const captchaEnabled = $derived(Boolean(faucetInfo.hcaptcha_sitekey));

  setToast({
    position: 'bottom-center',
    dismissible: true,
    pauseOnHover: true,
    closeOnClick: false,
    animate: { in: 'fadeIn', out: 'fadeOut' },
  });

  $effect(() => {
    document.title = `${faucetInfo.symbol} ${capitalize(faucetInfo.network)} Faucet`;
  });

  $effect(() => {
    if (!hcaptchaLoaded || widgetID !== null) {
      return;
    }
    try {
      widgetID = window.hcaptcha.render(captchaEl, {
        sitekey: faucetInfo.hcaptcha_sitekey,
        size: 'invisible',
      });
    } catch (error) {
      console.error('Failed to render hCaptcha:', error);
    }
  });

  onMount(() => {
    window[callbackName] = () => {
      hcaptchaLoaded = true;
    };

    fetch('/api/info')
      .then((res) => {
        if (!res.ok) throw new Error('Failed to fetch faucet info');
        return res.json();
      })
      .then((info) => {
        faucetInfo = info;
      })
      .catch(() => {
        toast({
          message:
            'Failed to load faucet information. Please refresh the page.',
          type: 'is-danger',
        });
      });

    return () => {
      delete window[callbackName];
    };
  });

  async function handleRequest(e) {
    e.preventDefault();
    if (isLoading) return;

    let address = input.trim();
    if (!address) {
      toast({
        message: 'Please enter an address or ENS name',
        type: 'is-warning',
      });
      return;
    }

    isLoading = true;
    try {
      if (address.endsWith('.eth')) {
        const res = await fetch(`https://api.ensdata.net/${address}`, {
          headers: { Accept: 'application/json' },
        });
        if (!res.ok) throw new Error(`ENS lookup failed (${res.status})`);
        const data = await res.json();
        if (!data?.address || !ETH_ADDRESS_RE.test(data.address)) {
          toast({
            message: 'Invalid ENS name or no address found',
            type: 'is-warning',
          });
          return;
        }
        address = data.address;
      } else if (!ETH_ADDRESS_RE.test(address)) {
        toast({ message: 'Invalid Ethereum address', type: 'is-warning' });
        return;
      } else if (!address.startsWith('0x')) {
        address = '0x' + address;
      }

      const headers = { 'Content-Type': 'application/json' };
      if (hcaptchaLoaded && widgetID !== null) {
        try {
          const { response } = await window.hcaptcha.execute(widgetID, {
            async: true,
          });
          headers['h-captcha-response'] = response;
        } catch {
          toast({
            message: 'Verification failed. Please try again.',
            type: 'is-warning',
          });
          return;
        }
      }

      const res = await fetch('/api/claim', {
        method: 'POST',
        headers,
        body: JSON.stringify({ address }),
      });
      const data = await res.json().catch(() => null);
      if (!res.ok) throw new Error(data?.msg || 'Request failed');
      toast({
        message: data?.msg || 'Transaction successful',
        type: 'is-success',
      });
      input = '';
    } catch (error) {
      toast({
        message: error.message || 'An unexpected error occurred.',
        type: 'is-danger',
      });
    } finally {
      isLoading = false;
    }
  }

  function capitalize(str) {
    if (!str) return '';
    return str.charAt(0).toUpperCase() + str.slice(1).toLowerCase();
  }
</script>

<svelte:head>
  {#if captchaEnabled}
    <script
      src="https://hcaptcha.com/1/api.js?onload={callbackName}&render=explicit"
      async
      defer
    ></script>
  {/if}
</svelte:head>

<main>
  <section class="hero is-info is-fullheight">
    <div class="hero-head">
      <nav class="navbar">
        <div class="container">
          <div class="navbar-brand">
            <a class="navbar-item" href="/">
              <span class="icon">{@html faucetIcon}</span>
              <span><b>{faucetInfo.symbol} Faucet</b></span>
            </a>
          </div>
          <div class="navbar-menu">
            <div class="navbar-end">
              <span class="navbar-item">
                <a
                  class="button is-white is-outlined"
                  href="https://github.com/chainflag/eth-faucet"
                  target="_blank"
                  rel="noopener"
                >
                  <span class="icon">{@html githubIcon}</span>
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
            Receive {faucetInfo.payout}
            {faucetInfo.symbol} per request
          </h1>
          <h2 class="subtitle">
            Serving from {faucetInfo.account}
          </h2>
          <div bind:this={captchaEl} data-size="invisible"></div>
          <div class="box">
            <form class="field is-grouped" onsubmit={handleRequest}>
              <p class="control is-expanded">
                <input
                  bind:value={input}
                  class="input is-rounded"
                  type="text"
                  placeholder="Enter your address or ENS name"
                  autocomplete="off"
                  autocapitalize="off"
                  spellcheck="false"
                />
              </p>
              <p class="control">
                <button
                  type="submit"
                  class="button is-primary is-rounded"
                  disabled={isLoading}
                  class:is-loading={isLoading}
                >
                  Request
                </button>
              </p>
            </form>
          </div>
        </div>
      </div>
    </div>
  </section>
</main>

<style>
  .hero.is-info {
    background:
      linear-gradient(rgba(0, 0, 0, 0.5), rgba(0, 0, 0, 0.5)),
      url('/background.jpg') center / cover no-repeat;
  }
  .hero .subtitle {
    padding: 3rem 0;
    line-height: 1.5;
  }
  .box {
    border-radius: 19px;
  }
</style>
