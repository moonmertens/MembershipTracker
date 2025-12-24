import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

function ListMembers() {
  const [members, setMembers] = useState(null);

  const [error, setError] = useState('');
//   const [success, setSuccess] = useState('');

  const navigate = useNavigate();

  useEffect(() => {
    const fetchMembers = async () => {
      try {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/get-all-members`);
        if (!response.ok) {
          throw new Error('Failed to fetch members');
        }
        const data = await response.json();
        setMembers(data);
      } catch (err) {
        setError('Could not load members. Is the server running?');
        console.error(err);
      }
    };

    fetchMembers();
  }, []);

  return (
    <div className="container">
      <h2>All Members</h2>
      
      {error && <p className="error-text">{error}</p>}

      {members && members.length === 0 && (
        <p>No members found.</p>
      )}

      {members && members.length > 0 && (
        <div className="table-container">
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>Phone</th>
                <th>Visits</th>
              </tr>
            </thead>
            <tbody>
              {members.map((member) => (
                <tr key={member.phone_number}>
                  <td>{member.name}</td>
                  <td>{member.phone_number}</td>
                  <td>{member.visits}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <br />
      <button onClick={() => navigate("/")}>
        Home
      </button>
    </div>
  );
}

export default ListMembers;
