<script lang="ts">
    // import { onMount } from 'svelte';
    import FAList from '$lib/components/FAList.svelte';
    import GraphViewer from '$lib/components/GraphViewer.svelte';
    import OperationsPanel from '$lib/components/OperationsPanel.svelte';
    import Toolbar from '$lib/components/Toolbar.svelte';
    import FADialog from '$lib/components/FADialog.svelte';
    import { api, type FA, type FARecord } from '$lib/api/client';

    let selectedIds: string[] = [];
    let currentFA: FA | null = null;
    let currentSVG: string = '';
    let currentTeX: string = '';
    let loading = false;

    // Dialog state
    let showFADialog = false;
    let editingFA: FA | null = null;
    let isEditing = false;

    async function handleFASelect(fa: FARecord) {
        loading = true;
        try {
            const result = await api.renderByUUID(fa.id);
            currentFA = fa.tuple;
            currentSVG = result.svg;
            currentTeX = result.tex;
        } catch (e) {
            console.error('Failed to render FA:', e);
        } finally {
            loading = false;
        }
    }

    async function handleOperationResult(fa: FA) {
        loading = true;
        try {
            const result = await api.renderFA(fa);
            currentFA = fa;
            currentSVG = result.svg;
            currentTeX = result.tex;
        } catch (e) {
            console.error('Failed to render result:', e);
        } finally {
            loading = false;
        }
    }

    // Toolbar handlers
    function handleAdd() {
        editingFA = null;
        isEditing = false;
        showFADialog = true;
    }

    function handleEdit() {
        if (selectedIds.length === 1) {
            // TODO: Get the FA data for the selected ID and set editingFA
            editingFA = currentFA;
            isEditing = true;
            showFADialog = true;
        }
    }

    async function handleRemove() {
        if (selectedIds.length > 0) {
            try {
                for (const id of selectedIds) {
                    await api.deleteFA(id);
                }
                selectedIds = [];
                // Refresh the FA list
                window.location.reload();
            } catch (e) {
                console.error('Failed to remove FA:', e);
            }
        }
    }

    let saveMessage = '';
    let showSaveMessage = false;

    async function handleSave() {
        if (currentFA) {
            try {
                await api.saveFA(currentFA);
                saveMessage = 'FA saved successfully!';
                showSaveMessage = true;
                setTimeout(() => {
                    showSaveMessage = false;
                }, 3000);
            } catch (e) {
                console.error('Failed to save FA:', e);
                saveMessage = 'Failed to save FA';
                showSaveMessage = true;
                setTimeout(() => {
                    showSaveMessage = false;
                }, 3000);
            }
        }
    }

    function handleDownload() {
        if (currentSVG) {
            const blob = new Blob([currentSVG], { type: 'image/svg+xml' });
            const url = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'finite-automaton.svg';
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            URL.revokeObjectURL(url);
        }
    }

    // FA Dialog handlers
    async function handleFADialogSave(event: CustomEvent<{ fa: FA; description?: string }>) {
        try {
            const { fa, description } = event.detail;
            if (isEditing && selectedIds.length > 0) {
                // Update existing FA
                await api.updateFA(selectedIds[0], fa, description);
            } else {
                // Create new FA
                await api.saveFA(fa, description);
            }
            showFADialog = false;
            // Refresh the FA list
            window.location.reload();
        } catch (e) {
            console.error('Failed to save FA:', e);
            alert('Failed to save FA');
        }
    }

    function handleFADialogCancel() {
        showFADialog = false;
    }
</script>

<div class="app">
    <header>
        <h1>RGXR</h1>
    </header>

    <Toolbar
        hasSelection={selectedIds.length > 0}
        hasCurrentFA={currentFA !== null}
        canSave={currentFA !== null}
        onadd={handleAdd}
        onedit={handleEdit}
        onremove={handleRemove}
        onsave={handleSave}
        ondownload={handleDownload}
    />

    <main>
        <div class="sidebar">
            <div class="sidebar-section">
                <FAList bind:selectedIds onSelect={handleFASelect} />
            </div>
            <div class="sidebar-section">
                <OperationsPanel {selectedIds} onResult={handleOperationResult} />
            </div>
        </div>

        <div class="content">
            {#if loading}
                <div class="loading">Rendering...</div>
            {:else}
                <GraphViewer svg={currentSVG} tex={currentTeX} />
            {/if}
        </div>
    </main>
</div>

<FADialog
    bind:open={showFADialog}
    {editingFA}
    {isEditing}
    onsave={handleFADialogSave}
    oncancel={handleFADialogCancel}
/>

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
        overflow: hidden;
    }

    .sidebar-section {
        flex: 1;
        display: flex;
        flex-direction: column;
        min-height: 0;
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
