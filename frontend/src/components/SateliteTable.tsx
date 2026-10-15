import { useEffect, useState } from "react";
import TableX from "./TableX";

type Row = Record<string, string | number | boolean | null>;

interface Satelite extends Row {
  tenant: number;
  id: number;
}

function SateliteTable() {
  const [satelites, setSatelites] = useState<Satelite[]>();

  useEffect(() => {
    async function loadSatelites() {
      const response = await fetch(`http://localhost:8000/satelite/`); //2
      const data = await response.json();
      console.log("lo que llega", data);
      setSatelites(data);
    }

    loadSatelites();
  }, []);

  return (
    <div>
      <h1>SATELITES</h1>
      {satelites && (
        <TableX atributes={Object.keys(satelites[0])} list={satelites}></TableX>
      )}
    </div>
  );
}

export default SateliteTable;
