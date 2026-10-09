import { useEffect, useState } from "react";
import ShowS from "./ShowS";

type Row = Record<string, string | number | boolean | null>;

interface Tenant extends Row {
  id: number;
  name: string;
  tax_id: string;
}

interface Props {
  tenant_id: number;
}

function TenantView({ tenant_id }: Props) {
  const [tenant, setTenant] = useState<Tenant>();

  useEffect(() => {
    async function loadTenant() {
      const response = await fetch(`http://localhost:8000/tenant/${tenant_id}`);
      const data = await response.json();
      console.log("lo que llega", data);
      setTenant(data);
    }

    loadTenant();
  }, []);

  let name = "Tenant: " + tenant?.id;
  return <div>{tenant && <ShowS s_name={name} s_data={tenant}></ShowS>}</div>;
}

export default TenantView;
