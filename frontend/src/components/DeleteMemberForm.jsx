import { useState } from 'react'

function DeleteMemberForm ({ members, onDeleteSuccess, onCancel }) { 
    const [selectedId, setSelectedId] = useState(''); 

    const handleDelete = async (e) => {
        e.preventDefault();
        if (!selectedId) return;

        try {
            const res = await fetch(`/member/${selectedId}`, {
                method: "DELETE",
            });
            if (res.ok) { 
                onDeleteSuccess();
            } else {
                alert("Failed to delete selected member ID"); 
            }
        } catch (error) {
            console.error("Error deleting member: ", error);
        }
    }; 

    return (
        <div style={{ border: "1px solid white", padding: "20px", margin: "20px" }}>
            <h3>Delete Member</h3>
            <form onSubmit={handleDelete}>
                <label>
                    Select Member to Delete: 
                    <select
                    value={selectedId}
                    onChange={(e) => setSelectedId(e.target.value)}
                    >
                        <option value="">--Select--</option>
                        {members.map((member) => (
                            <option key={member.id} value={member.id}>
                                {member.name} ({member.email})
                            </option>
                        ))}
                    </select>
                </label>
                <br />
                <button type="submit" disabled={!selectedId}>Delete</button>
                <button type="button" onClick={onCancel}>Cancel</button>
            </form>
        </div>
    );
}

export default DeleteMemberForm;