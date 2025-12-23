import { Link } from 'react-router-dom';

function Home() {
  return (
    <div className="container">
      <h1>Membership Tracker</h1>
      <p>Select an action below:</p>
      
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
    </div>
  );
}

export default Home;