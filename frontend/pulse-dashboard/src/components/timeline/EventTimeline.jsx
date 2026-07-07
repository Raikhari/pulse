function formatEventTime(timestamp) {
    return new Date(timestamp * 1000).toLocaleString([], {
        month: "short",
        day: "2-digit",
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit",
    });
}

function eventIcon(type) {
    switch (type) {
        case "reboot":
            return "🔄";
        case "high_cpu":
            return "🔥";
        case "cpu_normal":
            return "✅";
        case "high_ram":
            return "⚠️";
        case "ram_normal":
            return "✅";
        default:
            return "ℹ️";
    }
}

export default function EventTimeline({ events }) {
    if (!events || events.length === 0) {
        return (
            <div className="card">
                <h2>Recent Events</h2>
                <p className="muted">No events detected for this range.</p>
            </div>
        );
    }

    return (
        <div className="card">
            <h2>Recent Events</h2>

            <div className="timeline">
                {events.map((event, index) => (
                    <div className="timeline-item" key={`${event.timestamp}-${event.type}-${index}`}>
                        <div className="timeline-icon">
                            {eventIcon(event.type)}
                        </div>

                        <div className="timeline-content">
                            <div className="timeline-message">
                                {event.message}
                            </div>

                            <div className="timeline-meta">
                                {event.hostname} · {formatEventTime(event.timestamp)}
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
}
