import { useEffect, useState } from "react";
import ShowS from "./ShowS";
import TableX from "./TableX";

interface HyperRocket extends Row {
  id: number;
  fuel: number;
  parts: string;
  x: number;
  y: number;
  z: number;
  date: string;
  state: string;
}

type Row = Record<string, string | number | boolean | null>;

interface Rocket extends Row {
  satelite: number;
  id: number;
}

interface SateliteInfo {
  rocket: Rocket;
  h_rockets: HyperRocket[];
}

interface Props {
  rocket_id: number;
}

function RocketView({ rocket_id }: Props) {
  const [rocketInfo, setRocketInfo] = useState<SateliteInfo>();

  useEffect(() => {
    async function loadSateliteInfo() {
      const response = await fetch(
        `http://localhost:8000/hyper-rocket/${rocket_id}`,
      );
      const data = await response.json();
      console.log("lo que llega", data);
      setRocketInfo(data);
      console.log("lo que se parsea", rocketInfo);
    }

    loadSateliteInfo();
  }, []);
  useEffect(() => {
    console.log("aqui la info cambió:", rocketInfo);
  }, [rocketInfo]);
  let name = "Rocket: : " + rocketInfo?.rocket.id;
  return (
    <div>
      {rocketInfo?.rocket && (
        <ShowS s_name={name} s_data={rocketInfo.rocket}></ShowS>
      )}
      {rocketInfo?.h_rockets && (
        <TableX
          atributes={Object.keys(rocketInfo.h_rockets[0])}
          list={rocketInfo.h_rockets}
        ></TableX>
      )}
    </div>
  );
}

export default RocketView;
