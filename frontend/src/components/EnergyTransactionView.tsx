import { useEffect, useState } from "react";
import ShowS from "./ShowS";

type Row = Record<string, string | number | boolean | null>;

interface EnergyTransaction extends Row {
  volume: number;
  id: number;
  origin: number;
  destination: number;
  state: string;
  date: string;
}

interface Props {
  energyTransaction_id: number;
}

function EnergyTransactionView({ energyTransaction_id }: Props) {
  const [energyTransaction, setEnergyTransaction] =
    useState<EnergyTransaction>();

  useEffect(() => {
    async function loadEnergyTransactionInfo() {
      const response = await fetch(
        `http://localhost:8000/energy-transaction/${energyTransaction_id}`,
      );
      const data = await response.json();
      console.log("lo que llega", data);
      setEnergyTransaction(data);
      console.log("lo que se parsea", energyTransaction);
    }

    loadEnergyTransactionInfo();
  }, []);

  let name = "Energy Transaction: " + energyTransaction?.id;
  return (
    <div>
      {energyTransaction && (
        <ShowS s_name={name} s_data={energyTransaction}></ShowS>
      )}
    </div>
  );
}

export default EnergyTransactionView;
