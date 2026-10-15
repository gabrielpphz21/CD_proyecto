import { useEffect, useState } from "react";
import TableX from "./TableX";

type Row = Record<string, string | number | boolean | null>;

interface TakeOff extends Row {
  id: number;
  rocket: number;
  date: string;
  runway: number;
}

function TakeOffTable() {
  const [takeOffs, setTakeOffs] = useState<TakeOff[]>();

  useEffect(() => {
    async function loadTakeOffs() {
      const response = await fetch(`http://localhost:8000/take-off/`); //2
      const data = await response.json();
      console.log("lo que llega", data);
      setTakeOffs(data);
    }

    loadTakeOffs();
  }, []);

  return (
    <div>
      <h1>DESTINATIONS</h1>
      {takeOffs && (
        <TableX atributes={Object.keys(takeOffs[0])} list={takeOffs}></TableX>
      )}
    </div>
  );
}

export default TakeOffTable;
