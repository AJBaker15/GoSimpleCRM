import {useState} from 'react'
//creates an upload form from javascript FormData element and allows a .csv file upload. 
//also notifies the front end Dashboard to reload the members table if a new file is uploaded. 
function UploadForm({ onUploadSuccess}) {
    const [file, setFile] = useState(null)
    const [message, setMessage] = useState('')

    const handleSubmit = async (e) => {
        e.preventDefault()
        if (!file) {
            setMessage("Please select a CSV file to upload.")
            return
        }
        const formData = new FormData()
        formData.append('file', file)

        try {
            const res = await fetch('/upload', {
                method: 'POST',
                body: formData
            })
            const data = await res.json()
            if (res.ok) {
                setMessage(data.Message || 'Upload Successful')
                onUploadSuccess?.() //Reload Members
            } else {
                setMessage(data.error || 'Upload Failed')
            }
        } catch (err) {
            console.error("Upload Error: ", err)
            setMessage("An unexpected error occured.")
        }
    }
    return (
        <div>
            <h3>Upload CSV File</h3>
            <form onSubmit={handleSubmit}>
                <input type="file" accept=".csv" onChange={e => setFile(e.target.files[0])} />
                <button type="submit">Upload</button>
            </form>
        </div>
    )
}
export default UploadForm