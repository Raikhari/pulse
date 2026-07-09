import { useEffect, useState } from "react";
import { fetchConfig, saveConfig } from "../../api/metrics";

export default function ThresholdSettings() {
    const [config, setConfig] = useState(null);
    const [message, setMessage] = useState("");
    const [error, setError] = useState("");
    const [saving, setSaving] = useState(false);

    useEffect(() => {
        fetchConfig()
            .then(setConfig)
            .catch((err) => setError(err.message));
    }, []);

    function handleChange(event) {
        const { name, value } = event.target;

        setConfig((current) => ({
            ...current,
            [name]: value,
        }));
    }

    async function handleSubmit(event) {
        event.preventDefault();

        setSaving(true);
        setMessage("");
        setError("");

        try {
            const payload = {
                cpu_high_threshold: Number(config.cpu_high_threshold),
                cpu_normal_threshold: Number(config.cpu_normal_threshold),
                ram_high_threshold: Number(config.ram_high_threshold),
                ram_normal_threshold: Number(config.ram_normal_threshold),
            };

            const updated = await saveConfig(payload);

            setConfig(updated);
            setMessage("Thresholds saved.");
        } catch (err) {
            setError(err.message);
        } finally {
            setSaving(false);
        }
    }

    if (!config) {
        return (
            <div className="card threshold-settings">
                <h2>Threshold Settings</h2>
                <p>{error || "Loading settings..."}</p>
            </div>
        );
    }

    return (
        <div className="card threshold-settings">
            <h2>Threshold Settings</h2>

            <form onSubmit={handleSubmit}>
                <div className="threshold-grid">
                    <label>
                        CPU high
                        <input
                            type="number"
                            name="cpu_high_threshold"
                            value={config.cpu_high_threshold}
                            onChange={handleChange}
                            min="0"
                            max="100"
                            step="1"
                        />
                    </label>

                    <label>
                        CPU normal
                        <input
                            type="number"
                            name="cpu_normal_threshold"
                            value={config.cpu_normal_threshold}
                            onChange={handleChange}
                            min="0"
                            max="100"
                            step="1"
                        />
                    </label>

                    <label>
                        RAM high
                        <input
                            type="number"
                            name="ram_high_threshold"
                            value={config.ram_high_threshold}
                            onChange={handleChange}
                            min="0"
                            max="100"
                            step="1"
                        />
                    </label>

                    <label>
                        RAM normal
                        <input
                            type="number"
                            name="ram_normal_threshold"
                            value={config.ram_normal_threshold}
                            onChange={handleChange}
                            min="0"
                            max="100"
                            step="1"
                        />
                    </label>
                </div>

                <button type="submit" disabled={saving}>
                    {saving ? "Saving..." : "Save thresholds"}
                </button>

                {message && (
                    <p className="settings-success">{message}</p>
                )}

                {error && (
                    <p className="settings-error">{error}</p>
                )}
            </form>
        </div>
    );
}
