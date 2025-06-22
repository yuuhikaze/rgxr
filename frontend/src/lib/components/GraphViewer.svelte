<script lang="ts">
    import { onMount } from 'svelte';

    export let svg: string = '';
    export let showTeX: boolean = false;
    export let tex: string = '';

    let container: HTMLDivElement;

    $: if (container && svg) {
        container.innerHTML = svg;
    }
</script>

<div class="graph-viewer">
    <div class="controls">
        <label>
            <input type="checkbox" bind:checked={showTeX} />
            Show TeX Code
        </label>
    </div>

    {#if showTeX && tex}
        <div class="tex-viewer">
            <h3>TikZ/LaTeX Code</h3>
            <pre>{tex}</pre>
        </div>
    {/if}

    <div class="svg-container" bind:this={container}>
        {#if !svg}
            <p>No graph to display</p>
        {/if}
    </div>
</div>

<style>
    .graph-viewer {
        display: flex;
        flex-direction: column;
        gap: 1rem;
        padding: 1rem;
        border: 1px solid #ddd;
        border-radius: 4px;
    }

    .controls {
        display: flex;
        gap: 1rem;
        align-items: center;
    }

    .tex-viewer {
        background-color: #f5f5f5;
        padding: 1rem;
        border-radius: 4px;
        max-height: 300px;
        overflow-y: auto;
    }

    .tex-viewer pre {
        margin: 0;
        white-space: pre-wrap;
        font-family: 'Courier New', monospace;
        font-size: 0.875rem;
    }

    .svg-container {
        display: flex;
        justify-content: center;
        align-items: center;
        min-height: 300px;
        background-color: #fafafa;
        border-radius: 4px;
        overflow: auto;
    }

    .svg-container :global(svg) {
        max-width: 100%;
        height: auto;
    }
</style>
