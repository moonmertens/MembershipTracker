import { Link } from 'react-router-dom';

function Home() {
  return (
    <div className="container">
      <h1>Membership Tracker</h1>
      <p>Select an action below:</p>
      
      <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
        <Link to="/add">
          <button>
            Add New Member
          </button>
        </Link>

        <Link to="/update">
          <button>
            Update Member
          </button>
        </Link>

        <Link to="/list">
          <button>
            List Members
          </button>
        </Link>

        <Link to="/whatsapp">
          <button>
            Connect to WhatsApp
          </button>
        </Link>
      </div>
    </div>
  );
}

export default Home;