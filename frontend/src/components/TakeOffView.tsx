import { useEffect, useState } from "react";
import ShowS from "./ShowS";

type Row = Record<string, string | number | boolean | null>;

interface TakeOff extends Row {
  id: number;
  rocket: number;
  date: string;
  runway: number;
}

interface Props {
  takeOff_id: number;
}

function TakeOffView({ takeOff_id }: Props) {
  const [takeOff, setTakeOff] = useState<TakeOff>();

  useEffect(() => {
    async function loadTakeOff() {
      const response = await fetch(
        `http://localhost:8000/take-off/${takeOff_id}`,
      );
      const data = await response.json();
      console.log("lo que llega", data);
      setTakeOff(data);
    }

    loadTakeOff();
  }, []);

  let name = "Take Off: " + takeOff?.id;
  return <div>{takeOff && <ShowS s_name={name} s_data={takeOff}></ShowS>}</div>;
}

export default TakeOffView;
