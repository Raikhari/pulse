import {
	LineChart,
	Line,
	XAxis,
	YAxis,
	Tooltip,
	ResponsiveContainer,
} from "recharts";

export default function RamChart({ data }) {
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
		<YAxis />
		<Tooltip />
		<Line type="monotone" dataKey="ram" stroke="#22c55e" />
		</LineChart>
		</ResponsiveContainer>
	);
}
