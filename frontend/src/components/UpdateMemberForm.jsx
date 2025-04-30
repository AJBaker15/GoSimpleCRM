import { useState } from 'react';

function UpdateMemberForm({members, onUpdateSuccess, onCancel}) {
    const [selectedId, setSelectedId] = useState("");
    const [form, setForm] = useState({
        name: '', 
        street: '',
        city: '',
        state: '', 
        zip: '',
        county: '',
        phone: '',
        email: '',
        last_one_on_one: '',
        issues: '',
        due_date_pay: '',
        active: false,
    });

    function handleSelectChange(e) {
        const member = members.find(m => m.id === parseInt(e.target.value));
        if (member) {
            setSelectedId(member.id);
            setForm({
                name: member.name,
                street: member.street,
                city: member.city,
                state: member.state, 
                zip: member.zip,
                county: member.county,
                phone: member.phone,
                email: member.email,
                last_one_on_one: member.last_one_on_one.split('T')[0],
                issues: member.issues.join(', '),
                due_date_pay: member.due_date_pay.split('T')[0],
                active: member.active,
            });
        }
    }

    function handleChange(e) {
        const { name, value, type, checked } = e.target;
        setForm(prev => ({
            ...prev,
            [name]: type === 'checkbox' ? checked : value,
        }));
    }

    async function handleSubmit(e) {
        e.preventDefault();
        if (!selectedId) return;
    
        const convertToISO = (dateStr) => {
            const [year, month, day] = dateStr.split('-');
            return new Date(`${year}-${month}-${day}`).toISOString(); // YYYY-MM-DD -> ISO-8601
        };
    
        const updated = {
            ...form,
            last_one_on_one: convertToISO(form.last_one_on_one),
            due_date_pay: convertToISO(form.due_date_pay),
            issues: form.issues.split(',').map(i => i.trim()),
            id: selectedId,
        };
    
        try {
            const res = await fetch(`/member/${selectedId}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(updated),
            });
    
            if (res.ok) {
                alert('Member updated!');
                onUpdateSuccess();
            } else {
                alert("Failed to update member.");
            }
        } catch (err) {
            console.error("Update failed:", err);
        }
    }
    

    return (
        <div>
            <h3>Update Member</h3>
            <select value={selectedId} onChange={handleSelectChange}>
                <option value="">---Select Member---</option>
                {members.map(m => (
                    <option key={m.id} value={m.id}>
                        {m.name} ({m.email})
                    </option>
                ))}
            </select>

            {selectedId && (
                <form onSubmit={handleSubmit}>
                    <input name="name" value={form.name} onChange={handleChange} placeholder="Name" />
                    <input name="street" value={form.street} onChange={handleChange} placeholder="Street" />
                    <input name="city" value={form.city} onChange={handleChange} placeholder="City" />
                    <input name="state" value={form.state} onChange={handleChange} placeholder="State" />
                    <input name="zip" value={form.zip} onChange={handleChange} placeholder="Zip" />
                    <input name="county" value={form.county} onChange={handleChange} placeholder="County" />
                    <input name="phone" value={form.phone} onChange={handleChange} placeholder="Phone" />
                    <input name="email" value={form.email} onChange={handleChange} placeholder="Email" />
                    <input name="last_one_on_one" value={form.last_one_on_one} onChange={handleChange} placeholder="Last One on One" />
                    <input name="issues" value={form.issues} onChange={handleChange} placeholder="Issues" />
                    <input name="due_date_pay" value={form.due_date_pay} onChange={handleChange} placeholder="Due Date Payment" />
                    <label>
                        Active:
                        <input name="active" type="checkbox" checked={form.active} onChange={handleChange} />
                    </label>
                    <br />
                    <button type="submit">Update</button>
                    <button type="button" onClick={onCancel}>Cancel</button>
                </form>
            )}
        </div>
    );
}

export default UpdateMemberForm;
