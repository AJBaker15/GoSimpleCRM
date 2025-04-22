import {useEffect, useState} from 'react'

//creates the member table from a seeded .csv form in the backend. Added some extra code here to debug the 
//client server communication types. Making sure that Gin was sending a json Array. 
function MemberTable() {
    const [members, setMembers] = useState([])
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        async function fetchMembers() {
            try {
                const res = await fetch('/members')
                const data = await res.json()
                console.log("Fetched Data: ", data)
                //check to make sure Gin is sending an Array, create a fail safe if not
                if (Array.isArray(data)) {
                    setMembers(data)
                } else {
                    console.error("Expected an array but got:", data)
                    setMembers([]) //fail safe
                }
            } catch (err) {
                console.err('Error fetching members:', err)
            } finally {
                setLoading(false)
            }
        }
        fetchMembers()
    }, [])
    if (loading) return <p>Loading members....</p>
    if (!Array.isArray(members) || members.length === 0) return <p>No members found.</p>

    return (
        <div>
            <h2>Member List</h2>
            <table border="1" cellPadding="8" style={{borderCollape: 'collapse'}}>
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Street</th>
                        <th>City</th>
                        <th>State</th>
                        <th>Zip</th>
                        <th>County</th>
                        <th>Phone</th>
                        <th>Email</th>
                        <th>Last One-On-One</th>
                        <th>Issues</th>
                        <th>Due Date Payment</th>
                        <th>Active</th>
                    </tr>
                </thead>
                <tbody>
                    {members.map(member => (
                        <tr key={member.Id}>
                        <td>{member.Name}</td>
                        <td>{member.Street}</td>
                        <td>{member.City}</td>
                        <td>{member.State}</td>
                        <td>{member.Zip}</td>
                        <td>{member.County}</td>
                        <td>{member.Phone}</td>
                        <td>{member.Email}</td>
                        <td>{new Date(member.Last_one_on_one).toLocaleDateString()}</td>
                        <td>{Array.isArray(member.Issues) ? member.Issues.join(', ') : member.Issues || ''}</td>
                        <td>{new Date(member.Due_date_pay).toLocaleDateString}</td>
                        <td>{member.Active ? 'Yes' : 'No'}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    )
}
export default MemberTable