<script>
  const navigationList = [
    { title: 'Rollup', url: 'https://gateway.fm/presto' },
    { title: 'Stakeway', url: 'https://stakeway.com/' },
    { title: 'RPC', url: 'https://gateway.fm/rpc' },
    { title: 'Blog', url: 'https://gateway.fm/blog' },
    { title: 'About', url: 'https://gateway.fm/about' },
    { title: 'Careers', url: 'https://boards.eu.greenhouse.io/gatewayfm' },
  ];
</script>

<div>
  <nav class="navbar">
    {#each navigationList as { title, url }}
      <a
        class="navbar-item"
        href={url}
        target="_blank"
        rel="noopener noreferrer"
      >
        {title}
      </a>
    {/each}
  </nav>
</div>

<style>
  .navbar {
    display: flex;
    gap: 12px;
  }

  .navbar-item {
    text-decoration: none;
    color: #676e73 !important;
    transition: background 0.3s;
  }

  .navbar-item:hover {
    color: #8950fa !important;
    background-color: transparent !important;
  }
</style>
