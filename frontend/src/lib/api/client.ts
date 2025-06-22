export interface FA {
    alphabet: string[];
    states: string[];
    initial: string;
    acceptance: string[];
    transitions: (string | string[])[][];
}

export interface ConvertResponse {
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

    constructor(baseURL: string = '') {
        this.baseURL = baseURL;
    }

    // Convert FA to SVG/TeX
    async convertFA(fa: FA): Promise<ConvertResponse> {
        const response = await fetch(`${this.baseURL}/api/convert`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ fa }),
        });

        if (!response.ok) {
            throw new Error(`Conversion failed: ${response.statusText}`);
        }

        return response.json();
    }

    // Convert FA by UUID
    async convertByUUID(uuid: string): Promise<ConvertResponse> {
        const response = await fetch(`${this.baseURL}/api/convert`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ uuid }),
        });

        if (!response.ok) {
            throw new Error(`Conversion failed: ${response.statusText}`);
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
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ uuids }),
        });

        if (!response.ok) {
            throw new Error(`Union failed: ${response.statusText}`);
        }

        return response.json();
    }

    // Save FA to database
    async saveFA(fa: FA, renderPath: string, description?: string): Promise<void> {
        const id = crypto.randomUUID();

        const response = await fetch(`${this.baseURL}/pgapi/rpc/save_fa`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                id,
                tuple: fa,
                render: renderPath,
                description,
            }),
        });

        if (!response.ok) {
            throw new Error(`Failed to save FA: ${response.statusText}`);
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
