import { useEffect, useState } from "react";
import ShowS from "./ShowS";

type Row = Record<string, string | number | boolean | null>;

interface Props {
  destination_id: number;
}

interface Destination extends Row {
  id: number;
  planet: string;
  state: string;
}

function DestinationView({ destination_id }: Props) {
  const [destination, setDestination] = useState<Destination>();

  useEffect(() => {
    async function loadDestination() {
      const response = await fetch(
        `http://localhost:8000/destination/${destination_id}`,
      );
      const data = await response.json();
      console.log("lo que llega", data);
      setDestination(data);
    }

    loadDestination();
  }, []);

  let name = "Destination: " + destination?.id;
  return (
    <div>
      {destination && <ShowS s_name={name} s_data={destination}></ShowS>}
    </div>
  );
}

export default DestinationView;
