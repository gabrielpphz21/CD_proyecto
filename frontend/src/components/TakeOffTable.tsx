import { useEffect, useState } from "react";
import TableX from "./TableX";

type Row = Record<string, string | number | boolean | null>;

interface Destination extends Row {
  id: number;
  planet: string;
  state: string;
}

function DestinationTable() {
  const [destinations, setDestinations] = useState<Destination[]>();

  useEffect(() => {
    async function loadDestinations() {
      const response = await fetch(`http://localhost:8000/destination/`); //2
      const data = await response.json();
      console.log("lo que llega", data);
      setDestinations(data);
    }

    loadDestinations();
  }, []);

  return (
    <div>
      <h1>TAKEOFFS</h1>
      {destinations && (
        <TableX
          atributes={Object.keys(destinations[0])}
          list={destinations}
        ></TableX>
      )}
    </div>
  );
}

export default DestinationTable;
