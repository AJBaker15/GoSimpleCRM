import AddMemberForm from "./AddMemberForm";
import DeleteMemberForm from "./DeleteMemberForm";
import MemberTable from "./MemberTable";
import UploadForm from "./UploadForm";
import UpdateMemberForm from "./UpdateMemberForm";
import SearchVolunteerForm from "./SearchVolunteerForm";
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
    const [showSearchForm, setShowSearchForm] = useState(false);
    const [showInactive, setShowInactive] = useState(false);

    useEffect(() => {
        async function fetchMembers() {
            try{
                const endpoint = showInactive ? '/members/inactive' : '/members';
                const res = await fetch(endpoint);
                const data = await res.json();
                setMembers(data);
            } catch (err) {
                console.error("Error fetching members: ", err);
            }
        }
        fetchMembers();
    }, [reloadFlag, showInactive]);

    return (
        <div> 
            <div style ={{ marginBottom : '20px'}} >
                <button onClick={() => setShowAddForm(prev => !prev)}> Add Member</button>
                <button onClick={() => setShowUploadForm(prev => !prev)}> Upload CSV File</button>
                <button onClick={() => setShowUpdateForm(prev => !prev)}> Update Member</button>
                <button onClick={() => setShowDeleteForm(prev => !prev)}> Delete Member</button>
                <button onClick={() => setShowSearchForm(prev => !prev)}>Search Volunteers</button>
                <button onClick={() => {setShowInactive(prev => !prev); setReloadFlag(prev => !prev);}}> {showInactive ? 'Show All Members' : 'List Inactive Members'}</button>

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

            {showSearchForm && (
                <SearchVolunteerForm
                members={members}
                onCancel={() => setShowSearchForm(false)}
                />
            )}

            <MemberTable key={reloadFlag} members={members} />
        </div>
    );
}

export default Dashboard;