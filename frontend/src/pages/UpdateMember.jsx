import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

function UpdateMember() {

  // State
  const [phone, setPhone] = useState('');
  const [member, setMember] = useState(null);
  
  // Edit Form State
  const [editName, setEditName] = useState('');
  const [editVisits, setEditVisits] = useState('');

  // Feedback State
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  
  const navigate = useNavigate();

  const handleSearch = async () => {
    setMember(null);
    setError('');
    setSuccess('');

    try {
      const response = await fetch(`http://localhost:8080/get-member?phone_number=${phone}`);
      
      if (response.ok) {
        const data = await response.json();
        setMember(data);
        setEditName(data.name);
        setEditVisits(data.visits);
      } else {
        const data = await response.json();
        setError(data.error);
      }
    } catch {
      setError("Server is offline");
    }
  };

  const handleDelete = async () => {
    setError('');
    setSuccess('');

    if (!window.confirm(`Are you sure you want to delete ${member.name}?`)) {
      return;
    }

    try {
      const response = await fetch(`http://localhost:8080/delete-member?phone_number=${member.phone_number}`, {
        method: 'DELETE',
      });

      if (response.ok) {
        setSuccess("Member deleted")
        setMember(null);
        setPhone('');
        setEditName('');
        setEditVisits('');
      } else {
        const data = await response.json();
        setError(data.error)
      }
    } catch {
      setError("Could not connect to server");
    }
  };

  const performUpdate = async (newName, newVisits) => {
    setError('');
    setSuccess('');

    try {
      const response = await fetch("http://localhost:8080/update-member", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          phone_number: member.phone_number,
          name: newName,
          visits: newVisits
        })
      });

      if (response.ok) {
        setSuccess("Member updated successfully!");
        // Update local state
        setMember({ ...member, name: newName, visits: newVisits });
        setEditName(newName);
        setEditVisits(newVisits);
      } else {
        const data = await response.json();
        setError(data.error);
      }
    } catch {
      setError("Could not connect to server");
    }
  };

  const handleAdd = () => {
    performUpdate(editName, parseInt(editVisits || 0) + 1);
  };

  const handleUpdate = () => {
    performUpdate(editName, parseInt(editVisits));
  };

  // Page
  return (
    <div className="container">
      <h2>Update Member</h2>
      <p>Enter phone number to find member:</p>

      <input 
        type="number" 
        placeholder="Phone Number" 
        value={phone}
        onChange={(e) => setPhone(e.target.value)}
      />
      
      <button onClick={handleSearch}>Find Member</button>

      {/* Feedback Messages */}
      {error && <p className="error-text">{error}</p>}
      {success && <p className='success-text'>{success}</p>}

      {/* Member Details Card */}
      {member && (
        <div className="card">
          <h3>Member Details</h3>
          <p><strong>Phone:</strong> {member.phone_number}</p>

          <label>Name:</label>
          <input 
            type="text" 
            value={editName} 
            onChange={(e) => setEditName(e.target.value)} 
          />

          <label>Visits:</label>
          <input 
            type="number" 
            value={editVisits} 
            onChange={(e) => setEditVisits(e.target.value)} 
          />
          
          <button onClick={handleAdd}>
            Add 1 visit
          </button>

          <button onClick={handleUpdate}>
            Update Changes
          </button>
          
          <button onClick={handleDelete}>
            Delete Member
          </button>
        </div>
      )}

      <br /><br />
      <button onClick={() => navigate("/")}>
        Home
      </button>
    </div>
  );
}

export default UpdateMember;