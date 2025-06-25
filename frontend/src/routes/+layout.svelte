<script lang="ts">
    import '../app.css';
    import { api } from '$lib/api/client';

    let { children } = $props();
    api.setEditorToken();
    const token = localStorage.getItem('jwt');
    if (token) api.setToken(token);
</script>

{@render children()}
