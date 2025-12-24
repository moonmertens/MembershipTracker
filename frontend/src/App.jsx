import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Home from './pages/Home';
import AddMember from './pages/AddMember';
import UpdateMember from './pages/UpdateMember';
import ListMembers from './pages/ListMembers';

import './App.css';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/add" element={<AddMember />} />
        <Route path="/update" element={<UpdateMember />} />
        <Route path="/list" element={<ListMembers />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;