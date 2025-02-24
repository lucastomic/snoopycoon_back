
// src/components/LoginForm.js
import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';

function LoginForm() {
  const navigate = useNavigate();

  // Estados para los campos del formulario
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  // Función para manejar el inicio de sesión
  const handleLogin = async (e) => {
    e.preventDefault();

    // Validar que ambos campos estén completos
    if (!email || !password) {
      setError('Por favor, rellene ambos campos.');
      return;
    }

    try {
      // Enviar la solicitud de inicio de sesión a la API
      const response = await fetch('http://localhost:8080/api/auth/signin', {
        method: 'POST',
        credentials: "include",
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });

      const data = await response.json();
      if (response.ok) {
        console.log('Inicio de sesión exitoso, token:', data.token);
        // Almacenar el token en localStorage
        localStorage.setItem('authToken', data.token);
        // Redirigir al usuario a la página de dashboard
        navigate('/dashboard');
      } else {
        setError('Error de autenticación: ' + data.message);
      }
    } catch (error) {
      console.error('Hubo un error en la solicitud:', error);
      setError('Hubo un error al procesar la solicitud.');
    }
  };

  return (
    <div className="login-page">
      {error && <p style={{ color: 'red' }}>{error}</p>}

      <form onSubmit={handleLogin} className="login-form">
        <h2>Iniciar Sesión</h2>

        {/* Campo Email */}
        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="input"
          required
        />

        {/* Campo Contraseña */}
        <input
          type="password"
          placeholder="Contraseña"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className="input"
          required
        />

        {/* Botón para enviar el formulario */}
        <button type="submit" className="button">Iniciar Sesión</button>
      </form>

      {/* Enlace para registrarse */}
      <Link to="/registrarse" className="link">
        ¿No tienes cuenta? Regístrate aquí
      </Link>
    </div>
  );
}

export default LoginForm;








