// src/components/ProtectedRoute.js
import React from 'react';
import { Navigate } from 'react-router-dom';

function ProtectedRoute({ children }) {
  const authToken = localStorage.getItem('authToken');

  if (!authToken) {
    // Si no hay token, redirige al inicio de sesión
    return <Navigate to="/" replace />;
  }

  // Si hay token, permite acceder a la ruta protegida
  return children;
}

export default ProtectedRoute;
