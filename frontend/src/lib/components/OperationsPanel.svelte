<script lang="ts">
    import { api } from '$lib/api/client';

    export let selectedIds: string[] = [];
    export let onResult: (result: any) => void;

    let loading = false;
    let error: string | null = null;

    async function performOperation(operation: string) {
        if (selectedIds.length === 0) {
            error = 'Please select at least one FA';
            return;
        }

        loading = true;
        error = null;

        try {
            let result;

            switch (operation) {
                case 'union':
                    if (selectedIds.length < 2) {
                        throw new Error('Union requires at least 2 FAs');
                    }
                    result = await api.union(selectedIds);
                    break;

                case 'intersection':
                    if (selectedIds.length < 2) {
                        throw new Error('Intersection requires at least 2 FAs');
                    }
                    const response = await fetch('/api/intersection', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ uuids: selectedIds })
                    });
                    result = await response.json();
                    break;

                case 'complement':
                    if (selectedIds.length !== 1) {
                        throw new Error('Complement requires exactly 1 FA');
                    }
                    const compResponse = await fetch(`/api/complement?uuid=${selectedIds[0]}`);
                    result = await compResponse.json();
                    break;

                case 'minimize':
                    if (selectedIds.length !== 1) {
                        throw new Error('Minimize requires exactly 1 FA');
                    }
                    const minResponse = await fetch(`/api/minimize-dfa?uuid=${selectedIds[0]}`);
                    result = await minResponse.json();
                    break;

                default:
                    throw new Error(`Unknown operation: ${operation}`);
            }

            onResult(result);
        } catch (e) {
            error = e instanceof Error ? e.message : 'Operation failed';
        } finally {
            loading = false;
        }
    }
</script>

<div class="operations-panel">
    <h3>Operations</h3>

    {#if error}
        <p class="error">{error}</p>
    {/if}

    <div class="operations">
        <button
            on:click={() => performOperation('union')}
            disabled={loading || selectedIds.length < 2}
        >
            Union (∪)
        </button>

        <button
            on:click={() => performOperation('intersection')}
            disabled={loading || selectedIds.length < 2}
        >
            Intersection (∩)
        </button>

        <button
            on:click={() => performOperation('complement')}
            disabled={loading || selectedIds.length !== 1}
        >
            Complement (¬)
        </button>

        <button
            on:click={() => performOperation('minimize')}
            disabled={loading || selectedIds.length !== 1}
        >
            Minimize DFA
        </button>
    </div>

    <p class="selected-count">
        Selected: {selectedIds.length} FA{selectedIds.length !== 1 ? 's' : ''}
    </p>
</div>

<style>
    .operations-panel {
        padding: 1rem;
        border: 1px solid #ddd;
        border-radius: 4px;
    }

    .operations {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 0.5rem;
        margin: 1rem 0;
    }

    button {
        padding: 0.5rem 1rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        background-color: #fff;
        cursor: pointer;
        transition: background-color 0.2s;
    }

    button:hover:not(:disabled) {
        background-color: #f5f5f5;
    }

    button:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .error {
        color: #d32f2f;
        margin: 0.5rem 0;
    }

    .selected-count {
        text-align: center;
        color: #666;
        font-size: 0.875rem;
    }
</style>
