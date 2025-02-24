import React, { useState } from 'react';
import {  Link } from 'react-router-dom';

const Login = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const login = async (email, password) => {
    try {
      const response = await fetch('http://localhost:8080/api/auth/signin', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });

      const data = await response.json();

      if (data.status === 'success') {
        // Aquí manejarías el token de autenticación
        console.log('Login exitoso, token:', data.token);
        // Puedes almacenar el token en localStorage o un contexto global
        localStorage.setItem('authToken', data.token);
      } else {
        setError('Error de autenticación: ' + data.message);
      }
    } catch (error) {
      console.error('Hubo un error en la solicitud:', error);
      setError('Hubo un error al procesar la solicitud.');
    }
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    if (email && password) {
      login(email, password);
    } else {
      setError('Por favor, rellene ambos campos.');
    }
  };

  return (
    <div>
      <h2>Iniciar sesión</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Email:</label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <div>
          <label>Contraseña:</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        {error && <p style={{ color: 'red' }}>{error}</p>}
        <button type="submit">Iniciar sesión</button>
      </form>
      
    </div>
  );
};

export default Login;


