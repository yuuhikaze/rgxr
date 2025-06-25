<script lang="ts">
    import { onMount } from 'svelte';
    import { api, type FARecord } from '$lib/api/client';

    export let onSelect: (fa: FARecord) => void;
    export let selectedIds: string[] = [];

    let fas: FARecord[] = [];
    let loading = true;
    let error: string | null = null;

    onMount(async () => {
        try {
            fas = await api.getAllFAs();
        } catch (e) {
            error = e instanceof Error ? e.message : 'Failed to load FAs';
        } finally {
            loading = false;
        }
    });

    function toggleSelection(fa: FARecord) {
        if (selectedIds.includes(fa.id)) {
            selectedIds = selectedIds.filter((id) => id !== fa.id);
        } else {
            selectedIds = [...selectedIds, fa.id];
        }
    }
</script>

<div class="fa-list">
    <h2>Finite Automata</h2>

    {#if loading}
        <p>Loading...</p>
    {:else if error}
        <p class="error">{error}</p>
    {:else if fas.length === 0}
        <p>No finite automata found</p>
    {:else}
        <div class="list">
            {#each fas as fa}
                <div class="fa-item" class:selected={selectedIds.includes(fa.id)}>
                    <input
                        type="checkbox"
                        checked={selectedIds.includes(fa.id)}
                        on:change={() => toggleSelection(fa)}
                    />
                    <button on:click={() => onSelect(fa)}>
                        {fa.description || `FA ${fa.id.slice(0, 8)}`}
                    </button>
                </div>
            {/each}
        </div>
    {/if}
</div>

<style>
    .fa-list {
        display: flex;
        flex-direction: column;
        height: 100%;
        padding: 1rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        overflow: hidden;
    }

    .list {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        flex: 1;
        overflow-y: auto;
        min-height: 0;
    }

    .fa-item {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.5rem;
        border: 1px solid #eee;
        border-radius: 4px;
        transition: background-color 0.2s;
    }

    .fa-item:hover {
        background-color: #f5f5f5;
    }

    .fa-item.selected {
        background-color: #e3f2fd;
    }

    .fa-item button {
        flex: 1;
        text-align: left;
        background: none;
        border: none;
        cursor: pointer;
        padding: 0.25rem;
    }

    .error {
        color: #d32f2f;
    }
</style>
