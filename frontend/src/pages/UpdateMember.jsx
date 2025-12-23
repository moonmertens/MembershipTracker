import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

function UpdateMember() {
  const [phone, setPhone] = useState('');
  const [member, setMember] = useState(null);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleSearch = async () => {
    setMember(null);
    setError('');

    try {
      const response = await fetch(`http://localhost:8080/get-member?phone_number=${phone}`);
      
      if (response.ok) {
        const data = await response.json();
        setMember(data);
      } else {
        const data = await response.json();
        setError(data.error);
      }
    } catch {
      setError("Server is offline");
    }
  };

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

      {/* Consistent Error Display */}
      {error && <p className="error-text">{error}</p>}

      {/* Member Details Card */}
      {member && (
        <div className="card">
          <h3>Member Found</h3>
          <p><strong>Name:</strong> {member.name}</p>
          <p><strong>Phone:</strong> {member.phone_number}</p>
          <p><strong>Visits:</strong> {member.visits}</p>
        </div>
      )}

      <br /><br />
      <button onClick={() => navigate("/")}>
        Back to Home
      </button>
    </div>
  );
}

export default UpdateMember;