<script lang="ts">
    import type { FA } from '$lib/api/client';

    interface Props {
        open?: boolean;
        editingFA?: FA | null;
        isEditing?: boolean;
        onsave?: (event: CustomEvent<{ fa: FA; description?: string }>) => void;
        oncancel?: (event: CustomEvent<void>) => void;
    }

    let {
        open = $bindable(false),
        editingFA = $bindable(null),
        isEditing = $bindable(false),
        onsave,
        oncancel
    }: Props = $props();

    // Form state
    let alphabet = $state('');
    let states = $state('');
    let initialState = $state('');
    let acceptanceStates = $state<string[]>([]);
    let description = $state('');
    let selectedState = $state('');
    let selectedCell = $state<{ row: number; col: number } | null>(null);
    let table = $state<(string | string[])[][]>([]);

    // State types for visual representation
    const STATE_TYPES = {
        NORMAL: 0,
        ACCEPTANCE: 1,
        INITIAL: 2,
        BOTH: 3
    };

    // Initialize form when editing - using $effect for Svelte 5
    $effect(() => {
        if (open && editingFA && isEditing) {
            alphabet = editingFA.alphabet.join(', ');
            states = editingFA.states.join(', ');
            initialState = editingFA.initial;
            acceptanceStates = [...editingFA.acceptance];
            table = editingFA.transitions.map((row) => [...row]);
        } else if (open && !isEditing) {
            // Reset form for new FA
            alphabet = '';
            states = '';
            initialState = '';
            acceptanceStates = [];
            description = '';
            selectedState = '';
            selectedCell = null;
            table = [];
        }
    });

    // Update table when alphabet or states change - using $effect for Svelte 5
    $effect(() => {
        if (alphabet && states) {
            const stateList = states
                .split(',')
                .map((s) => s.trim())
                .filter((s) => s);
            const symbolList = alphabet
                .split(',')
                .map((s) => s.trim())
                .filter((s) => s);

            if (!isEditing || table.length === 0) {
                table = stateList.map(() => symbolList.map(() => ''));
            }
        }
    });

    function getStateType(state: string): number {
        const isAcceptance = acceptanceStates.includes(state);
        const isInitial = initialState === state;

        if (isAcceptance && isInitial) return STATE_TYPES.BOTH;
        if (isAcceptance) return STATE_TYPES.ACCEPTANCE;
        if (isInitial) return STATE_TYPES.INITIAL;
        return STATE_TYPES.NORMAL;
    }

    function setStateType(state: string, type: number): void {
        // Remove from acceptance states if needed
        if (type === STATE_TYPES.NORMAL || type === STATE_TYPES.INITIAL) {
            acceptanceStates = acceptanceStates.filter((s) => s !== state);
        }

        // Add to acceptance states if needed
        if (type === STATE_TYPES.ACCEPTANCE || type === STATE_TYPES.BOTH) {
            if (!acceptanceStates.includes(state)) {
                acceptanceStates = [...acceptanceStates, state];
            }
        }

        // Set or unset as initial state
        if (type === STATE_TYPES.INITIAL || type === STATE_TYPES.BOTH) {
            initialState = state;
        } else if (initialState === state) {
            initialState = '';
        }
    }

    function cycleStateType(state: string): void {
        const currentType = getStateType(state);
        const nextType = (currentType + 1) % 4;
        setStateType(state, nextType);
    }

    function getStateStyles(state: string): string {
        const stateType = getStateType(state);
        let classes = '';

        switch (stateType) {
            case STATE_TYPES.ACCEPTANCE:
                classes += 'acceptance-state ';
                break;
            case STATE_TYPES.INITIAL:
                classes += 'initial-state ';
                break;
            case STATE_TYPES.BOTH:
                classes += 'acceptance-state initial-state ';
                break;
        }

        return classes;
    }

    function handleCellClick(rowIdx: number, colIdx: number) {
        selectedCell = { row: rowIdx, col: colIdx };
        if (selectedState) {
            if (!table[rowIdx]) table[rowIdx] = [];
            
            // Get current cell value
            const currentValue = table[rowIdx][colIdx];
            
            // If cell is empty, add the selected state
            if (!currentValue || currentValue === '') {
                table[rowIdx][colIdx] = selectedState;
            } else {
                // Convert to array if it's not already
                let statesArray: string[];
                if (Array.isArray(currentValue)) {
                    statesArray = [...currentValue];
                } else {
                    statesArray = [currentValue];
                }
                
                // Toggle the selected state
                const stateIndex = statesArray.indexOf(selectedState);
                if (stateIndex === -1) {
                    // Add state if not present
                    statesArray.push(selectedState);
                } else {
                    // Remove state if present
                    statesArray.splice(stateIndex, 1);
                }
                
                // Update table with new value
                if (statesArray.length === 0) {
                    table[rowIdx][colIdx] = '';
                } else if (statesArray.length === 1) {
                    table[rowIdx][colIdx] = statesArray[0];
                } else {
                    table[rowIdx][colIdx] = statesArray;
                }
            }
        }
    }

    function handleStateSelect(state: string) {
        selectedState = state;
    }

    function handleKeyDown(event: KeyboardEvent) {
        if (event.key === 'Backspace' && selectedCell) {
            const { row, col } = selectedCell;
            if (table[row] && table[row][col] !== undefined) {
                table[row][col] = '';
            }
            event.preventDefault();
        }
    }

    function handleSave() {
        const stateList = states
            .split(',')
            .map((s) => s.trim())
            .filter((s) => s);
        const symbolList = alphabet
            .split(',')
            .map((s) => s.trim())
            .filter((s) => s);

        const fa: FA = {
            alphabet: symbolList,
            states: stateList,
            initial: initialState,
            acceptance: acceptanceStates,
            transitions: table
        };

        onsave?.(
            new CustomEvent('save', { detail: { fa, description: description || undefined } })
        );
    }

    function handleCancel() {
        oncancel?.(new CustomEvent('cancel'));
    }
</script>

<svelte:window on:keydown={handleKeyDown} />

{#if open}
    <div class="modal-overlay" on:click={handleCancel}>
        <div class="modal-content" on:click|stopPropagation>
            <div class="modal-header">
                <h2>{isEditing ? 'Edit' : 'Create'} Finite Automaton</h2>
                <button class="close-btn" on:click={handleCancel}>
                    <svg
                        width="24"
                        height="24"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                    >
                        <path d="M18 6L6 18M6 6l12 12" />
                    </svg>
                </button>
            </div>

            <div class="modal-body">
                <p class="description">Both NFA and DFA creation are supported.</p>

                <div class="form-section">
                    <div class="form-group">
                        <label for="description">Description (optional)</label>
                        <input
                            id="description"
                            class="form-input"
                            placeholder="Description for this FA"
                            bind:value={description}
                        />
                    </div>

                    <div class="form-group">
                        <label for="alphabet">Alphabet</label>
                        <input
                            id="alphabet"
                            class="form-input"
                            placeholder="0, 1, @e"
                            bind:value={alphabet}
                        />
                    </div>

                    <div class="form-group">
                        <label for="states">States</label>
                        <input
                            id="states"
                            class="form-input"
                            placeholder="q0, q1, q2"
                            bind:value={states}
                        />
                    </div>
                </div>

                {#if alphabet && states}
                    <div class="transition-table">
                        <h3>Transition Table</h3>
                        <p class="table-help">
                            Click a state to select it, then click cells to toggle transitions. Click the same cell multiple times to add/remove states (for NFA). Scroll on states to change their type.
                        </p>

                        <table>
                            <thead>
                                <tr>
                                    <th class="state-header">δ</th>
                                    {#each alphabet.split(',').map((s) => s.trim()) as symbol}
                                        <th
                                            class="symbol-header"
                                            class:selected={selectedState === symbol}
                                            on:click={() => handleStateSelect(symbol)}
                                        >
                                            {symbol === '@e' ? 'ε' : symbol}
                                        </th>
                                    {/each}
                                </tr>
                            </thead>
                            <tbody>
                                {#each states.split(',').map((s) => s.trim()) as state, rowIdx}
                                    <tr>
                                        <td
                                            class="state-cell {getStateStyles(state)}"
                                            class:selected={selectedState === state}
                                            on:click={() => handleStateSelect(state)}
                                            on:wheel|preventDefault={(e) => cycleStateType(state)}
                                        >
                                            {state}
                                        </td>
                                        {#each alphabet
                                            .split(',')
                                            .map((s) => s.trim()) as symbol, colIdx}
                                            <td
                                                class="transition-cell"
                                                class:selected={selectedCell?.row === rowIdx &&
                                                    selectedCell?.col === colIdx}
                                                on:click={() => handleCellClick(rowIdx, colIdx)}
                                            >
                                                {#if table[rowIdx]?.[colIdx]}
                                                    {Array.isArray(table[rowIdx][colIdx])
                                                        ? table[rowIdx][colIdx].join(', ')
                                                        : table[rowIdx][colIdx]}
                                                {:else}
                                                    ∅
                                                {/if}
                                            </td>
                                        {/each}
                                    </tr>
                                {/each}
                            </tbody>
                        </table>
                    </div>
                {/if}
            </div>

            <div class="modal-footer">
                <button class="btn btn-secondary" on:click={handleCancel}> Cancel </button>
                <button
                    class="btn btn-primary"
                    on:click={handleSave}
                    disabled={!alphabet || !states || !initialState}
                >
                    {isEditing ? 'Update' : 'Create'}
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
        padding: 1rem;
    }

    .modal-content {
        background: white;
        border-radius: 8px;
        box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
        max-width: 90vw;
        max-height: 90vh;
        width: 800px;
        overflow: hidden;
        display: flex;
        flex-direction: column;
    }

    .modal-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1.5rem;
        border-bottom: 1px solid #e0e0e0;
    }

    .modal-header h2 {
        margin: 0;
        font-size: 1.25rem;
        font-weight: 600;
    }

    .close-btn {
        background: none;
        border: none;
        cursor: pointer;
        padding: 0.25rem;
        color: #6b7280;
        border-radius: 4px;
    }

    .close-btn:hover {
        background-color: #f3f4f6;
    }

    .modal-body {
        padding: 1.5rem;
        overflow-y: auto;
        flex: 1;
    }

    .description {
        color: #6b7280;
        margin-bottom: 1.5rem;
        font-size: 0.875rem;
    }

    .form-section {
        display: flex;
        flex-direction: column;
        gap: 1rem;
        margin-bottom: 1.5rem;
    }

    .form-group {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .form-group label {
        font-weight: 500;
        color: #374151;
        font-size: 0.875rem;
    }

    .form-input {
        padding: 0.75rem;
        border: 1px solid #d1d5db;
        border-radius: 4px;
        font-size: 0.875rem;
        transition: border-color 0.2s;
    }

    .form-input:focus {
        outline: none;
        border-color: #3b82f6;
        box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
    }

    .transition-table {
        margin-top: 1.5rem;
    }

    .transition-table h3 {
        margin: 0 0 0.5rem 0;
        font-size: 1rem;
        font-weight: 600;
    }

    .table-help {
        font-size: 0.75rem;
        color: #6b7280;
        margin-bottom: 1rem;
    }

    table {
        width: 100%;
        border-collapse: collapse;
        border: 1px solid #e0e0e0;
    }

    th,
    td {
        border: 1px solid #e0e0e0;
        padding: 0.5rem;
        text-align: center;
        min-width: 60px;
    }

    .state-header {
        background-color: #f9fafb;
        font-weight: 600;
    }

    .symbol-header {
        background-color: #f3f4f6;
        cursor: pointer;
        transition: background-color 0.2s;
    }

    .symbol-header:hover {
        background-color: #e5e7eb;
    }

    .symbol-header.selected {
        background-color: #dbeafe;
        color: #1d4ed8;
    }

    .state-cell {
        background-color: #f3f4f6;
        cursor: pointer;
        transition: background-color 0.2s;
        position: relative;
    }

    .state-cell:hover {
        background-color: #e5e7eb;
    }

    .state-cell.selected {
        background-color: #dbeafe;
        color: #1d4ed8;
    }

    .state-cell.acceptance-state {
        position: relative;
    }

    .state-cell.acceptance-state::after {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        border: 3px double #374151;
        pointer-events: none;
        box-sizing: border-box;
    }

    .state-cell.initial-state::before {
        content: '→';
        position: absolute;
        left: 2px;
        top: 50%;
        transform: translateY(-50%);
        font-size: 0.75rem;
        color: #059669;
    }

    .transition-cell {
        cursor: pointer;
        transition: background-color 0.2s;
    }

    .transition-cell:hover {
        background-color: #f9fafb;
    }

    .transition-cell.selected {
        background-color: #fef3c7;
    }

    .modal-footer {
        display: flex;
        justify-content: flex-end;
        gap: 0.75rem;
        padding: 1.5rem;
        border-top: 1px solid #e0e0e0;
        background-color: #f9fafb;
    }

    .btn {
        padding: 0.75rem 1.5rem;
        border-radius: 4px;
        font-size: 0.875rem;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.2s;
        border: 1px solid;
    }

    .btn-secondary {
        background-color: white;
        color: #374151;
        border-color: #d1d5db;
    }

    .btn-secondary:hover {
        background-color: #f9fafb;
    }

    .btn-primary {
        background-color: #3b82f6;
        color: white;
        border-color: #3b82f6;
    }

    .btn-primary:hover:not(:disabled) {
        background-color: #2563eb;
    }

    .btn-primary:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }
</style>
