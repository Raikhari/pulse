function formatUptime(seconds) {
    if (seconds == null) return "--";

    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);

    if (days > 0) {
        return `${days}d ${hours}h`;
    }

    if (hours > 0) {
        return `${hours}h ${minutes}m`;
    }

    return `${minutes}m`;
}

function getStatus(timestamp) {
    if (!timestamp) {
        return {
            label: "Unknown",
            className: "status-unknown",
        };
    }

    const ageSeconds = Date.now() / 1000 - timestamp;

    if (ageSeconds < 15) {
        return {
            label: "Reporting",
            className: "status-online",
        };
    }

    if (ageSeconds < 60) {
        return {
            label: "Delayed",
            className: "status-delayed",
        };
    }

    return {
        label: "Not reporting",
        className: "status-offline",
    };
}

function formatLastSeen(timestamp) {
    if (!timestamp) return "--";

    const seconds = Math.max(
        0,
        Math.floor(Date.now() / 1000 - timestamp)
    );

    if (seconds < 60) return `${seconds}s ago`;

    const minutes = Math.floor(seconds / 60);
    if (minutes < 60) return `${minutes}m ago`;

    const hours = Math.floor(minutes / 60);
    return `${hours}h ago`;
}

export default function HostStatusCard({ latest }) {
    if (!latest) {
        return (
            <div className="card host-status-card">
                <h2>Host Status</h2>
                <p>No current host data available.</p>
            </div>
        );
    }

    const status = getStatus(latest.timestamp);

    return (
        <div className="card host-status-card">
            <div className="host-status-header">
                <h2>{latest.hostname}</h2>

                <span className={`status-badge ${status.className}`}>
                    {status.label}
                </span>
            </div>

            <div className="host-status-grid">
                <div>
                    <span className="host-label">Last seen</span>
                    <strong>{formatLastSeen(latest.timestamp)}</strong>
                </div>

                <div>
                    <span className="host-label">Uptime</span>
                    <strong>{formatUptime(latest.uptime)}</strong>
                </div>

                <div>
                    <span className="host-label">CPU</span>
                    <strong>{latest.cpu.toFixed(1)}%</strong>
                </div>

                <div>
                    <span className="host-label">RAM</span>
                    <strong>{latest.ram.toFixed(1)}%</strong>
                </div>
            </div>
        </div>
    );
}
