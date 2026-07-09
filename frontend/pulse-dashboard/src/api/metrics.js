const BASE_URL = "http://192.168.1.211:8080";

export async function fetchMetrics(host, hours = 24) {
	const res = await fetch(`${BASE_URL}/metrics?host=${host}&limit=200&hours=${hours}`);
	return res.json();
}

export async function fetchHosts() {
	const res = await fetch(`${BASE_URL}/debug/hosts`);
	return res.json();
}

export async function fetchStats(host) {
	const res = await fetch(
		`${BASE_URL}/stats?host=${host}`
	);

	return await res.json();
}

export async function fetchEvents(host, hours = 24) {
	const res = await fetch(
		`${BASE_URL}/events?host=${host}&hours=${hours}`
	);

	if (!res.ok) {
		return[];
	}

	const data = await res.json();

	return await Array.isArray(data) ? data : [];
}

export async function fetchLatest(host) {
	const response = await fetch(
		`${BASE_URL}/metrics/latest?host=${encodeURIComponent(host)}`
	);

	if (!response.ok) {
		throw new Error("Failed to fetch latest metrics");
	}

	return response.json();
}

export async function fetchConfig() {
    const response = await fetch(`${BASE_URL}/config`);

    if (!response.ok) {
        throw new Error("Failed to fetch config");
    }

    return response.json();
}

export async function saveConfig(config) {
    const response = await fetch(`${BASE_URL}/config`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(config),
    });

    if (!response.ok) {
        const message = await response.text();
        throw new Error(message || "Failed to save config");
    }

    return response.json();
}
