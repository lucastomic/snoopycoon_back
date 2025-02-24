import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import VentanaEmergente from "./VentanaEmergente";
import "./Dashboard.css";

function Dashboard() {
  const navigate = useNavigate();
  const [historialBusquedas, setHistorialBusquedas] = useState([]);
  const [busquedaSeleccionada, setBusquedaSeleccionada] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");
  const [creandoBusqueda, setCreandoBusqueda] = useState(false);
  const [editandoBusqueda, setEditandoBusqueda] = useState(null);

  const [scrapeData, setScrapeData] = useState(null);
  const [isScraping, setIsScraping] = useState(false);
  const [scrapeError, setScrapeError] = useState("");

  const [query, setQuery] = useState("");

  const fetchScrapeData = async () => {
    if (!query) {
      setScrapeError("Introduce un término de búsqueda.");
      return;
    }

    setIsScraping(true);
    setScrapeError("");

    try {
      const response = await fetch(`http://localhost:8080/api/scrape?query=${query}`);
      const data = await response.json();

      if (response.ok) {
        setScrapeData(data);
      } else {
        setScrapeError(data.error || "Error al obtener datos de scraping.");
      }
    } catch (error) {
      console.error("Error en el scraping:", error);
      setScrapeError("Error en la solicitud.");
    } finally {
      setIsScraping(false);
    }
  };

  function ActualizarHistorial() {
    const obtenerHistorial = async () => {
      setIsLoading(true);
      setError("");

      try {
        const respuesta = await fetch("http://localhost:8080/api/listeners", {
          method: "GET",
          credentials: "include",
          headers: {
            "Content-Type": "application/json",
          },
        });

        const data = await respuesta.json();

        if (respuesta.ok) {
          setHistorialBusquedas(data);
        } else {
          setError(data.message || "Error al obtener el historial.");
        }
      } catch (err) {
        console.error("Error:", err);
        setError("Error al cargar el historial.");
      } finally {
        setIsLoading(false);
      }
    };
    obtenerHistorial();
  }

  useEffect(() => {
    ActualizarHistorial();
  }, []);

  const eliminarBusqueda = async (id) => {
    if (!id) {
      console.error("❌ Error: ID inválido en eliminarBusqueda", id);
      return;
    }

    try {
      const respuesta = await fetch(`http://localhost:8080/api/listeners/${id}`, {
        method: "DELETE",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (respuesta.ok) {
        setHistorialBusquedas(historialBusquedas.filter((b) => b.ID !== id));
        console.log("✅ Eliminado correctamente:", id);
      } else {
        console.error("❌ Error al eliminar:", respuesta.status);
      }
    } catch (err) {
      console.error("❌ Error en la solicitud de eliminación:", err);
    }
  };

  return (
    <div className="dashboard-page">
      <div className="sidebar">
        <div className="historial-busquedas">
          <h3>Historial de Búsquedas</h3>
          {isLoading && <p>Cargando historial...</p>}
          {!isLoading && historialBusquedas.length === 0 && (
            <p>No tienes búsquedas guardadas.</p>
          )}
          {!isLoading && historialBusquedas.length > 0 && (
            <ul>
              {historialBusquedas.map((busqueda) => (
                <div
                  key={busqueda.ID}
                  className={`busqueda-item ${
                    busquedaSeleccionada?.ID === busqueda.ID ? "seleccionada" : ""
                  }`}
                >
                  <button
                    className="boton-historial"
                    onClick={() => {
                      setQuery(busqueda.name);
                      setBusquedaSeleccionada(busqueda);
                      fetchScrapeData();
                    }}
                  >
                    {busqueda.name} <br />
                    <span className="rango-precio">
                      {busqueda.price_min}€ - {busqueda.price_max}€
                    </span>
                  </button>
                </div>
              ))}
            </ul>
          )}
        </div>

        
        <button className="boton-crear" onClick={() => setCreandoBusqueda(true)}>
          Crear Nuevo
        </button>
      </div>

      <div className="dashboard-content">
        <div className="dashboard-header">
          <h2>Bienvenido al Dashboard</h2>
          <button onClick={() => navigate("/")} className="button">
            Cerrar Sesión
          </button>
        </div>

        <div className="busqueda-container">
          <input
            type="text"
            placeholder="Buscar productos..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
        </div>

        
        <div className="acciones-busqueda">
          {busquedaSeleccionada ? (
            <>
              <button
                className="boton-editar"
                onClick={() => setEditandoBusqueda(busquedaSeleccionada)}
              >
                ✏️ Editar
              </button>

              <button
                className="boton-eliminar"
                onClick={() => eliminarBusqueda(busquedaSeleccionada.ID)}
              >
                ❌ Eliminar
              </button>
            </>
          ) : (
            <button onClick={fetchScrapeData}>Buscar en Wallapop + Vinted</button>
          )}
        </div>

        
        <div className="estadisticas">
          {creandoBusqueda ? (
            <VentanaEmergente
              cerrar={() => setCreandoBusqueda(false)}
              setHistorialBusquedas={setHistorialBusquedas}
              actualizarHistorial={ActualizarHistorial}
            />
          ) : editandoBusqueda ? (
            <VentanaEmergente
              cerrar={() => setEditandoBusqueda(null)}
              setHistorialBusquedas={setHistorialBusquedas}
              actualizarHistorial={ActualizarHistorial}
              busquedaInicial={editandoBusqueda}
            />
          ) : (
            <>
              <h3>STATS</h3>
              {isScraping ? (
                <p>Cargando datos...</p>
              ) : scrapeError ? (
                <p style={{ color: "red" }}>{scrapeError}</p>
              ) : scrapeData ? (
                <>
                  <h4>📌 Resultados de Wallapop</h4>
                  <p>
                    Total productos:{" "}
                    <strong>{scrapeData.wallapop?.total_items || 0}</strong>
                  </p>
                  <p>
                    Precio medio:{" "}
                    <strong>
                      {scrapeData.wallapop?.average_price?.toFixed(2) || 0}€
                    </strong>
                  </p>

                  <h4>📌 Resultados de Vinted</h4>
                  <p>
                    Total productos:{" "}
                    <strong>{scrapeData.vinted?.total_items || 0}</strong>
                  </p>
                  <p>
                    Precio medio:{" "}
                    <strong>{scrapeData.vinted?.average_price?.toFixed(2) || 0}€</strong>
                  </p>
                </>
              ) : (
                <p>No hay datos de scraping disponibles.</p>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  );
}

export default Dashboard;



































