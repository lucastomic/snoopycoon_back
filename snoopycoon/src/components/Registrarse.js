import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Link } from 'react-router-dom';


function Registrarse() {
  // Estados para cada campo del formulario
  const navigate = useNavigate();

  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');


  
  // Función para manejar el registro
  const handleRegister = async (e) => {
    e.preventDefault();

    // Validar que todos los campos estén completos
    if (!name || !email || !password || !confirmPassword) {
      setError('Por favor, rellene todos los campos.');
      return;
    }

    // Validar que las contraseñas coincidan
    if (password !== confirmPassword) {
      setError('Las contraseñas no coinciden.');
      return;
    }

    try {
      // Enviar la solicitud de registro a la API
      const response = await fetch('http://localhost:8080/api/auth/signup', {
        method: 'POST',
        credentials: "include",
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name, email, password }),
      });

      const data = await response.json();

      if (data.token !== undefined) {
        console.log('Registro exitoso:', data.message);
        // Aquí puedes redirigir al usuario o mostrar un mensaje de éxito
      } else {
        setError(data.error);
      }
    } catch (error) {
      console.error('Hubo un error en la solicitud:', error);
      setError('Hubo un error al procesar la solicitud.');
    }
  };
  const handleGoBack = () => {
    navigate(-1);
}


  return (
    <div className="register-page">
      {error && <p style={{ color: 'red' }}>{error}</p>}

      <form onSubmit={handleRegister} className="register-form">
        <h2>Registrarse</h2>

        {/* Campo Nombre */}
        <input
          type="text"
          placeholder="Nombre"
          value={name}
          onChange={(e) => setName(e.target.value)}
          style={{ width: '300px', height: '40px', fontSize: '16px', padding: '10px' }}
          required
        />

        {/* Campo Email */}
        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          style={{ width: '300px', height: '40px', fontSize: '16px', padding: '10px' }}
          required
        />

        {/* Campo Contraseña */}
        <input
          type="password"
          placeholder="Contraseña"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          style={{ width: '300px', height: '40px', fontSize: '16px', padding: '10px' }}
          required
        />

        {/* Confirmar Contraseña */}
        <input
          type="password"
          placeholder="Confirmar Contraseña"
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          style={{ width: '300px', height: '40px', fontSize: '16px', padding: '10px' }}
          required
        />

        {/* Botón para enviar el formulario */}
        <button type="submit">Registrarse</button>
      </form>
      <Link to="/">Volver a la página de inicio</Link>
    </div>
  );
}

export default Registrarse;


