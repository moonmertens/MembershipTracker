import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Home from './pages/Home';
import AddMember from './pages/AddMember';
import UpdateMember from './pages/UpdateMember'; // 1. Import the new page

import './App.css';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/add" element={<AddMember />} />
        <Route path="/update" element={<UpdateMember />} /> {/* 2. Add the route */}
      </Routes>
    </BrowserRouter>
  );
}

export default App;