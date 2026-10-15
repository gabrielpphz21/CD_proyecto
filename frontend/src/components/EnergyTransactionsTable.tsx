import { useEffect, useState } from "react";
import TableX from "./TableX";

type Row = Record<string, string | number | boolean | null>;

interface EnergyTransaction extends Row {
  volume: number;
  id: number;
  origin: number;
  destination: number;
  state: string;
  date: string;
}

function EnergyTransactionsTable() {
  const [energyTransactions, setEnergyTransactions] =
    useState<EnergyTransaction[]>();

  useEffect(() => {
    async function loadSatelites() {
      const response = await fetch(`http://localhost:8000/energy-transaction/`); //2
      const data = await response.json();
      console.log("lo que llega", data);
      setEnergyTransactions(data);
    }

    loadSatelites();
  }, []);

  return (
    <div>
      <h1>ENERGY TRANSACTION</h1>
      {energyTransactions && (
        <TableX
          atributes={Object.keys(energyTransactions[0])}
          list={energyTransactions}
        ></TableX>
      )}
    </div>
  );
}

export default EnergyTransactionsTable;
