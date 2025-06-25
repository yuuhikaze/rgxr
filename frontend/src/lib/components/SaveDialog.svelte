<script lang="ts">
    interface Props {
        open?: boolean;
        onsave?: (description: string) => void;
        oncancel?: () => void;
    }

    let { open = false, onsave, oncancel }: Props = $props();

    let description = $state('');

    function handleSave() {
        if (description.trim()) {
            onsave?.(description.trim());
            description = '';
        }
    }

    function handleCancel() {
        oncancel?.();
        description = '';
    }

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === 'Enter' && !event.shiftKey) {
            event.preventDefault();
            handleSave();
        } else if (event.key === 'Escape') {
            handleCancel();
        }
    }
</script>

{#if open}
    <div class="modal-overlay" onclick={handleCancel}>
        <div class="modal" onclick={(e) => e.stopPropagation()}>
            <div class="modal-header">
                <h2>Save Finite Automaton</h2>
                <button class="close-btn" onclick={handleCancel}>
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <line x1="18" y1="6" x2="6" y2="18"></line>
                        <line x1="6" y1="6" x2="18" y2="18"></line>
                    </svg>
                </button>
            </div>

            <div class="modal-body">
                <div class="form-group">
                    <label for="description">Description</label>
                    <input
                        id="description"
                        type="text"
                        bind:value={description}
                        placeholder="Enter a description for this FA..."
                        onkeydown={handleKeydown}
                        autofocus
                    />
                    <p class="help-text">Provide a brief description to help identify this finite automaton later.</p>
                </div>
            </div>

            <div class="modal-footer">
                <button class="btn btn-secondary" onclick={handleCancel}>
                    Cancel
                </button>
                <button 
                    class="btn btn-primary" 
                    onclick={handleSave}
                    disabled={!description.trim()}
                >
                    Save FA
                </button>
            </div>
        </div>
    </div>
{/if}

<style>
    .modal-overlay {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background-color: rgba(0, 0, 0, 0.5);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 1000;
    }

    .modal {
        background: white;
        border-radius: 8px;
        box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
        min-width: 400px;
        max-width: 500px;
        width: 90vw;
        max-height: 90vh;
        overflow: hidden;
        display: flex;
        flex-direction: column;
    }

    .modal-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 1.5rem 1.5rem 1rem;
        border-bottom: 1px solid #e0e0e0;
    }

    .modal-header h2 {
        margin: 0;
        font-size: 1.25rem;
        font-weight: 600;
        color: #1f2937;
    }

    .close-btn {
        background: none;
        border: none;
        cursor: pointer;
        padding: 0.25rem;
        border-radius: 4px;
        color: #6b7280;
        transition: all 0.2s ease;
    }

    .close-btn:hover {
        background-color: #f3f4f6;
        color: #374151;
    }

    .modal-body {
        padding: 1.5rem;
        flex: 1;
        overflow-y: auto;
    }

    .form-group {
        margin-bottom: 1rem;
    }

    .form-group label {
        display: block;
        margin-bottom: 0.5rem;
        font-weight: 500;
        color: #374151;
        font-size: 0.875rem;
    }

    .form-group input {
        width: 100%;
        padding: 0.75rem;
        border: 1px solid #d1d5db;
        border-radius: 4px;
        font-size: 0.875rem;
        transition: border-color 0.2s ease;
        box-sizing: border-box;
    }

    .form-group input:focus {
        outline: none;
        border-color: #3b82f6;
        box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
    }

    .help-text {
        margin: 0.5rem 0 0;
        font-size: 0.75rem;
        color: #6b7280;
        line-height: 1.4;
    }

    .modal-footer {
        display: flex;
        gap: 0.75rem;
        padding: 1rem 1.5rem 1.5rem;
        border-top: 1px solid #e0e0e0;
        justify-content: flex-end;
    }

    .btn {
        padding: 0.5rem 1rem;
        border-radius: 4px;
        font-size: 0.875rem;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.2s ease;
        border: 1px solid transparent;
        min-width: 80px;
    }

    .btn-secondary {
        background-color: #f9fafb;
        color: #374151;
        border-color: #d1d5db;
    }

    .btn-secondary:hover {
        background-color: #f3f4f6;
        border-color: #9ca3af;
    }

    .btn-primary {
        background-color: #3b82f6;
        color: white;
        border-color: #3b82f6;
    }

    .btn-primary:hover:not(:disabled) {
        background-color: #2563eb;
        border-color: #2563eb;
    }

    .btn-primary:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }
</style>
