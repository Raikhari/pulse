const BASE_URL = "http://localhost:8080";

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
		`http://localhost:8080/stats?host=${host}`
	);

	return await res.json();
}

export async function fetchEvents(host, hours = 24) {
    const res = await fetch(
        `http://localhost:8080/events?host=${host}&hours=${hours}`
    );

    return await res.json();
}

export async function fetchLatest(host) {
    const response = await fetch(
        `http://localhost:8080/metrics/latest?host=${encodeURIComponent(host)}`
    );

    if (!response.ok) {
        throw new Error("Failed to fetch latest metrics");
    }

    return response.json();
}
