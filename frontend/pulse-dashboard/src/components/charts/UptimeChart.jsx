import {
	LineChart,
	Line,
	XAxis,
	YAxis,
	Tooltip,
	ResponsiveContainer,
} from "recharts";

export default function UptimeChart({ data }) {
	return (
		<ResponsiveContainer width="100%" height={300}>
		<LineChart data={data}>
		<XAxis
		dataKey="date"
		tickFormatter={(value) =>
			value.toLocaleTimeString([], {
				hour: "2-digit",
				minute: "2-digit",
			})
		}
		/>
		<YAxis 
		tickFormatter={(v) => `${Math.round(v)}h`}
		/>
		<Tooltip />
		<Line type="monotone" dataKey="uptimeHours" stroke="#a855f7" />
		</LineChart>
		</ResponsiveContainer>
	);
}

