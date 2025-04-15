<script>
  // @ts-nocheck

  import { onMount } from 'svelte';
  import { getAddress } from '@ethersproject/address';
  import { CloudflareProvider } from '@ethersproject/providers';
  import { setDefaults as setToast, toast } from 'bulma-toast';

  import BaseDesign from './components/BaseDesign.svelte';
  import Redesign from './components/Redesign.svelte';

  let input = null;

  let faucetInfo = {
    account: '0x0000000000000000000000000000000000000000',
    network: 'testnet',
    payout: 1000000000,
    symbol: 'ETH',
    hcaptcha_sitekey: '',
    logo_url: '/gatewayfm-logo.svg',
    background_url: 'background.jpg',
    frontend_type: 'redesign',
    paid_customer: false,
  };

  let mounted = false;
  let hcaptchaLoaded = false;

  function gweiToEth(gwei) {
    let str = gwei.toString();
    let len = str.length;

    // Add leading zeros if necessary
    if (len <= 9) {
      str = '0'.repeat(9 - len) + str;
      len = str.length;
    }

    // Insert decimal point
    str = str.slice(0, len - 9) + '.' + str.slice(len - 9);

    // Add leading zero if necessary
    if (str.startsWith('.')) {
      str = '0' + str;
    }

    // Remove trailing zeros
    str = str.replace(/0+$/, '');

    // Remove trailing decimal point
    if (str.endsWith('.')) {
      str = str.slice(0, -1);
    }

    return str;
  }

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

  window.hcaptchaOnLoad = () => {
    hcaptchaLoaded = true;
  };

  $: baseFrontendType = faucetInfo.frontend_type === 'base';
  $: redesignFrontendType = faucetInfo.frontend_type === 'redesign';

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
    message: '',
    position: 'bottom-center',
    dismissible: true,
    pauseOnHover: true,
    closeOnClick: false,
    animate: { in: 'fadeIn', out: 'fadeOut' },
  });

  async function handleRequest(input) {
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

{#if baseFrontendType}
  <BaseDesign {faucetInfo} {input} {handleRequest} {gweiToEth} />
{:else if redesignFrontendType}
  <Redesign {faucetInfo} {input} {handleRequest} {gweiToEth} />
{/if}
