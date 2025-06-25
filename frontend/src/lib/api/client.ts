export interface FA {
    alphabet: string[];
    states: string[];
    initial: string;
    acceptance: string[];
    transitions: (string | string[])[][];
}

export interface RenderResponse {
    id: string;
    svg: string;
    tex: string;
    dot: string;
}

export interface FARecord {
    id: string;
    description?: string;
    tuple: FA;
    render: string;
    created_at: string;
}

export class APIClient {
    private baseURL: string;
    private token: string | null = null;

    constructor(baseURL: string = '') {
        this.baseURL = baseURL;
    }

    async setEditorToken() { // change this later
        const email = 'editor@example.com'
        const pass = 'securepassword'
        const res = await fetch(`${this.baseURL}/pgapi/rpc/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, pass })
        });
        const data = await res.json();
        this.setToken(data.token);
    }

    public setToken(token: any) {
        localStorage.setItem('jwt', token);
        this.token = token;
    }

    private authHeaders(): HeadersInit {
        console.log("YO TOKEN: " + this.token)
        const headers: HeadersInit = {
            'Content-Type': 'application/json'
        };
        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }
        return headers;
    }

    // Render FA to SVG/TeX
    async renderFA(fa: FA): Promise<RenderResponse> {
        const response = await fetch(`${this.baseURL}/api/render`, {
            method: 'POST',
            headers: this.authHeaders(),
            body: JSON.stringify({ fa })
        });

        if (!response.ok) {
            throw new Error(`Conversion failed: ${response.statusText}`);
        }

        return response.json();
    }

    // Render FA by UUID
    async renderByUUID(uuid: string): Promise<RenderResponse> {
        const response = await fetch(`${this.baseURL}/api/render`, {
            method: 'POST',
            headers: this.authHeaders(),
            body: JSON.stringify({ uuid })
        });

        if (!response.ok) {
            throw new Error(`Render failed: ${response.statusText}`);
        }

        return response.json();
    }

    // Get all FAs from PostgREST
    async getAllFAs(): Promise<FARecord[]> {
        const response = await fetch(`${this.baseURL}/pgapi/finite_automatas`);

        if (!response.ok) {
            throw new Error(`Failed to fetch FAs: ${response.statusText}`);
        }

        return response.json();
    }

    // Get FA by UUID
    async getFA(uuid: string): Promise<FARecord> {
        const response = await fetch(`${this.baseURL}/pgapi/finite_automatas?id=eq.${uuid}`);

        if (!response.ok) {
            throw new Error(`Failed to fetch FA: ${response.statusText}`);
        }

        const results = await response.json();
        if (results.length === 0) {
            throw new Error('FA not found');
        }

        return results[0];
    }

    // Union multiple FAs
    async union(uuids: string[]): Promise<FA> {
        const response = await fetch(`${this.baseURL}/api/union`, {
            method: 'POST',
            headers: this.authHeaders(),
            body: JSON.stringify({ uuids })
        });

        if (!response.ok) {
            throw new Error(`Union failed: ${response.statusText}`);
        }

        return response.json();
    }

    // Concatenation of multiple FAs
    async concatenation(uuids: string[]): Promise<FA> {
        const response = await fetch(`${this.baseURL}/api/concatenation`, {
            method: 'POST',
            headers: this.authHeaders(),
            body: JSON.stringify({ uuids })
        });

        if (!response.ok) {
            throw new Error(`Concatenation failed: ${response.statusText}`);
        }

        return response.json();
    }

    // Convert FA to regular expression
    async faToRegex(uuid: string): Promise<string> {
        const response = await fetch(`${this.baseURL}/api/fa-to-regex`, {
            method: 'POST',
            headers: this.authHeaders(),
            body: JSON.stringify({ uuid })
        });

        if (!response.ok) {
            throw new Error(`FA to regex conversion failed: ${response.statusText}`);
        }

        const result = await response.json();
        return result.regex;
    }

    // Convert regular expression to NFA
    async regexToNFA(regex: string): Promise<FA> {
        const response = await fetch(`${this.baseURL}/api/regex-to-nfa`, {
            method: 'POST',
            headers: this.authHeaders(),
            body: JSON.stringify({ regex })
        });

        if (!response.ok) {
            throw new Error(`Regex to NFA conversion failed: ${response.statusText}`);
        }

        return response.json();
    }

    // Convert NFA to DFA
    async nfaToDFA(uuid: string): Promise<FA> {
        const response = await fetch(`${this.baseURL}/api/nfa-to-dfa`, {
            method: 'POST',
            headers: this.authHeaders(),
            body: JSON.stringify({ uuid })
        });

        if (!response.ok) {
            throw new Error(`NFA to DFA conversion failed: ${response.statusText}`);
        }

        return response.json();
    }

    // Run a string through an FA
    async runString(uuid: string, input: string): Promise<{ accepted: boolean; path: string[] }> {
        const response = await fetch(`${this.baseURL}/api/run-string`, {
            method: 'POST',
            headers: this.authHeaders(),
            body: JSON.stringify({ uuid, string: input })
        });

        if (!response.ok) {
            throw new Error(`Run string failed: ${response.statusText}`);
        }

        return response.json();
    }

    // Save FA to database
    async saveFA(fa: FA, description?: string): Promise<void> {
        const id = crypto.randomUUID();

        // First render the FA to get the render path
        const renderResult = await this.renderFA(fa);

        const response = await fetch(`${this.baseURL}/pgapi/finite_automatas`, {
            method: 'POST',
            headers: this.authHeaders(),
            body: JSON.stringify({
                id,
                tuple: fa,
                render: renderResult.id,
                description
            })
        });

        if (!response.ok) {
            throw new Error(`Failed to save FA: ${response.statusText}`);
        }
    }

    // Update FA in database
    async updateFA(uuid: string, fa: FA, description?: string): Promise<void> {
        // First render the FA to get the render path
        const renderResult = await this.renderFA(fa);

        const response = await fetch(`${this.baseURL}/pgapi/finite_automatas?id=eq.${uuid}`, {
            method: 'PATCH',
            headers: this.authHeaders(),
            body: JSON.stringify({
                tuple: fa,
                render: renderResult.id,
                description
            })
        });

        if (!response.ok) {
            throw new Error(`Failed to update FA: ${response.statusText}`);
        }
    }

    // Delete FA from database
    async deleteFA(uuid: string): Promise<void> {
        const response = await fetch(`${this.baseURL}/pgapi/finite_automatas?id=eq.${uuid}`, {
            headers: this.authHeaders(),
            method: 'DELETE'
        });

        if (!response.ok) {
            throw new Error(`Failed to delete FA: ${response.statusText}`);
        }
    }

    // Get TeX code
    async getTeX(uuid: string): Promise<string> {
        const response = await fetch(`${this.baseURL}/api/tex/${uuid}`);

        if (!response.ok) {
            throw new Error(`Failed to get TeX: ${response.statusText}`);
        }

        return response.text();
    }

    // Get SVG
    async getSVG(uuid: string): Promise<string> {
        const response = await fetch(`${this.baseURL}/api/svg/${uuid}`);

        if (!response.ok) {
            throw new Error(`Failed to get SVG: ${response.statusText}`);
        }

        return response.text();
    }
}

export const api = new APIClient();
