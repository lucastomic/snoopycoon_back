
// src/App.js
import './App.css';
import LoginForm from './components/LoginForm'; // Asegúrate de que el archivo se llama LoginForm.js
import React from 'react';
import { BrowserRouter as Router, Route, Link, Routes } from 'react-router-dom';
import Registrarse from './components/Registrarse';
import Dashboard from './components/Dashboards';
import ProtectedRoute from './components/ProtectedRoute'; // Importa ProtectedRoute

function App() {
  return (
    <Router>
      <div className="App">
        <div className="container">
          <div className="content">
            <p></p>
          </div>
          <img
            src="images/logo.png"
            alt="Imagen"
            width="150"
            height="auto"
            className="image"
          />
        </div>

        <header className="App-header">
          <p></p>
          
          {/* Definir rutas */}
          <Routes>
            <Route path="/" element={<LoginForm />} />
            <Route path="/registrarse" element={<Registrarse />} />
            <Route 
              path="/dashboard" 
              element={
                <ProtectedRoute>
                  <Dashboard />
                </ProtectedRoute>
              } 
            /> {/* Ruta protegida para Dashboard */}
          </Routes>
          
          <Link to="/registrarse" className="link">
            ¿No tienes cuenta? Regístrate aquí
          </Link>
        </header>
      </div>
    </Router>
  );
}

export default App;



