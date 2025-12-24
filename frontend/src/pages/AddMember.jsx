import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

function AddMember() {
  const [phone, setPhone] = useState('');
  const [name, setName] = useState('');
  const [visits, setVisits] = useState('');

  const [success, setSuccess] = useState('');
  const [error, setError] = useState('');

  const navigate = useNavigate(); 

  const handleSubmit = async () => {
    setError(''); // Clear previous errors and successes
    setSuccess('')

    try {
      const response = await fetch(`${import.meta.env.VITE_API_URL}/add-member`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          phone_number: parseInt(phone),
          name: name,
          visits: parseInt(visits)
        })
      });

      if (response.ok) {
        setSuccess("User added")
      } else {
        const data = await response.json();
        setError(data.error);
      }
    } catch {
      setError("Could not connect to server");
    }
  };

  const handleClear = () => {
    setPhone('');
    setName('');
    setVisits('');
    setSuccess('');
    setError('');
  };

  return (
    <div className="container">
      <h2>Register New Member</h2>
      
      <input 
        type="number" 
        placeholder="Phone Number" 
        value={phone}
        onChange={(e) => setPhone(e.target.value)}
      />
      
      <input 
        type="text" 
        placeholder="Full Name" 
        value={name}
        onChange={(e) => setName(e.target.value)}
      />

      <input 
        type="number" 
        placeholder="Visits" 
        value={visits}
        onChange={(e) => setVisits(e.target.value)}
      />

      {/* Feedback Messages */}
      {error && <p className="error-text">{error}</p>}
      {success && <p style={{color: 'green'}}>{success}</p>}

      <button onClick={handleSubmit}>Add Member</button>

      <button onClick={handleClear}>
        Clear
      </button>
      
      <button onClick={() => navigate("/")}>
        Home
      </button>
    </div>
  );
}

export default AddMember;