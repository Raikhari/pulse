import { useState, useEffect } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from './assets/vite.svg'
import heroImg from './assets/hero.png'
import './App.css'
import CpuChart from "./components/CpuChart";
import RamChart from "./components/RamChart";
import HostSelector from "./components/HostSelector";
import { fetchMetrics, fetchHosts } from "./api/metrics";

export default function App() {
	const [hosts, setHosts] = useState([]);
	const [host, setHost] = useState("");
	const [data, setData] = useState([]);

	async function load() {
		if (!host) return;

		const metrics = await fetchMetrics(host);

		const formatted = metrics.map(m => ({
			...m,
			time: new Date(m.timestamp * 1000).toLocaleTimeString()
		}));

		setData(formatted);
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
		<HostSelector hosts={hosts} selected={host} onChange={setHost} />
		</header>
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
	);
}

