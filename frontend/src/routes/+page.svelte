<script lang="ts">
    import { onMount } from 'svelte';
    import FAList from '$lib/components/FAList.svelte';
    import GraphViewer from '$lib/components/GraphViewer.svelte';
    import OperationsPanel from '$lib/components/OperationsPanel.svelte';
    import { api, type FA, type FARecord } from '$lib/api/client';

    let selectedIds: string[] = [];
    let currentFA: FA | null = null;
    let currentSVG: string = '';
    let currentTeX: string = '';
    let loading = false;

    async function handleFASelect(fa: FARecord) {
        loading = true;
        try {
            const result = await api.convertByUUID(fa.id);
            currentFA = fa.tuple;
            currentSVG = result.svg;
            currentTeX = result.tex;
        } catch (e) {
            console.error('Failed to convert FA:', e);
        } finally {
            loading = false;
        }
    }

    async function handleOperationResult(fa: FA) {
        loading = true;
        try {
            const result = await api.convertFA(fa);
            currentFA = fa;
            currentSVG = result.svg;
            currentTeX = result.tex;

            // Optionally save the result
            // await api.saveFA(fa, `/data/images/${result.id}.svg`, 'Operation result');
        } catch (e) {
            console.error('Failed to convert result:', e);
        } finally {
            loading = false;
        }
    }
</script>

<div class="app">
    <header>
        <h1>RGXR - Regular Expression Visualizer</h1>
    </header>

    <main>
        <div class="sidebar">
            <FAList bind:selectedIds onSelect={handleFASelect} />
            <OperationsPanel {selectedIds} onResult={handleOperationResult} />
        </div>

        <div class="content">
            {#if loading}
                <div class="loading">Converting...</div>
            {:else}
                <GraphViewer svg={currentSVG} tex={currentTeX} />
            {/if}
        </div>
    </main>
</div>

<style>
    .app {
        display: flex;
        flex-direction: column;
        height: 100vh;
    }

    header {
        background-color: #1976d2;
        color: white;
        padding: 1rem;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }

    header h1 {
        margin: 0;
        font-size: 1.5rem;
    }

    main {
        display: flex;
        flex: 1;
        overflow: hidden;
    }

    .sidebar {
        width: 300px;
        padding: 1rem;
        display: flex;
        flex-direction: column;
        gap: 1rem;
        background-color: #f5f5f5;
        overflow-y: auto;
    }

    .content {
        flex: 1;
        padding: 1rem;
        overflow-y: auto;
    }

    .loading {
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100%;
        font-size: 1.25rem;
        color: #666;
    }
</style>
