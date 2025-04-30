import {useState} from 'react'

function AddMemberForm({ onAddSuccess, onCancel}) {
    const [formData, setFormData] = useState({
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
        active: false
    });

    const handleChange = (e) => {
        const {name, value, type, checked} = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: type === 'checkbox' ? checked :value
        }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
    
        const convertToISO = (dateStr) => {
            const [month, day, year] = dateStr.split('/');
            const fullYear = year.length === 2 ? '20' + year : year;
            return new Date(`${fullYear}-${month.padStart(2, '0')}-${day.padStart(2, '0')}`).toISOString();
        };
    
        const payload = {
            ...formData,
            last_one_on_one: convertToISO(formData.last_one_on_one),
            due_date_pay: convertToISO(formData.due_date_pay),
            issues: formData.issues.split(',').map(issue => issue.trim())
        };
    
        try {
            const res = await fetch("/members/add", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(payload),
            });
    
            if (res.ok) {
                onAddSuccess();
            } else {
                alert("Failed to add member.");
            }
        } catch (error) {
            console.error("Error adding member:", error);
        }
    };
    
    return (
        <div style={{ border: "1px solid white", padding: "20px", margin: "20px" }}>
            <h3>Add New Member</h3>
            <form onSubmit={handleSubmit}>
                <input name = "name" placeholder = "Name" value={formData.name} onChange={handleChange} /><br/>
                <input name = "street" placeholder = "Street" value={formData.street} onChange={handleChange} /><br/>
                <input name = "city" placeholder = "City" value={formData.city} onChange={handleChange} /><br/>
                <input name = "state" placeholder = "State" value={formData.state} onChange={handleChange} /><br />
                <input name = "zip" placeholder = "Zip" value={formData.zip} onChange={handleChange} /><br/>
                <input name = "county" placeholder = "County" value={formData.county} onChange={handleChange} /><br/>
                <input name = "phone" placeholder = "Phone" value={formData.phone} onChange={handleChange} /><br/>
                <input name = "email" placeholder = "Email" value={formData.email} onChange={handleChange} /><br/>
                <input name = "last_one_on_one" placeholder = "Last One-on-One" value={formData.last_one_on_one} onChange={handleChange} /><br/>
                <input name = "issues" placeholder = "Issues" value={formData.issues} onChange={handleChange} /><br/>
                <input name = "due_date_pay" placeholder = "Last Dues Payment" value={formData.due_date_pay} onChange={handleChange} /><br/>
                <label>
                    Active:
                    <input type="checkbox" name = "active" checked = {formData.active} onChange={handleChange} />
                </label><br />
                <button type="submit">Submit</button>
                <button type="button" onClick={onCancel}>Cancel</button>
            </form>
        </div>
    );
}

export default AddMemberForm;