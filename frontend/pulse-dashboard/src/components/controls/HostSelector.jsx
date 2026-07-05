export default function HostSelector({ hosts, selected, onChange }) {
  return (
    <select value={selected} onChange={(e) => onChange(e.target.value)}>
      {hosts.map((h) => (
        <option key={h} value={h}>
          {h}
        </option>
      ))}
    </select>
  );
}
