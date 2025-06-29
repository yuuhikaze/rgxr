<script lang="ts">
    import { api } from '$lib/api/client';

    export let selectedIds: string[] = [];
    export let onResult: (result: any) => void;

    let loading = false;
    let error: string | null = null;
    let regexInput = '';
    let stringInput = '';
    let runResult: { accepted: boolean; path: string[] } | null = null;
    let activeTab: 'unary' | 'binary' | 'regex' = 'unary';

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
                        throw new Error('Deterministic union requires at least 2 FAs');
                    }
                    result = await api.boolean(selectedIds, 'union');
                    onResult(result);
                    break;

                case 'n-union':
                    if (selectedIds.length < 2) {
                        throw new Error('Non-deterministic union requires at least 2 FAs');
                    }
                    result = await api.nboolean(selectedIds, 'union');
                    onResult(result);
                    break;

                case 'intersection':
                    if (selectedIds.length < 2) {
                        throw new Error('Deterministic intersection requires at least 2 FAs');
                    }
                    result = await api.boolean(selectedIds, 'intersection');
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
    <div class="panel-header">
        <h3>Operations</h3>
        <p class="selected-count">
            Selected: {selectedIds.length} FA{selectedIds.length !== 1 ? 's' : ''}
        </p>
    </div>

    {#if error}
        <p class="error">{error}</p>
    {/if}

    <!-- Tab Navigation -->
    <div class="tab-nav">
        <button
            class="tab-button"
            class:active={activeTab === 'unary'}
            on:click={() => (activeTab = 'unary')}
        >
            Unary Ops
        </button>
        <button
            class="tab-button"
            class:active={activeTab === 'binary'}
            on:click={() => (activeTab = 'binary')}
        >
            Binary Ops
        </button>
        <button
            class="tab-button"
            class:active={activeTab === 'regex'}
            on:click={() => (activeTab = 'regex')}
        >
            Regex/String
        </button>
    </div>

    <!-- Tab Content -->
    <div class="tab-content">
        {#if activeTab === 'unary'}
            <div class="operations-grid">
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
        {/if}

        {#if activeTab === 'binary'}
            <div class="operations-grid">
                <button
                    on:click={() => performOperation('union')}
                    disabled={loading || selectedIds.length < 2}
                >
                    Union (∪)
                </button>

                <button
                    on:click={() => performOperation('n-union')}
                    disabled={loading || selectedIds.length < 2}
                >
                    N-Union (∪)
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
            </div>
        {/if}

        {#if activeTab === 'regex'}
            <div class="regex-string-section">
                <div class="input-section">
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

                <div class="input-section">
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
                        <div
                            class="run-result"
                            class:accepted={runResult.accepted}
                            class:rejected={!runResult.accepted}
                        >
                            <p>
                                <strong>Result:</strong>
                                {runResult.accepted ? 'Accepted' : 'Rejected'}
                            </p>
                            <p><strong>Path:</strong> {runResult.path.join(' → ')}</p>
                        </div>
                    {/if}
                </div>
            </div>
        {/if}
    </div>
</div>

<style>
    .operations-panel {
        display: flex;
        flex-direction: column;
        height: 100%;
        padding: 1rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        overflow: hidden;
    }

    .panel-header {
        flex-shrink: 0;
        margin-bottom: 1rem;
    }

    h3,
    h4 {
        margin: 0 0 0.5rem 0;
    }

    .selected-count {
        color: #666;
        font-size: 0.875rem;
        margin: 0.5rem 0 0 0;
    }

    .error {
        color: #d32f2f;
        margin: 0.5rem 0;
        flex-shrink: 0;
    }

    /* Tab Navigation */
    .tab-nav {
        display: flex;
        border-bottom: 1px solid #ddd;
        margin-bottom: 1rem;
        flex-shrink: 0;
    }

    .tab-button {
        flex: 1;
        padding: 0.75rem 1rem;
        border: none;
        border-bottom: 2px solid transparent;
        background-color: transparent;
        cursor: pointer;
        font-size: 0.875rem;
        font-weight: 500;
        transition: all 0.2s;
        color: #666;
    }

    .tab-button:hover {
        background-color: #f5f5f5;
        color: #333;
    }

    .tab-button.active {
        color: #007bff;
        border-bottom-color: #007bff;
        background-color: #f8f9fa;
    }

    /* Tab Content */
    .tab-content {
        flex: 1;
        overflow-y: auto;
        min-height: 0;
    }

    /* Operations Grid */
    .operations-grid {
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
    }

    .operations-grid button {
        width: 100%;
        padding: 0.75rem 1rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        background-color: #fff;
        cursor: pointer;
        transition: all 0.2s;
        font-size: 0.875rem;
        text-align: left;
        min-height: 44px; /* Ensure minimum touch target */
    }

    .operations-grid button:hover:not(:disabled) {
        background-color: #f5f5f5;
        border-color: #007bff;
        transform: translateY(-1px);
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }

    .operations-grid button:disabled {
        opacity: 0.5;
        cursor: not-allowed;
        transform: none;
        box-shadow: none;
    }

    /* Regex/String Section */
    .regex-string-section {
        display: flex;
        flex-direction: column;
        gap: 1.5rem;
    }

    .input-section {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .input-group {
        display: flex;
        gap: 0.5rem;
    }

    .input-group input {
        flex: 1;
        padding: 0.5rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        font-size: 0.875rem;
        min-width: 0; /* Allow input to shrink */
    }

    .input-group button {
        padding: 0.5rem 1rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        background-color: #fff;
        cursor: pointer;
        transition: all 0.2s;
        font-size: 0.875rem;
        white-space: nowrap;
        flex-shrink: 0;
    }

    .input-group button:hover:not(:disabled) {
        background-color: #f5f5f5;
        border-color: #007bff;
    }

    .input-group button:disabled {
        opacity: 0.5;
        cursor: not-allowed;
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

    /* Responsive adjustments */
    @media (max-width: 768px) {
        .operations-panel {
            padding: 0.75rem;
        }

        .tab-button {
            padding: 0.5rem 0.75rem;
            font-size: 0.8rem;
        }

        .operations-grid button {
            padding: 0.5rem 0.75rem;
            font-size: 0.8rem;
        }
    }
</style>
