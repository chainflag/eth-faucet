<script>
  let isOpen = false;

  const navigationList = [
    { title: 'Rollup', url: 'https://gateway.fm/presto' },
    { title: 'Stakeway', url: 'https://stakeway.com/' },
    { title: 'RPC', url: 'https://gateway.fm/rpc' },
    { title: 'Blog', url: 'https://gateway.fm/blog' },
    { title: 'About', url: 'https://gateway.fm/about' },
    { title: 'Careers', url: 'https://boards.eu.greenhouse.io/gatewayfm' },
  ];
</script>

<div class="navbar-container">
  <button class="burger-menu" on:click={() => (isOpen = !isOpen)}>
    ☰
  </button>

  <nav class="navbar" class:open={isOpen}>
    {#each navigationList as { title, url }}
      <a class="navbar-item" href={url} target="_blank" rel="noopener noreferrer">
        {title}
      </a>
    {/each}
  </nav>
</div>

<style>
  .navbar-container {
    position: relative;
  }

  .burger-menu {
    display: none;
    font-size: 34px;
    background: none;
    border: none;
    cursor: pointer;
    padding: 10px;
  }

  .navbar {
    display: flex;
    gap: 12px;
  }

  .navbar-item {
    text-decoration: none;
    color: #676e73 !important;
    transition: color 0.3s;
  }

  .navbar-item:hover {
    color: #8950fa !important;
  }

  @media (max-width: 768px) {
    .burger-menu {
      display: block;
      color: #8950fa;
    }

    .navbar {
      display: none;
      flex-direction: column;
      position: absolute;
      top: 50px;
      left: 0;
      background: white;
      /* min-width: 100%; */
      padding: 10px;
      box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    }

    .navbar.open {
      display: flex;
    }
  }
</style>
