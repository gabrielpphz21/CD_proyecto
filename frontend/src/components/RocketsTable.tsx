import { useEffect, useState } from "react";
import TableX from "./TableX";

type Row = Record<string, string | number | boolean | null>;

interface Rocket extends Row {
  tenant: number;
  id: number;
}

function RocketTable() {
  const [rockets, setRockets] = useState<Rocket[]>();

  useEffect(() => {
    async function loadRockets() {
      const response = await fetch(`http://localhost:8000/rocket/`); //2
      const data = await response.json();
      console.log("lo que llega", data);
      setRockets(data);
    }

    loadRockets();
  }, []);

  return (
    <div>
      <h1>ROCKETS</h1>
      {rockets && (
        <TableX atributes={Object.keys(rockets[0])} list={rockets}></TableX>
      )}
    </div>
  );
}

export default RocketTable;
