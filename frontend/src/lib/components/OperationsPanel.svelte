<script lang="ts">
    import { api } from '$lib/api/client';

    export let selectedIds: string[] = [];
    export let onResult: (result: any) => void;

    let loading = false;
    let error: string | null = null;
    let regexInput = '';
    let stringInput = '';
    let runResult: { accepted: boolean; path: string[] } | null = null;

    async function performOperation(operation: string) {
        if (selectedIds.length === 0 && operation !== 'regex-to-nfa') {
            error = 'Please select at least one FA';
            return;
        }

        loading = true;
        error = null;
        runResult = null;

        try {
            let result;

            switch (operation) {
                case 'union':
                    if (selectedIds.length < 2) {
                        throw new Error('Union requires at least 2 FAs');
                    }
                    result = await api.union(selectedIds);
                    onResult(result);
                    break;

                case 'intersection':
                    if (selectedIds.length < 2) {
                        throw new Error('Intersection requires at least 2 FAs');
                    }
                    const intersectionResponse = await fetch('/api/intersection', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ uuids: selectedIds, mode: 'intersection' })
                    });
                    result = await intersectionResponse.json();
                    onResult(result);
                    break;

                case 'concatenation':
                    if (selectedIds.length < 1) {
                        throw new Error('Concatenation requires at least 1 FA');
                    }
                    result = await api.concatenation(selectedIds);
                    onResult(result);
                    break;

                case 'complement':
                    if (selectedIds.length !== 1) {
                        throw new Error('Complement requires exactly 1 FA');
                    }
                    const compResponse = await fetch(`/api/complement?uuid=${selectedIds[0]}`);
                    result = await compResponse.json();
                    onResult(result);
                    break;

                case 'minimize':
                    if (selectedIds.length !== 1) {
                        throw new Error('Minimize requires exactly 1 FA');
                    }
                    const minResponse = await fetch(`/api/minimize-dfa?uuid=${selectedIds[0]}`);
                    result = await minResponse.json();
                    onResult(result);
                    break;

                case 'nfa-to-dfa':
                    if (selectedIds.length !== 1) {
                        throw new Error('NFA to DFA requires exactly 1 FA');
                    }
                    result = await api.nfaToDFA(selectedIds[0]);
                    onResult(result);
                    break;

                case 'fa-to-regex':
                    if (selectedIds.length !== 1) {
                        throw new Error('FA to Regex requires exactly 1 FA');
                    }
                    const regex = await api.faToRegex(selectedIds[0]);
                    // Display the regex result
                    error = `Regular Expression: ${regex}`;
                    break;

                case 'regex-to-nfa':
                    if (!regexInput) {
                        throw new Error('Please enter a regular expression');
                    }
                    result = await api.regexToNFA(regexInput);
                    onResult(result);
                    break;

                case 'run-string':
                    if (selectedIds.length !== 1) {
                        throw new Error('Run string requires exactly 1 FA');
                    }
                    if (!stringInput) {
                        throw new Error('Please enter a string to test');
                    }
                    runResult = await api.runString(selectedIds[0], stringInput);
                    break;

                default:
                    throw new Error(`Unknown operation: ${operation}`);
            }
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
            on:click={() => performOperation('concatenation')}
            disabled={loading || selectedIds.length < 1}
        >
            Concatenation
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

        <button
            on:click={() => performOperation('nfa-to-dfa')}
            disabled={loading || selectedIds.length !== 1}
        >
            NFA to DFA
        </button>

        <button
            on:click={() => performOperation('fa-to-regex')}
            disabled={loading || selectedIds.length !== 1}
        >
            FA to Regex
        </button>
    </div>

    <div class="regex-section">
        <h4>Regex to NFA</h4>
        <div class="input-group">
            <input
                type="text"
                bind:value={regexInput}
                placeholder="Enter regular expression"
                disabled={loading}
            />
            <button
                on:click={() => performOperation('regex-to-nfa')}
                disabled={loading || !regexInput}
            >
                Convert
            </button>
        </div>
    </div>

    <div class="run-section">
        <h4>Test String</h4>
        <div class="input-group">
            <input
                type="text"
                bind:value={stringInput}
                placeholder="Enter test string"
                disabled={loading}
            />
            <button
                on:click={() => performOperation('run-string')}
                disabled={loading || selectedIds.length !== 1 || !stringInput}
            >
                Run
            </button>
        </div>
        
        {#if runResult}
            <div class="run-result" class:accepted={runResult.accepted} class:rejected={!runResult.accepted}>
                <p><strong>Result:</strong> {runResult.accepted ? 'Accepted' : 'Rejected'}</p>
                <p><strong>Path:</strong> {runResult.path.join(' → ')}</p>
            </div>
        {/if}
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

    h3, h4 {
        margin: 0 0 0.5rem 0;
    }

    .operations {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 0.5rem;
        margin: 1rem 0;
    }

    .regex-section, .run-section {
        margin-top: 1.5rem;
        padding-top: 1rem;
        border-top: 1px solid #e0e0e0;
    }

    .input-group {
        display: flex;
        gap: 0.5rem;
        margin-top: 0.5rem;
    }

    .input-group input {
        flex: 1;
        padding: 0.5rem;
        border: 1px solid #ddd;
        border-radius: 4px;
    }

    .run-result {
        margin-top: 1rem;
        padding: 0.75rem;
        border-radius: 4px;
        font-size: 0.875rem;
    }

    .run-result.accepted {
        background-color: #d4edda;
        border: 1px solid #c3e6cb;
        color: #155724;
    }

    .run-result.rejected {
        background-color: #f8d7da;
        border: 1px solid #f5c6cb;
        color: #721c24;
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
        margin-top: 1rem;
    }
</style>
