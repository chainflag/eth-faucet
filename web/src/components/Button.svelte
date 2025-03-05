<script>
  export let style = 'primary'
  export let disabled = false
  export let icon
</script>

<button {disabled} class={style} >
  <slot>
{#if icon}
{icon}
{/if}

  </slot>
</button>

<style>
  button {
    color: #fff;
    gap: 16px;
    border-radius: 96px;
    border-width: 0px;
    padding: 12px 24px;
    cursor: pointer;
     background-color: #8950FA;
  }

  .primary {
    background-color: #8950FA;
  }

  .secondary {
    background: #FA50C1;
  }

  button:hover {
    transform: scale(0.98);
    opacity: 0.9;
  }

  button:disabled{
    background-color: #cccccc;
    color: #333;
    cursor:auto;
  }

  button:disabled:hover{
    transform:scale(1);
    opacity: 1;
  }

</style>