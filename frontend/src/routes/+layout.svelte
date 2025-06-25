<script lang="ts">
    import '../app.css';
    import { onMount } from 'svelte';
    import { api } from '$lib/api/client';

    let { children } = $props();

    onMount(() => {
        const token = localStorage.getItem('jwt');
        if (token) {
            api.setToken(token);
        } else {
            api.setEditorToken();
        }
    });
</script>

{@render children()}
