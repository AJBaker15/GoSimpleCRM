import MemberTable from "./MemberTable";
import UploadForm from "./UploadForm";
import {useState} from 'react'

//displays the members table as the main interface after logging in. 
function Dashboard() {
//need this to change/reload the members table if there is a new upload -- manipulating the useState.
const [reloadFlag, setReloadFlag] = useState(false);
const [showUploadForm, setShowUploadForm] = useState(false);

    return (
        <div> 
            <div style ={{ marginBottom : '20px'}} >
                <button onClick={() => alert('Add member not yet implemented.')}> Add Member</button>
                <button onClick={() => setShowUploadForm(prev => !prev)}> Upload CSV File</button>
                <button onClick={() => alert('Update Member not implemented yet. ')}>Update Member</button>
                <button onClick={() => alert('Delete Member not implemented yet.')}> Delete Member</button>
            </div>

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