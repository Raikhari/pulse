export default function StatCard({ title, value, unit }) {
    return (
        <div className="card stat-card">
            <h3>{title}</h3>

            <div className="stat-value">
                {value}
                <span>{unit}</span>
            </div>
        </div>
    );
}
