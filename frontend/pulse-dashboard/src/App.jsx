import { useState, useEffect } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from './assets/vite.svg'
import heroImg from './assets/hero.png'
import './App.css'
import CpuChart from "./components/charts/CpuChart";
import RamChart from "./components/charts/RamChart";
import HostSelector from "./components/controls/HostSelector";
import StatCard from "./components/cards/StatCard";
import { fetchMetrics, fetchHosts, fetchStats } from "./api/metrics";

export default function App() {
	const [hosts, setHosts] = useState([]);
	const [host, setHost] = useState("");
	const [data, setData] = useState([]);
	const [stats, setStats] = useState(null);

	async function load() {
		if (!host) return;

		const metrics = await fetchMetrics(host);

		const formatted = metrics.map(m => ({
			...m,
			date: new Date(m.timestamp * 1000),
			timestamp: m.timestamp,
		}));

		setData(formatted);

		const statData = await fetchStats(host);
		setStats(statData);
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
	}, [host]);

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
		<div className="grid">
		<div className="card">
		<h2>CPU</h2>
		<CpuChart data={data} />
		</div>

		<div className="card">
		<h2>RAM</h2>
		<RamChart data={data} />
		</div>
		</div>
		</div>
		</div>

	);
}

