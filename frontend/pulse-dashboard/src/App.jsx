import { useState, useEffect } from 'react'
import './App.css'
import CpuChart from "./components/charts/CpuChart";
import RamChart from "./components/charts/RamChart";
import LoadChart from "./components/charts/LoadChart";
import HostSelector from "./components/controls/HostSelector";
import StatCard from "./components/cards/StatCard";
import EventTimeline from "./components/timeline/EventTimeline";
import HostStatusCard from "./components/cards/HostStatusCard";
import { 
	fetchMetrics, 
	fetchHosts, 
	fetchStats, 
	fetchEvents,
	fetchLatest,
} from "./api/metrics";

export default function App() {
	const [hosts, setHosts] = useState([]);
	const [host, setHost] = useState("");
	const [data, setData] = useState([]);
	const [stats, setStats] = useState(null);
	const [hours, setHours] = useState(24);
	const [events, setEvents] = useState([]);
	const [latest, setLatest] = useState(null);

	async function load() {
		if (!host) return;

		const metrics = await fetchMetrics(host, hours);

		const formatted = metrics.map(m => ({
			...m,
			uptimeHours: m.uptime/3600,
			date: new Date(m.timestamp * 1000),
			timestamp: m.timestamp,
		}));

		setData(formatted);

		const statData = await fetchStats(host);
		setStats(statData);

		const eventData = await fetchEvents(host, hours);
		console.log("EVENT DATA:", eventData);
		console.log("EVENT COUNT:", eventData.length);
		setEvents(eventData);

		const latestData = await fetchLatest(host);
		setLatest(latestData);
	}

	useEffect(() => {
		fetchHosts().then((h) => {
			setHosts(h);
			setHost(h[0]);
		});
	}, []);

	useEffect(() => {
		load();
		const interval = setInterval(load, 5000);
		return () => clearInterval(interval);
	}, [host, hours]);

	return (
		<div style={{ padding: 20 }}>
		<div className="dashboard">
		<header className="topbar">
		<h1>Pulse Dashboard</h1>
		<HostSelector
		hosts={hosts}
		selected={host}
		onChange={setHost}
		/>
		</header>

		<HostStatusCard latest={latest} />

		<div className="stats-grid">

		<div className="range-selector">
		<button
		className={hours === 1 ? "active" : ""}
		onClick={() => setHours(1)}
		>
		1h
		</button>

		<button
		className={hours === 6 ? "active" : ""}
		onClick={() => setHours(6)}
		>
		6h
		</button>

		<button
		className={hours === 24 ? "active" : ""}
		onClick={() => setHours(24)}
		>
		24h
		</button>

		<button
		className={hours === 168 ? "active" : ""}
		onClick={() => setHours(168)}
		>
		7d
		</button>
		</div>


		<div className="stats-row">

		<StatCard
		title="Avg CPU"
		value={stats ? stats.avg_cpu.toFixed(1) : "--"}
		unit="%"
		/>

		<StatCard
		title="Avg RAM"
		value={stats ? stats.avg_ram.toFixed(1) : "--"}
		unit="%"
		/>

		<StatCard
		title="Avg Load"
		value={stats ? stats.avg_load1.toFixed(2) : "--"}
		unit=""
		/>

		<StatCard
		title="Samples"
		value={stats ? stats.samples : "--"}
		unit=""
		/>
		</div>
		</div>
		<div className="grid">

		<div className="card">
		<h2>CPU</h2>
		<CpuChart data={data} />
		</div>

		<div className="card">
		<h2>RAM</h2>
		<RamChart data={data} />
		</div>

		<div className="card">
		<h2>Load</h2>
		<LoadChart data={data} />
		</div>

		</div>
		<EventTimeline events={events} />
		</div>
		</div>

	);
}

