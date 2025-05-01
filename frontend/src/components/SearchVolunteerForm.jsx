import { useState } from 'react';

function SearchVolunteerForm({ members }) {
    const [selectedIssue, setSelectedIssue] = useState('');
    const [results, setResults ] = useState([]);

    //find all of the unique issues in the members' issue list
    const allIssues = members.flatMap(m => m.issues || []);
    //put the issues in a set which does not allow duplicates, sort them.
    const uniqueIssues = Array.from(new Set(allIssues)).sort();

    async function handleSearch(e) {
        e.preventDefault();
        if(!selectedIssue) return;

        try {
            const res = await fetch(`/members/search?issue=${encodeURIComponent(selectedIssue)}`);
            const data = await res.json();
            setResults(data);
        } catch (err) {
            console.error('Search error: ', err);
        }
    }

    return (
        <div style={{ border: '1px solid white', padding: '20px', margin: '20px' }}>
            <h3> Search Volunteers by Issue </h3>
            <form onSubmit={handleSearch}>
                <label>
                    Select Issue: 
                    <select value={selectedIssue} onChange={(e) => setSelectedIssue(e.target.value)}>
                        <option value="">--Select--</option>
                        {uniqueIssues.map(issue => (
                            <option key={issue} value={issue}>{issue}</option>
                        ))}
                    </select>
                </label>
                <button type="submit" disabled={!selectedIssue}>Search</button>
            </form>
            {results.length > 0 && (
                <div>
                    <h4>Matching Members:</h4>
                    <ul>
                        {results.map(m => (
                            <li key={m.id}>{m.name} ({m.email}) - Issues: {Array.isArray(m.issues) ? m.issues.join(', ') : m.issues}</li>
                        ))}
                    </ul>
                </div>
            )}
        </div>
    );
}

export default SearchVolunteerForm;