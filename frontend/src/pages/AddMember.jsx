import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

function AddMember() {
  const [phone, setPhone] = useState('');
  const [name, setName] = useState('');
  const [error, setError] = useState(''); 
  const navigate = useNavigate(); 

  const handleSubmit = async () => {
    setError(''); // Clear previous errors

    try {
      const response = await fetch("http://localhost:8080/add-member", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          phone_number: parseInt(phone),
          name: name,
          visits: 0
        })
      });

      if (response.ok) {
        // Success: Redirect immediately
        navigate("/"); 
      } else {
        const data = await response.json();
        setError(data.error);
      }
    } catch {
      setError("Could not connect to server");
    }
  };

  return (
    <div className="container">
      <h2>Register New Member</h2>
      
      <input 
        type="number" 
        placeholder="Phone Number (e.g. 81234567)" 
        value={phone}
        onChange={(e) => setPhone(e.target.value)}
      />
      
      <input 
        type="text" 
        placeholder="Full Name" 
        value={name}
        onChange={(e) => setName(e.target.value)}
      />

      {/* Consistent Error Display */}
      {error && <p className="error-text">{error}</p>}

      <button onClick={handleSubmit}>Save Member</button>
      
      <br /><br />
      <button onClick={() => navigate("/")}>
        Back to Home
      </button>
    </div>
  );
}

export default AddMember;