import React from 'react';

function Boton({ paramFunc, children }) {
  return (
    <button className="button" onClick={paramFunc}>
      {children}
    </button>
  );
}

export default Boton;






