import { useState } from 'react';
import { Link } from 'react-router-dom';

function Broadcast() {
  const [message, setMessage] = useState('');
  const [image, setImage] = useState(null);
  const [success, setSuccess] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleImageChange = (e) => {
    if (e.target.files && e.target.files[0]) {
      setImage(e.target.files[0]);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSuccess('');
    setError('');
    setLoading(true);

    const formData = new FormData();
    formData.append('message', message);
    if (image) {
      formData.append('image', image);
    }

    try {
      const response = await fetch(`${import.meta.env.VITE_API_URL}/broadcast-message`, {
        method: 'POST',
        body: formData,
      });

      const data = await response.json();

      if (response.ok) {
        setSuccess(data.message);
        setMessage('');
        setImage(null);
        document.getElementById('fileInput').value = "";
      } else {
        setError(data.error || 'Failed to send broadcast');
      }
    } catch {
      setError('Could not connect to server');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container">
      <h2>Broadcast Message</h2>
      <p>Send a message to all members.</p>
      
      <form onSubmit={handleSubmit}>
        <div>
          <label>Message Text:</label>
          <textarea
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            rows="4"
            placeholder="Enter your message here..."
          />
        </div>

        <div>
          <label>Image (Optional):</label>
          <input
            id="fileInput"
            type="file"
            accept="image/*"
            onChange={handleImageChange}
          />
        </div>

        <button type="submit" disabled={loading}>
          {loading ? 'Sending...' : 'Send Broadcast'}
        </button>
      </form>

      {error && <p className="error-text">{error}</p>}
      {success && <p className="success-text">{success}</p>}

      <div>
        <Link to="/">
          <button>Back to Home</button>
        </Link>
      </div>
    </div>
  );
}

export default Broadcast;
