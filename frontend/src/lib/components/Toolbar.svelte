<script lang="ts">
    interface Props {
        hasSelection?: boolean;
        hasCurrentFA?: boolean;
        canSave?: boolean;
        onadd?: () => void;
        onedit?: () => void;
        onremove?: () => void;
        onsave?: () => void;
        ondownload?: () => void;
    }

    let { 
        hasSelection = false, 
        hasCurrentFA = false, 
        canSave = false,
        onadd,
        onedit,
        onremove,
        onsave,
        ondownload
    }: Props = $props();
</script>

<div class="toolbar">
    <div class="toolbar-group">
        <button class="toolbar-btn" onclick={() => onadd?.()} title="Add new FA">
            <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
            >
                <path d="M12 5v14M5 12h14" />
            </svg>
            Add
        </button>

        <button
            class="toolbar-btn"
            onclick={() => onedit?.()}
            disabled={!hasSelection}
            title="Edit selected FA"
        >
            <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
            >
                <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                <path d="m18.5 2.5 3 3L12 15l-4 1 1-4 9.5-9.5z" />
            </svg>
            Edit
        </button>

        <button
            class="toolbar-btn"
            onclick={() => onremove?.()}
            disabled={!hasSelection}
            title="Remove selected FA"
        >
            <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
            >
                <path
                    d="M3 6h18M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2m3 0v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6h14z"
                />
                <path d="M10 11v6M14 11v6" />
            </svg>
            Remove
        </button>

        <div class="separator"></div>

        <button
            class="toolbar-btn"
            onclick={() => onsave?.()}
            disabled={!canSave}
            title="Save current FA to database"
        >
            <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
            >
                <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z" />
                <polyline points="17,21 17,13 7,13 7,21" />
                <polyline points="7,3 7,8 15,8" />
            </svg>
            Save
        </button>

        <button
            class="toolbar-btn"
            onclick={() => ondownload?.()}
            disabled={!hasCurrentFA}
            title="Download current SVG"
        >
            <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
            >
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                <polyline points="7,10 12,15 17,10" />
                <line x1="12" y1="15" x2="12" y2="3" />
            </svg>
            Download
        </button>
    </div>
</div>

<style>
    .toolbar {
        display: flex;
        align-items: center;
        padding: 0.75rem 1rem;
        background-color: #ffffff;
        border-bottom: 1px solid #e0e0e0;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    }

    .toolbar-group {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }

    .toolbar-btn {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.5rem 0.75rem;
        background: none;
        border: 1px solid #e0e0e0;
        border-radius: 4px;
        cursor: pointer;
        font-size: 0.875rem;
        color: #374151;
        transition: all 0.2s ease;
    }

    .toolbar-btn:hover:not(:disabled) {
        background-color: #f3f4f6;
        border-color: #d1d5db;
    }

    .toolbar-btn:active:not(:disabled) {
        transform: translateY(1px);
        background-color: #e5e7eb;
    }

    .toolbar-btn:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .toolbar-btn svg {
        flex-shrink: 0;
    }

    .separator {
        width: 1px;
        height: 24px;
        background-color: #e0e0e0;
        margin: 0 0.25rem;
    }
</style>
