import { useEffect, useState } from "react";
import ShowS from "./ShowS";
import TableX from "./TableX";

interface HyperSatelites extends Row {
  id: number;
  energy: number;
  parts: string;
  x: number;
  y: number;
  z: number;
  date: string;
}

type Row = Record<string, string | number | boolean | null>;

interface Satelite extends Row {
  tenant: number;
  id: number;
}

interface SateliteInfo {
  satelite: Satelite;
  h_satelite: HyperSatelites[];
}

interface Props {
  satelite_id: number;
}

function SateliteView({ satelite_id }: Props) {
  const [sateliteInfo, setSateliteInfo] = useState<SateliteInfo>();

  useEffect(() => {
    async function loadSateliteInfo() {
      const response = await fetch(
        `http://localhost:8000/hyper-satelite/${satelite_id}`,
      ); //2
      const data = await response.json();
      console.log("lo que llega", data);
      setSateliteInfo(data);
      console.log("lo que se parsea", sateliteInfo);
    }

    loadSateliteInfo();
  }, []);
  useEffect(() => {
    console.log("aqui la info cambió:", sateliteInfo);
  }, [sateliteInfo]);
  let name = "Satelite: " + sateliteInfo?.satelite.id;
  return (
    <div>
      {sateliteInfo?.satelite && (
        <ShowS s_name={name} s_data={sateliteInfo.satelite}></ShowS>
      )}
      {sateliteInfo?.h_satelite && (
        <TableX
          atributes={Object.keys(sateliteInfo.h_satelite[0])}
          list={sateliteInfo.h_satelite}
        ></TableX>
      )}
    </div>
  );
}

export default SateliteView;
