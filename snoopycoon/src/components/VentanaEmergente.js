import { useState } from "react";

const VentanaEmergente = ({ cerrar, setHistorialBusquedas ,actualizarHistorial, busquedaInicial}) => {
  const [name, setName] = useState(busquedaInicial && busquedaInicial.name);
  const [category, setCategory] = useState(busquedaInicial && busquedaInicial.category);
  const [price_min, setPrice_min] = useState(busquedaInicial && busquedaInicial.price_min);
  const [price_max, setPrice_max] = useState(busquedaInicial && busquedaInicial.price_max);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const handleBuscar = async (e) => {
    e.preventDefault();
    setError('');

    if (!name || !category || price_min === '' || price_max === '') {
      setError('Por favor, complete todos los campos.');
      return;
    }

    if (Number(price_min) > Number(price_max)) {
      setError('El precio mínimo no puede ser mayor que el máximo.');
      return;
    }

    try {
      setIsLoading(true);
      const token = localStorage.getItem('authToken');

      const endpoint = busquedaInicial ? 'http://localhost:8080/api/listeners/' + busquedaInicial.ID : 'http://localhost:8080/api/listeners'
      const respuesta = await fetch(endpoint, {
        method: busquedaInicial ? "PUT" : "POST",
        credentials: "include",
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
          name,
          category,
          price_min: Number(price_min),
          price_max: Number(price_max),
        }),
      });
       console.log(respuesta)

      const data = await respuesta.json();

      if (respuesta.ok) {
        setHistorialBusquedas(data); 
        actualizarHistorial();
        setName('');
        setCategory('');
        setPrice_min('');
        setPrice_max('');
        cerrar(); 
      } else {
        setError(data.message || 'Error al guardar la búsqueda.');
      }
    } catch (err) {
      console.error('Error:', err);
      setError('Hubo un error al procesar la solicitud.');
    } finally {
      setIsLoading(false);
    }
    
    
  };

  return (
    <div className="ventana-emergente">
      <h2>Buscar Productos</h2>
      <form onSubmit={handleBuscar}>
        {error && <p style={{ color: 'red' }}>{error}</p>}

        <input type="text" placeholder="Nombre" value={name} onChange={(e) => setName(e.target.value)} required />
        <select value={category} onChange={(e) => setCategory(e.target.value)} required>
          <option value="">Seleccione una categoría</option>
          <option value="electronica">Electrónica</option>
          <option value="ropa">Ropa</option>
          <option value="hogar">Hogar</option>
          <option value="libros">Libros</option>
        </select>
        <input type="number" placeholder="Precio mínimo" value={price_min} onChange={(e) => setPrice_min(e.target.value)} required />
        <input type="number" placeholder="Precio máximo" value={price_max} onChange={(e) => setPrice_max(e.target.value)} required />

        <button type="submit" disabled={isLoading}>{isLoading ? 'Buscando...' : (busquedaInicial) ? 'Guardar' : 'Buscar'}</button>
      </form>

      <button onClick={cerrar}>Cerrar</button>
    </div>
  );
};

export default VentanaEmergente;






