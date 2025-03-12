<script>
  import copyIcon from '../assets/copyIcon.svg';
  export let text = 'Copy';
  let tooltipVisible = false;

  function copyToClipboard() {
    navigator.clipboard.writeText(text);
    tooltipVisible = true;
    setTimeout(() => tooltipVisible = false, 1000); // Hide tooltip after 1s
  }
</script>

<button class="button-copy" on:click={copyToClipboard}>
  <img src={copyIcon} alt={text} />
  {#if tooltipVisible}
    <span class="tooltip">Copied!</span>
  {/if}
</button>

<style>
  .button-copy {
    position: relative;
    background: none;
    border: none;
    cursor: pointer;
  }

  img {
    width: 16px;
    height: 16px;
  }

  .tooltip {
    position: absolute;
    bottom: 25px;
    left: 50%;
    transform: translateX(-50%);
    background: rgba(0, 0, 0, 0.5);
    color: white;
    padding: 5px 8px;
    border-radius: 4px;
    font-size: 12px;
    white-space: nowrap;
  }
</style>
