<script>
  import { onMount } from 'svelte';
  import { getAddress } from '@ethersproject/address';
  import { CloudflareProvider } from '@ethersproject/providers';
  import { setDefaults as setToast, toast } from 'bulma-toast';

  let input = null;
  let faucetInfo = {
    account: '0x0000000000000000000000000000000000000000',
    network: 'testnet',
    payout: 1,
    symbol: 'ETH',
    hcaptcha_sitekey: '',
  };

  let mounted = false;
  let hcaptchaLoaded = false;

  onMount(async () => {
    const res = await fetch('/api/info');
    faucetInfo = await res.json();
    mounted = true;
  });

  window.hcaptchaOnLoad = () => {
    hcaptchaLoaded = true;
  };

  $: document.title = `${faucetInfo.symbol} ${capitalize(
    faucetInfo.network,
  )} Faucet`;

  let widgetID;
  $: if (mounted && hcaptchaLoaded) {
    widgetID = window.hcaptcha.render('hcaptcha', {
      sitekey: faucetInfo.hcaptcha_sitekey,
    });
  }

  setToast({
    position: 'bottom-center',
    dismissible: true,
    pauseOnHover: true,
    closeOnClick: false,
    animate: { in: 'fadeIn', out: 'fadeOut' },
  });

  async function handleRequest() {
    let address = input;
    if (address === null) {
      toast({ message: 'input required', type: 'is-warning' });
      return;
    }

    if (address.endsWith('.eth')) {
      try {
        const provider = new CloudflareProvider();
        address = await provider.resolveName(address);
        if (!address) {
          toast({ message: 'invalid ENS name', type: 'is-warning' });
          return;
        }
      } catch (error) {
        toast({ message: error.reason, type: 'is-warning' });
        return;
      }
    }

    try {
      address = getAddress(address);
    } catch (error) {
      toast({ message: error.reason, type: 'is-warning' });
      return;
    }

    try {
      let headers = {
        'Content-Type': 'application/json',
      };

      if (hcaptchaLoaded) {
        const { response } = await window.hcaptcha.execute(widgetID, {
          async: true,
        });
        headers['h-captcha-response'] = response;
      }

      const res = await fetch('/api/claim', {
        method: 'POST',
        headers,
        body: JSON.stringify({
          address,
        }),
      });

      let { msg } = await res.json();
      let type = res.ok ? 'is-success' : 'is-warning';
      toast({ message: msg, type });
    } catch (err) {
      console.error(err);
    }
  }

  function capitalize(str) {
    const lower = str.toLowerCase();
    return str.charAt(0).toUpperCase() + lower.slice(1);
  }
</script>

<svelte:head>
  {#if mounted && faucetInfo.hcaptcha_sitekey}
    <script
      src="https://hcaptcha.com/1/api.js?onload=hcaptchaOnLoad&render=explicit"
      async
      defer
    ></script>
  {/if}
</svelte:head>

<!-- Site header -->
<header class="absolute w-full z-30">
  <div class="max-w-6xl mx-auto px-4 sm:px-6">
    <div class="flex items-center justify-between h-16 md:h-20">

      <!-- Site branding -->
      <div class="flex-1">
        <!-- Logo -->
        <a class="inline-flex items-center" href="index.html" aria-label="Cruip">
          <img class="max-w-none" src="/images/stratis_logo_white.svg" width="38" height="38" alt="Stellar">
          <span class="ml-3 hidden md:block">Stratis Auroria Faucet</span>
        </a>
      </div>

      <!-- Desktop sign in links -->
      <ul class="flex-1 flex justify-end items-center">
        <li class="ml-6">
          <a class="btn-sm text-slate-300 hover:text-white transition duration-150 ease-in-out w-full group [background:linear-gradient(theme(colors.slate.900),_theme(colors.slate.900))_padding-box,_conic-gradient(theme(colors.slate.400),_theme(colors.slate.700)_25%,_theme(colors.slate.700)_75%,_theme(colors.slate.400)_100%)_border-box] relative before:absolute before:inset-0 before:bg-slate-800/30 before:rounded-full before:pointer-events-none"
            href="https://github.com/stratisproject/strax-faucet">
            <span class="relative inline-flex items-center">
              View Source <span
                class="tracking-normal text-purple-500 group-hover:translate-x-0.5 transition-transform duration-150 ease-in-out ml-1">-&gt;</span>
            </span>
          </a>
        </li>
      </ul>
    </div>
  </div>
</header>

<!-- Page content -->
<main class="grow">
  <!-- Features #2 -->
  <section class="relative">

    <!-- Particles animation -->
    <div class="absolute left-1/2 -translate-x-1/2 top-0 -z-10 w-80 h-80 -mt-24 -ml-32">
      <div class="absolute inset-0 -z-10" aria-hidden="true">
        <canvas data-particle-animation data-particle-quantity="6" data-particle-staticity="30"></canvas>
      </div>
    </div>

    <div class="max-w-6xl mx-auto px-4 sm:px-6">
      <div class="pt-16 md:pt-32">

        <!-- Particles animation -->
      <div class="absolute inset-0 -z-10" aria-hidden="true">
        <canvas data-particle-animation></canvas>
      </div>

      <!-- Illustration -->
      <div class="absolute inset-0 -z-10 -mx-28 rounded-b-[3rem] pointer-events-none overflow-hidden"
        aria-hidden="true">
        <div class="absolute left-1/2 -translate-x-1/2 bottom-0 -z-10">
          <img src="/images/glow-bottom.svg" class="max-w-none" width="2146" height="774" alt="Hero Illustration">
        </div>
      </div>

        <!-- Section header -->
        <div class="max-w-xl mx-auto text-center md:px-1 px-5 pb-20 md:pb-20">
          <h2
            class="h2 bg-clip-text text-transparent bg-gradient-to-r from-slate-200/60 via-slate-200 to-slate-200/60 pb-4">
            Receive {faucetInfo.payout} {faucetInfo.symbol} per request.</h2>
          <p class="md:text-lg text-sm text-slate-400">Serving from {faucetInfo.account}</p>
          <div id="hcaptcha" data-size="invisible"></div>
          <div class="py-12 px-12 w-xl mx-auto">
            <div class="mb-40">
              <div class="space-y-2">
                <div>
                  <label class="block text-sm text-slate-300 font-medium mb-1" for="address">tSTRAX Address</label>
                  <input id="address" class="form-input w-full h-3 p-3 rounded-full" type="text" bind:value={input} required placeholder="Enter your address or ENS name" />
                </div>
              </div>
              <div class="mt-2">
                <button on:click={handleRequest} type="button"
                  class="btn text-slate-300 hover:text-white transition duration-150 ease-in-out w-full group [background:linear-gradient(theme(colors.slate.900),_theme(colors.slate.900))_padding-box,_conic-gradient(theme(colors.slate.400),_theme(colors.slate.700)_25%,_theme(colors.slate.700)_75%,_theme(colors.slate.400)_100%)_border-box] relative before:absolute before:inset-0 before:bg-slate-800/30 before:rounded-full before:pointer-events-none">
                  Request <span
                    class="tracking-normal text-purple-300 group-hover:translate-x-0.5 transition-transform duration-150 ease-in-out ml-1">-&gt;</span>
                </button>
              </div>
              <p>&nbsp;</p>
              <p>&nbsp;</p>
              <p>&nbsp;</p>
              <p>&nbsp;</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</main>

<!-- Site footer -->
<footer>
  <div class="max-w-6xl mx-auto px-4 sm:px-6">

    <!-- Blocks -->
    <div class="grid sm:grid-cols-12 gap-8 py-8 md:py-4">

      <!-- 1st block -->
      <div class="sm:col-span-12 lg:col-span-4 order-1 lg:order-none">
        <div class="h-full flex flex-col sm:flex-row lg:flex-col justify-between">
          <div class="mb-4 sm:mb-0">
            <div class="mb-4">
              <!-- Logo -->
              <a class="inline-flex" href="index.html" aria-label="Cruip">
                <img src="/images/stratis_logo_white.svg" width="38" height="38" alt="Stellar">
              </a>
            </div>
            <div class="text-sm text-slate-300">&copy; Stratis Platform <span class="text-slate-500">-</span> All
              rights
              reserved.</div>
          </div>
          <!-- Social links -->
          <ul class="flex">
            <li>
              <a class="flex justify-center items-center text-purple-500 hover:text-purple-400 transition duration-150 ease-in-out"
                href="#0" aria-label="Twitter">
                <svg class="w-8 h-8 fill-current" viewBox="0 0 32 32" xmlns="http://www.w3.org/2000/svg">
                  <path
                    d="m13.063 9 3.495 4.475L20.601 9h2.454l-5.359 5.931L24 23h-4.938l-3.866-4.893L10.771 23H8.316l5.735-6.342L8 9h5.063Zm-.74 1.347h-1.457l8.875 11.232h1.36l-8.778-11.232Z" />
                </svg>
              </a>
            </li>
            <li class="ml-2">
              <a class="flex justify-center items-center text-purple-500 hover:text-purple-400 transition duration-150 ease-in-out"
                href="#0" aria-label="Dev.to">
                <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32">
                  <path class="w-8 h-8 fill-current"
                    d="M12.29 14.3a.69.69 0 0 0-.416-.155h-.623v3.727h.623a.689.689 0 0 0 .416-.156.543.543 0 0 0 .21-.466v-2.488a.547.547 0 0 0-.21-.462ZM22.432 8H9.568C8.704 8 8.002 8.7 8 9.564v12.872A1.568 1.568 0 0 0 9.568 24h12.864c.864 0 1.566-.7 1.568-1.564V9.564A1.568 1.568 0 0 0 22.432 8Zm-8.925 9.257a1.631 1.631 0 0 1-1.727 1.687h-1.657v-5.909h1.692a1.631 1.631 0 0 1 1.692 1.689v2.533ZM17.1 14.09h-1.9v1.372h1.163v1.057H15.2v1.371h1.9v1.056h-2.217a.72.72 0 0 1-.74-.7v-4.471a.721.721 0 0 1 .7-.739H17.1v1.054Zm3.7 4.118c-.471 1.1-1.316.88-1.694 0l-1.372-5.172H18.9l1.058 4.064 1.056-4.062h1.164l-1.378 5.17Z" />
                </svg>
              </a>
            </li>
            <li class="ml-2">
              <a class="flex justify-center items-center text-purple-500 hover:text-purple-400 transition duration-150 ease-in-out"
                href="#0" aria-label="Github">
                <svg class="w-8 h-8 fill-current" viewBox="0 0 32 32" xmlns="http://www.w3.org/2000/svg">
                  <path
                    d="M16 8.2c-4.4 0-8 3.6-8 8 0 3.5 2.3 6.5 5.5 7.6.4.1.5-.2.5-.4V22c-2.2.5-2.7-1-2.7-1-.4-.9-.9-1.2-.9-1.2-.7-.5.1-.5.1-.5.8.1 1.2.8 1.2.8.7 1.3 1.9.9 2.3.7.1-.5.3-.9.5-1.1-1.8-.2-3.6-.9-3.6-4 0-.9.3-1.6.8-2.1-.1-.2-.4-1 .1-2.1 0 0 .7-.2 2.2.8.6-.2 1.3-.3 2-.3s1.4.1 2 .3c1.5-1 2.2-.8 2.2-.8.4 1.1.2 1.9.1 2.1.5.6.8 1.3.8 2.1 0 3.1-1.9 3.7-3.7 3.9.3.4.6.9.6 1.6v2.2c0 .2.1.5.6.4 3.2-1.1 5.5-4.1 5.5-7.6-.1-4.4-3.7-8-8.1-8z" />
                </svg>
              </a>
            </li>
          </ul>
        </div>
      </div>

    </div>

  </div>
</footer>

<!--<main>
  <section class="hero is-info is-fullheight">
    <div class="hero-head bg-purple-900">
      <nav class="navbar">
        <div class="container">
          <div class="navbar-brand">
            <a class="navbar-item" href="../..">
              <img src="stratis_logo_white.svg" class="mr-2" width="40" height="40 /" />
              <span>Stratis Auroria Faucet</span>
            </a>
          </div>
          <div id="navbarMenu" class="navbar-menu">
            <div class="navbar-end">
              <span class="navbar-item">
                <a
                  class="button is-white is-outlined"
                  href="https://github.com/stratisproject/strax-faucet"
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
          <h1 class="text-purple-900 text-3xl mb-5">
            Receive {faucetInfo.payout}
            {faucetInfo.symbol} per request
          </h1>
          <h2 class="text-xl text-purple-700 mb-5">
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
                  placeholder="Enter your address or ENS name"
                />
              </p>
              <p class="control">
                <button
                  on:click={handleRequest}
                  class="button bg-purple-900 text-white hover:text-white hover:bg-purple-600 is-rounded"
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
-->

