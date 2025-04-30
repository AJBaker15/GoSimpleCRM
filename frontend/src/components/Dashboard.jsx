import AddMemberForm from "./AddMemberForm";
import DeleteMemberForm from "./DeleteMemberForm";
import MemberTable from "./MemberTable";
import UploadForm from "./UploadForm";
import UpdateMemberForm from "./UpdateMemberForm";
import {useState, useEffect} from 'react'

//displays the members table as the main interface after logging in. 
function Dashboard() {
//need this to change/reload the members table if there is a new upload -- manipulating the useState.
    const [reloadFlag, setReloadFlag] = useState(false);
    const [members, setMembers] = useState([]);
    const [showUploadForm, setShowUploadForm] = useState(false);
    const [showAddForm, setShowAddForm] = useState(false);
    const [showDeleteForm, setShowDeleteForm] = useState(false);
    const [showUpdateForm, setShowUpdateForm] = useState(false);

    useEffect(() => {
        async function fetchMembers() {
            try{
                const res = await fetch('/members');
                const data = await res.json();
                setMembers(data);
            } catch (err) {
                console.error("Error fetching members: ", err);
            }
        }
        fetchMembers();
    }, [reloadFlag]);

    return (
        <div> 
            <div style ={{ marginBottom : '20px'}} >
                <button onClick={() => setShowAddForm(prev => !prev)}> Add Member</button>
                <button onClick={() => setShowUploadForm(prev => !prev)}> Upload CSV File</button>
                <button onClick={() => setShowUpdateForm(prev => !prev)}> Update Member</button>
                <button onClick={() => setShowDeleteForm(prev => !prev)}> Delete Member</button>
                <button onClick={() => alert("This is not implemented yet.")}>Search Volunteers</button>
            </div>

            {showAddForm && (
                <AddMemberForm
                onAddSuccess={() => {
                    setReloadFlag(!reloadFlag); 
                    setShowAddForm(false); 
                }}
                onCancel={() => setShowAddForm(false)}
                />
            )}

            {showUpdateForm && (
                <UpdateMemberForm
                members={members}
                onUpdateSuccess={() => {
                    setReloadFlag(!reloadFlag); 
                    setShowUpdateForm(false); 
                }}
                onCancel={() => setShowUpdateForm(false)}
                />
            )}  

            {showDeleteForm && (
                <DeleteMemberForm
                members={members}
                onDeleteSuccess={() => {
                    setReloadFlag(!reloadFlag);
                    setShowDeleteForm(false);
                }}
                onCancel={() => setShowDeleteForm(false)}
                />
            )}

            {showUploadForm && (
                <UploadForm
                    onUploadSuccess={() => {
                        setReloadFlag(!reloadFlag);
                        setShowUploadForm(false);
                    }}
                />
            )}

            <MemberTable key={reloadFlag} />
        </div>
    );
}

export default Dashboard